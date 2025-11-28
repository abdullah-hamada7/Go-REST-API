# Commands
```bash 

curl localhost:8080/books --include --header  "Content-Type: application/json" --request "GET"

curl localhost:8080/books --include --header  "Content-Type: application/json" -d @body.json --request "POST"

curl localhost:8080/books/1 --include --header  "Content-Type: application/json" --request "GET"

curl localhost:8080/books/1 --include --header  "Content-Type: application/json" -d @body.json --request "PUT"

curl localhost:8080/books/1 --include --header  "Content-Type: application/json" -d @body.json --request "PATCH"

curl localhost:8080/books/1 --include --header  "Content-Type: application/json" --request "DELETE"
```