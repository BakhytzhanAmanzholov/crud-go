run:
	echo "Run crud application"
	docker-compose -f docker-compose.yml up

build:
	echo "Build application"
	docker-compose -f docker-compose.yml build

swag:
	swag init -g server.go