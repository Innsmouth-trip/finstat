build:
	docker build -f ./deployments/Dockerfile -t finstat-app:0.0.1 .

run: build
	docker-compose -f ./deployments/docker-compose.yaml -p finstat-app --env-file .env up