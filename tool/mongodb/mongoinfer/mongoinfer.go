package mongoinfer

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo"
)

// Specifically chosen based on available bsonTypes
type Count [21]int

type Counts struct {
	Count
	Fields map[string]*Counts
	Array  *Counts
}

func (c Counts) String() string {
	var sb strings.Builder
	sb.WriteString("{[")
	first := true
	for i, cnt := range c.Count {
		if cnt > 0 {
			if !first {
				sb.WriteString(", ")
			}
			first = false
			sb.WriteString(fmt.Sprintf("%v %v", bsontype.Type(i), cnt))
		}
	}
	sb.WriteString("]")
	if len(c.Fields) > 0 {
		sb.WriteString(fmt.Sprintf(", fields %v", c.Fields))
	}
	if c.Array != nil {
		sb.WriteString(fmt.Sprintf(", array %v", c.Array))
	}
	sb.WriteString("}")
	return sb.String()
}

type BsonType struct {
	bsontype.Type
	Fields map[string]*BsonType
	Array  *BsonType
}

func (t BsonType) String() string {
	var sb strings.Builder
	sb.WriteString("{")
	if t.Fields != nil {
		sb.WriteString(fmt.Sprintf("doc %v", t.Fields))
	} else if t.Array != nil {
		sb.WriteString("arr ")
		sb.WriteString(t.Array.String())
	} else {
		sb.WriteString(t.Type.String())
	}
	sb.WriteString("}")
	return sb.String()
}

var NumberTypePriority = []bsontype.Type{
	// bson.TypeDecimal128,
	bson.TypeDouble,
	bson.TypeInt64,
	bson.TypeInt32,
}

func PickBestBsonType(count Count) bsontype.Type {
	count2 := count
	// Combine "number" priorities
	for i, t := range NumberTypePriority {
		if count2[t] != 0 {
			for j := i + 1; j < len(NumberTypePriority); j++ {
				otherIdx := NumberTypePriority[j]
				count2[t] += count2[otherIdx]
				count2[otherIdx] = 0
			}
		}
	}

	var sum int
	var bestCount int
	var bestType bsontype.Type
	for i, c := range count2 {
		if c > bestCount {
			bestCount = c
			bestType = bsontype.Type(i) // mapping is currently straight mapping
		}
		sum += c
	}

	return bestType
}

func CountsToBsonType(counts *Counts) *BsonType {
	bestType := PickBestBsonType(counts.Count)
	if bestType == bson.TypeEmbeddedDocument {
		fields := map[string]*BsonType{}
		for k, v := range counts.Fields {
			fields[k] = CountsToBsonType(v)
		}
		return &BsonType{
			Type:   bestType,
			Fields: fields,
		}

	} else if bestType == bson.TypeArray {
		return &BsonType{
			Type:  bestType,
			Array: CountsToBsonType(counts.Array),
		}
	}

	return &BsonType{
		Type: bestType,
	}
}

func TypeIndex(v bson.RawValue) int {
	if v.Type <= 20 {
		return int(v.Type)
	}
	return 0
}

func AddArrToCounts(arr bson.Raw, counts *Counts) error {
	rawValues, err := arr.Values()
	if err != nil {
		return err
	}
	for _, value := range rawValues {
		counts.Count[TypeIndex(value)] += 1
		if fieldDoc, ok := value.DocumentOK(); ok {
			if err := AddDocToCounts(fieldDoc, counts); err != nil {
				return err
			}
		} else if fieldArr, ok := value.ArrayOK(); ok {
			if counts.Array == nil {
				counts.Array = &Counts{
					Count:  Count{},
					Fields: map[string]*Counts{},
				}
			}
			if err := AddArrToCounts(fieldArr, counts.Array); err != nil {
				return err
			}
		}
	}
	return nil
}

func AddDocToCounts(doc bson.Raw, counts *Counts) error {
	elems, err := doc.Elements()
	if err != nil {
		return err
	}
	fields := counts.Fields
	for _, elem := range elems {
		key := elem.Key()
		value := elem.Value()
		if fieldCount, ok := fields[key]; ok {
			fieldCount.Count[TypeIndex(value)] += 1
		} else {
			newCount := Count{}
			newCount[TypeIndex(value)] += 1
			fields[key] = &Counts{
				Count:  newCount,
				Fields: map[string]*Counts{},
			}
		}
		if fieldDoc, ok := value.DocumentOK(); ok {
			if err := AddDocToCounts(fieldDoc, fields[key]); err != nil {
				return err
			}
		} else if fieldArr, ok := value.ArrayOK(); ok {
			if fields[key].Array == nil {
				fields[key].Array = &Counts{
					Count:  Count{},
					Fields: map[string]*Counts{},
				}
			}
			if err := AddArrToCounts(fieldArr, fields[key].Array); err != nil {
				return err
			}
		}
	}
	return nil
}

func BsonTypeFromSamples(ctx context.Context, col *mongo.Collection, numSamples int) (*BsonType, error) {
	res, err := col.Aggregate(ctx, mongo.Pipeline{
		{{"$sample", bson.D{{"size", numSamples}}}},
	})
	if err != nil {
		fmt.Println("hello")
		return nil, err
	}
	defer res.Close(ctx)

	counts := Counts{
		Count:  Count{},
		Fields: map[string]*Counts{},
	}

	for res.Next(ctx) {
		if err := AddDocToCounts(res.Current, &counts); err != nil {
			return nil, err
		}
		counts.Count[bson.TypeEmbeddedDocument] += 1
	}
	if res.Err() != nil {
		return nil, err
	}

	return CountsToBsonType(&counts), nil
}
