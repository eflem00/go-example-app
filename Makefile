build:
	go build -o bin/main main.go

run:
	go run main.go

docker-build:
	docker build --tag go-example-app .

docker-run:
	docker run -d -p 8080:8080 go-example-app

docker-pg:
	docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=postgres postgres