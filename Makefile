.PHONY: dev prod

dev:
	docker-compose -f docker-compose.dev.yml up --build

prod:
	docker-compose up -d --build

swagger:
	swag init --parseDependency --parseInternal -g cmd/api/main.go