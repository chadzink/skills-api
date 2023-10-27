.PHONY: up down up-pg up-airflow

dev:
	docker-compose -f docker-compose.yaml up -d

down:
	docker-compose -f docker-compose.yaml down

rm-image:
	docker image rm skills-api-web

logs:
	docker-compose -f docker-compose.yaml logs -f