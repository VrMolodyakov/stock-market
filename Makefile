test:
	go test ./... -cover
start:
	docker-compose up -d

build:
	docker-compose build

stop:
	docker-compose down