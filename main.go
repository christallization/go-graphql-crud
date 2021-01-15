package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
	"github.com/graphql-go/graphql"
)

// Document contains infomation about one document

type Document struct {
	ID    int64   `json:"id"`
	Name  string  `json:"name,omitempty"`
	File  string  `json:"file,omitempty"`
}

var documents = []Document{
	{
		ID:    1,
		Name:  "Document one",
		File:  "a23hkjhl03209n2lh34sd009f92h3h4120098fwejk13h342h...",
	},
	{
		ID:    2,
		Name:  "Document 2",
		File:  "a23hkjhl03209n2lh34sd009f92h3h4120098fwejk13h342h...",
	},
	{
		ID:    3,
		Name:  "Document 3",
		File:  "a23hkjhl03209n2lh34sd009f92h3h4120098fwejk13h342h...",
	},
}

var documentType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Document",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"file": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			/* Get (read) single document by id
			   http://localhost:8080/document?query={document(id:1){name,file}}
			*/
			"document": &graphql.Field{
				Type:        documentType,
				Description: "Get document by id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"].(int)
					if ok {
						// Find document
						for _, document := range documents {
							if int(document.ID) == id {
								return document, nil
							}
						}
					}
					return nil, nil
				},
			},
			/* Get (read) documents list
			   http://localhost:8080/document?query={list{id,name,file}}
			*/
			"list": &graphql.Field{
				Type:        graphql.NewList(documentType),
				Description: "Get document list",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return documents, nil
				},
			},
		},
	}
)

var mutationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		/* Create new document item
		http://localhost:8080/document?query=mutation+_{create(name:"Test File",file:"test.pdf"){id,name,file}}
		*/
		"create": &graphql.Field{
			Type:        documentType,
			Description: "Create new document",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"file": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				rand.Seed(time.Now().UnixNano())
				document := Document{
					ID:    int64(rand.Intn(100000)), // generate random ID
					Name:  params.Args["name"].(string),
					File:  params.Args["file"].(string),
				}
				documents = append(documents, document)
				return document, nil
			},
		},
		/* Update document by id
		   http://localhost:8080/document?query=mutation+_{update(id:1,name:"test name"file:"test2.pdf"){id,name,file}}
		*/
		"update": &graphql.Field{
			Type:        documentType,
			Description: "Update document by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"file": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, _ := params.Args["id"].(int)
				name, nameOk := params.Args["name"].(string)
				file, fileOk := params.Args["file"].(string)
				document := Document{}
				for i, p := range documents {
					if int64(id) == p.ID {
						if nameOk {
							documents[i].Name = name
						}
						if fileOk {
							documents[i].File = file
						}
						document = documents[i]
						break
					}
				}
				return document, nil
			},
		},
		/* Delete document by id
		   http://localhost:8080/document?query=mutation+_{delete(id:1){id,name,file}}
		*/
		"delete": &graphql.Field{
			Type:        documentType,
			Description: "Delete document by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, _ := params.Args["id"].(int)
				document := Document{}
				for i, p := range documents {
					if int64(id) == p.ID {
						document = documents[i]
						// Remove from document list
						documents = append(documents[:i], documents[i+1:]...)
					}
				}
				return document, nil
			},
		},
	},
})

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	},
)

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("errors: %v", result.Errors)
	}
	return result
}

func main() {
	http.HandleFunc("/document", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)
	})
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}