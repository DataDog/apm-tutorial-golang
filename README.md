# apm-tutorial-golang
Tutorial for golang application used by APM Ecosystems team for agent setup in popular configurations

# Installation

Navigate to calendar/ and run `go run cmd/main.go` to run calendar application

Navigate to notes/ and run `go run cmd/main.go` to run notes application

# Example requests:

`curl localhost:8080/notes`

[]

`curl -X POST 'localhost:8080/notes?desc=hello'`

{"id":1,"description":"hello"}

`curl localhost:8080/notes/1`

{"id":1,"description":"hello"}

`curl localhost:8080/notes`

[{"id":1,"description":"hello"}]

`curl -X PUT ‘localhost:8080/notes/1/desc=UpdatedNote’`

{"id":1,"description":"UpdatedNote"}

`curl DELETE ‘localhost:8080/notes/1’`

Deleted

`curl -X POST 'localhost:8080/notes?desc=hello_again&add_date=y'`

{"id":2,"description":"hello_again with date 2022-11-06"}
