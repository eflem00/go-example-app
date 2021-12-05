build:
	go build -o bin/main main.go

run:
	go run main.go

docker-build:
	docker build --tag go-example-app .

docker-run:
	docker run -d -p 8081:8081 go-example-app