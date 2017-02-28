run:
	docker-compose down
	docker-compose build
	docker-compose up

build:
	docker build -t teachbot-api -f conf/Dockerfile-API .

dev:
	docker run -d -p 5432:5432 --name db -e POSTGRES_PASSWORD=teachpass -e POSTGRES_DB=usersdev -e POSTGRES_USER=teach postgres
	docker run -p 8080:8080 -v $(PWD):/go/src/github.com/IrfanFaizullabhoy/teacher -v mounted-volume:/mounted-volume --name teachbot-api -it --link db:db --env-file .env teachbot-api

kill: 
	docker kill db | true
	docker kill teachbot-api | true
	docker rm db
	docker rm -v teachbot-api

enter:
	docker exec -it teachbot-api /bin/sh