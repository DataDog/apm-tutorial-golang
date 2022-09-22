# apm-tutorial-golang
Tutorial for golang application used by APM Ecosystems team for agent setup in popular configurations

# Installation

Inside the cmd folder there are two directories for notes and calendar.

A makefile is attached for convenience. Use `make` to build the programs and `make run` to run both the notes and calendar applications. Use `make exit` to kill the running applications in the background.

Otherwise, build and run the main.go file inside each directory to run the corresponding application

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
