build:
	go build -o bin/main main.go

run:
	go run main.go

docker-build:
	docker build --tag go-example-app .

docker-run:
	docker run -d -p 8080:8080 go-example-app