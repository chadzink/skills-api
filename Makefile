.PHONY: up down up-pg up-airflow

dev:
	docker-compose -f docker-compose.yaml up -d

down:
	docker-compose -f docker-compose.yaml down

