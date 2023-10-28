.PHONY: up down up-pg up-airflow

dev:
	docker-compose -f docker-compose.yaml up -d

down:
	docker-compose -f docker-compose.yaml down

clean:
	docker image rm skills-api-web

logs:
	docker-compose -f docker-compose.yaml logs -f