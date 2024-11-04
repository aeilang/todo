POSTGRES_URL="postgres://user:password@localhost:5432/mydb?sslmode=disable"

migrate_init:
	migrate create -seq -ext sql -dir=./migration ${name}

migrate_up: migrate_drop
	migrate -database=${POSTGRES_URL} -path="./migration" up

migrate_drop:
	migrate -database=${POSTGRES_URL} -path="./migration" drop -f

postgres_up:
	docker run -d \
	--name postgresdb \
	-e POSTGRES_USER=user \
	-e POSTGRES_PASSWORD=password \
	-e POSTGRES_DB=mydb \
	-p 5432:5432 \
	postgres:alpine