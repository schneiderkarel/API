up:
	docker network create local-network
	docker-compose up -d --remove-orphans postgres

	./ping_postgres.sh postgres 60

	docker-compose up -d --build --remove-orphans migrate-up

	docker-compose up -d --remove-orphans --build user-service

down:
	docker-compose down
	docker network rm local-network

migrate-up:
	docker-compose up -d --remove-orphans migrate-up
