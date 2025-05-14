package mongoinfer

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
)

var DefaultServer = `server:
  cleartext: true
  hostport: :8090
  shutdowntimeout: 10s
  reflect: true
  health: true
  descriptors: out.pb # TODO
`

func DefaultInterceptor(url string) string {
	return `interceptors:
  - name: MongoDB
    config:
      url: ` + url + "\n"
}

type NamedType struct {
	Name       string
	Type       *BsonType
	Database   string
	Collection string
}

type DapiCrud struct {
	pkg     string
	service string
	ts      []NamedType
}

func NewDapiCrud(pkg string, service string, ts []NamedType) *DapiCrud {
	return &DapiCrud{
		pkg:     pkg,
		service: service,
		ts:      ts,
	}
}

func (g *DapiCrud) GenerateDapiCfg() string {
	var builder strings.Builder
	builder.WriteString("services:\n")
	builder.WriteString(fmt.Sprintf("  %v.%v:\n", g.pkg, strcase.ToCamel(g.service)))
	builder.WriteString("    endpoints:\n")

	for _, t := range g.ts {
		// TODO: remove assumption _id is an object id
		// TODO: support bytes better
		name := strcase.ToCamel(t.Name)
		builder.WriteString(fmt.Sprintf("      List%v:\n", name))
		builder.WriteString(fmt.Sprintf("        database: %v\n", t.Database))
		builder.WriteString(fmt.Sprintf("        collection: %v\n", t.Collection))
		builder.WriteString(`        options: 'options.FindOptions {Limit: req.limit}'
        paginatedfind:
          filter: '{}'
          cursor: 'req.cursor'
`)
		builder.WriteString(fmt.Sprintf("      Get%v:\n", name))
		builder.WriteString(fmt.Sprintf("        database: %v\n", t.Database))
		builder.WriteString(fmt.Sprintf("        collection: %v\n", t.Collection))
		builder.WriteString(`        mapresponse: '{ "data": resp }'
        findone:
          filter: '{ "_id": ObjectID(req._id) }'
`)
		builder.WriteString(fmt.Sprintf("      Create%v:\n", name))
		builder.WriteString(fmt.Sprintf("        database: %v\n", t.Database))
		builder.WriteString(fmt.Sprintf("        collection: %v\n", t.Collection))
		builder.WriteString("        insertone:\n")
		builder.WriteString("          document: '{ ")

		var builder2 strings.Builder
		kvs := ToSortedKV(t.Type.Fields)
		first := true
		for _, kv := range kvs {
			k := kv.Key
			if k == "_id" {
				continue
			}
			if !first {
				builder2.WriteString(", ")
			}
			first = false
			builder2.WriteString(fmt.Sprintf("?\"%v\": req.?data.?%v", k, k))
		}
		createAndUpdateInner := builder2.String()
		builder.WriteString(createAndUpdateInner)

		builder.WriteString("}'\n")
		builder.WriteString(fmt.Sprintf("      Update%v:\n", name))
		builder.WriteString(fmt.Sprintf("        database: %v\n", t.Database))
		builder.WriteString(fmt.Sprintf("        collection: %v\n", t.Collection))
		builder.WriteString("        replaceone:\n")
		builder.WriteString("          filter: '{ \"_id\": ObjectID(req.data._id) }'\n")
		builder.WriteString("          replacement: '{ ")
		builder.WriteString(createAndUpdateInner)
		builder.WriteString("}'\n")
		builder.WriteString(fmt.Sprintf("      Delete%v:\n", name))
		builder.WriteString(fmt.Sprintf("        database: %v\n", t.Database))
		builder.WriteString(fmt.Sprintf("        collection: %v\n", t.Collection))
		builder.WriteString(`        deleteone:
          filter: '{ "_id": ObjectID(req._id) }'
`)
	}

	return builder.String()
}

func (g *DapiCrud) GenerateServices() string {
	var builder strings.Builder
	builder.WriteString("service ")
	builder.WriteString(strcase.ToCamel(g.service))
	builder.WriteString(" {\n")

	for _, t := range g.ts {
		name := strcase.ToCamel(t.Name)
		for _, prefix := range []string{"List", "Get", "Create", "Update", "Delete"} {
			n := prefix + name
			builder.WriteString(fmt.Sprintf("  rpc %[1]v(%[1]vRequest) returns (%[1]vResponse);\n", n))
		}
	}

	builder.WriteString("}\n\n")

	for _, t := range g.ts {
		name := strcase.ToCamel(t.Name)

		builder.WriteString(fmt.Sprintf("message List%vRequest {\n  string cursor = 1;\n  int32 limit = 2;\n}\n\n", name))
		builder.WriteString(fmt.Sprintf("message List%vResponse {\n  repeated %v data = 1;\n  string next_cursor = 2;\n  int32 limit = 3;\n}\n\n", name, name))
		builder.WriteString(fmt.Sprintf("message Get%vRequest {\n  string _id = 1;\n}\n\n", name))
		builder.WriteString(fmt.Sprintf("message Get%vResponse {\n  %v data = 1;\n}\n\n", name, name))
		builder.WriteString(fmt.Sprintf("message Create%vRequest {\n  %v data = 1;\n}\n\n", name, name))
		builder.WriteString(fmt.Sprintf("message Create%vResponse {\n  string _id = 1;\n}\n\n", name))
		builder.WriteString(fmt.Sprintf("message Update%vRequest {\n  %v data = 1;\n}\n\n", name, name))
		builder.WriteString(fmt.Sprintf("message Update%vResponse {}\n\n", name))
		builder.WriteString(fmt.Sprintf("message Delete%vRequest {\n  string _id = 1;\n}\n\n", name))
		builder.WriteString(fmt.Sprintf("message Delete%vResponse {}\n\n", name))
	}

	return builder.String()
}
