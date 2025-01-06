IMAGE_NAME=pastebin_app
DOCKER_USERNAME=judyas
VERSION=latest

run:
	go run ./cmd/main.go

docker-build:
	docker build -t $(DOCKER_USERNAME)/$(IMAGE_NAME):$(VERSION) .

docker-run:
	docker run -p 8080:8080 --env-file .env $(DOCKER_USERNAME)/$(IMAGE_NAME):$(VERSION)

docker-push:
	docker push $(DOCKER_USERNAME)/$(IMAGE_NAME):$(VERSION)

up:
	docker-compose up -d

down:
	docker-compose down