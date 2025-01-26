docker-run:
	docker compose up -f build/docker-compsoe.yaml

build-image:
	docker build -f build/Dockerfile

run:
	go run ./cmd