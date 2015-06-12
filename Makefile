TAG?=latest

all: moviedb-frontend
	docker build -t jamesclonk/moviedb-frontend:${TAG} .
	rm moviedb-frontend

moviedb-frontend: main.go
	GOARCH=amd64 GOOS=linux go build -o moviedb-frontend

test:
	GOARCH=amd64 GOOS=linux go test -v ./...
