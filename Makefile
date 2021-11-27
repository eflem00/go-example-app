build:
	go build -o bin/main main.go

run:
	go run main.go

docker-build:
	docker build --build-arg AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} --build-arg AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} --build-arg AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION} --tag go-example-app .

docker-run:
	docker run -d -p 8080:8080 go-example-app