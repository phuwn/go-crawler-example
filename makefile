.PHONY: remove-infras init run
POSTGRES_CONTAINER?=crawlexample_db

remove-infras:
	docker-compose stop; docker-compose  rm -f

init: remove-infras
	@docker-compose  up -d 
	@echo "Waiting for database connection..."
	@while ! docker exec $(POSTGRES_CONTAINER) pg_isready -h localhost -p 5432 > /dev/null; do \
		sleep 1; \
	done
	@sql-migrate up -config=dbconfig.yml -env="local"

run:
	@go run main.go