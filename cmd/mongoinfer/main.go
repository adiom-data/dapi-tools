package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/adiom-data/dapi-tools/tool/mongodb/mongoinfer"
	"github.com/iancoleman/strcase"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type stringList []string

func (s *stringList) String() string {
	return strings.Join(*s, ",")
}

func (s *stringList) Set(value string) error {
	*s = append(*s, strings.Split(value, ",")...)
	return nil
}

func main() {
	var namespaces stringList
	var databases stringList
	flag.Var(&databases, "database", "Database to generate for (can set multiple times or comma separated)")
	flag.Var(&namespaces, "namespace", "Fully qualified namespace to generate for (can set multiple times or comma separated)")
	samples := flag.Int("samples", 1000, "Number of samples to use")
	url := flag.String("url", "mongodb://localhost:27017", "The mongodb url")
	pkg := flag.String("package", "mypkg", "The package name")
	service := flag.String("service", "MyService", "Service name to generate")
	protoFile := flag.String("proto-file", "", "Protobuf source file to generate")
	dapiConfigFile := flag.String("dapi-config-file", "", "Partial Dapi config.yml file to generate")
	server := flag.Bool("server", false, "Include default server config")
	flag.Parse()

	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(*url))
	if err != nil {
		panic(err)
	}

	allNamespaces := map[string]struct{}{}
	for _, n := range namespaces {
		allNamespaces[n] = struct{}{}
	}
	if len(databases) == 0 && len(namespaces) == 0 {
		dbs, err := client.ListDatabaseNames(ctx, bson.M{})
		if err != nil {
			panic(err)
		}
		for _, db := range dbs {
			if db == "local" || db == "admin" || db == "config" {
				continue
			}
			databases = append(databases, db)
		}
	}
	for _, d := range databases {
		db := client.Database(d)
		cols, err := db.ListCollectionNames(ctx, bson.M{})
		if err != nil {
			panic(err)
		}
		for _, col := range cols {
			n := fmt.Sprintf("%v.%v", d, col)
			if _, ok := allNamespaces[n]; !ok {
				allNamespaces[n] = struct{}{}
				namespaces = append(namespaces, n)
			}
		}
	}

	var builder strings.Builder
	builder.WriteString(mongoinfer.Header(*pkg))
	var namedTypes []mongoinfer.NamedType

	for _, namespace := range namespaces {
		database, collection, ok := strings.Cut(namespace, ".")
		if !ok {
			slog.Error("Could not parse namespace", "namespace", namespace)
			continue
		}
		col := client.Database(database).Collection(collection)

		bsonType, err := mongoinfer.BsonTypeFromSamples(ctx, col, *samples)
		if err != nil {
			panic(err)
		}

		if err := mongoinfer.BsonTypeToProto(&builder, strcase.ToCamel(namespace), bsonType); err != nil {
			panic(err)
		}

		namedTypes = append(namedTypes, mongoinfer.NamedType{
			Name:       namespace,
			Type:       bsonType,
			Database:   database,
			Collection: collection,
		})
	}

	dapiCrud := mongoinfer.NewDapiCrud(*pkg, *service, namedTypes)

	builder.WriteString(dapiCrud.GenerateServices())
	protoContents := builder.String()

	dapiCfg := mongoinfer.DefaultInterceptor(*url) + dapiCrud.GenerateDapiCfg()
	if *server {
		dapiCfg = mongoinfer.DefaultServer + dapiCfg
	}

	if *dapiConfigFile != "" {
		_ = os.MkdirAll(filepath.Dir(*dapiConfigFile), 0755)
		if err := os.WriteFile(*dapiConfigFile, []byte(dapiCfg), 0644); err != nil {
			panic(err)
		}
	}

	if *protoFile != "" {
		_ = os.MkdirAll(filepath.Dir(*protoFile), 0755)
		if err := os.WriteFile(*protoFile, []byte(protoContents), 0644); err != nil {
			panic(err)
		}
	}

	if *dapiConfigFile == "" && *protoFile == "" {
		fmt.Println(dapiCfg)
		fmt.Println("--- --- ---")
		fmt.Println(protoContents)
	}
}
