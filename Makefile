.PHONY: dev prod

dev:
	docker-compose -f docker-compose.dev.yml up --build

prod:
	docker-compose up -d --build