package mongoinfer

import (
	"fmt"
	"slices"
	"strings"

	"github.com/iancoleman/strcase"
	"go.mongodb.org/mongo-driver/bson"
)

func Header(pkg string) string {
	return `syntax = "proto3";

` + `package ` + pkg + `;

import "google/protobuf/timestamp.proto";

message MongoBinary {
  int32 _subtype = 1;
  bytes _data = 2;
}

`
}

type KV[T any] struct {
	Key   string
	Value T
}

func ToSortedKV(m map[string]*BsonType) []KV[*BsonType] {
	var kvs []KV[*BsonType]
	for k, v := range m {
		kvs = append(kvs, KV[*BsonType]{
			Key:   k,
			Value: v,
		})
	}
	slices.SortFunc(kvs, func(l, r KV[*BsonType]) int {
		return strings.Compare(l.Key, r.Key)
	})
	return kvs
}

func BsonTypeToProto(builder *strings.Builder, prefix string, t *BsonType) error {
	if t.Fields != nil {
		for k, v := range t.Fields {
			if err := BsonTypeToProto(builder, prefix+strcase.ToCamel(k), v); err != nil {
				return err
			}
		}
	} else if t.Array != nil {
		if t.Array.Array != nil {
			// not currently supported- skipping
			return nil
		}
		return BsonTypeToProto(builder, prefix, t.Array)
	} else {
		return nil
	}

	builder.WriteString("message ")
	builder.WriteString(prefix)
	builder.WriteString(" {\n")

	count := 1

	var kvs []KV[*BsonType]
	for k, v := range t.Fields {
		kvs = append(kvs, KV[*BsonType]{
			Key:   k,
			Value: v,
		})
	}
	slices.SortFunc(kvs, func(l, r KV[*BsonType]) int {
		return strings.Compare(l.Key, r.Key)
	})
	for _, kv := range kvs {
		k := kv.Key
		v := kv.Value

		if v.Type == bson.TypeArray {
			if v.Array.Type == bson.TypeArray {
				// not currently supported- skipping
				continue
			}
			builder.WriteString("  repeated ")
			v = v.Array
		} else {
			builder.WriteString("  ")
		}

		switch v.Type {
		case bson.TypeObjectID:
			builder.WriteString("string")
		case bson.TypeBinary:
			builder.WriteString("MongoBinary")
		case bson.TypeBoolean:
			builder.WriteString("bool")
		case bson.TypeDateTime:
			builder.WriteString("google.protobuf.Timestamp")
		case bson.TypeDouble:
			builder.WriteString("double")
		case bson.TypeDecimal128:
			builder.WriteString("string")
		case bson.TypeInt32:
			builder.WriteString("int32")
		case bson.TypeInt64:
			builder.WriteString("int64")
		case bson.TypeString:
			builder.WriteString("string")
		case bson.TypeEmbeddedDocument:
			builder.WriteString(prefix + strcase.ToCamel(k))
		default:
			// unknown, skipping for now
			builder.WriteString("XXUNKNOWNXX")
		}

		builder.WriteString(" ")
		builder.WriteString(k)
		builder.WriteString(fmt.Sprintf(" = %v;\n", count))
		count += 1
	}

	builder.WriteString("}\n\n")
	return nil
}
