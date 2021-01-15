# go-graphql-crud
This a simple Go CRUD example using GraphQL


Implement create, read, update and delete on Go.

To run the program:

1. Run the example: `go run main.go`

## Create

`http://localhost:8080/document?query=mutation+_{create(name:"Document Test",file:"2021-01-13-00-00-skfnsk82y4fbusnfkisn"){id,name,file}}`

## Read

* Get single document by id: `http://localhost:8080/document?query={document(id:1){name,file}}`
* Get document list: `http://localhost:8080/document?query={list{id,name,file}}`

## Update

`http://localhost:8080/document?query=mutation+_{update(id:1,price:3.95){id,name,file}}`

## Delete

`http://localhost:8080/document?query=mutation+_{delete(id:1){id,name,file}}`
