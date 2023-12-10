proto-generate:
	@protoc --go_out=. ./api/protobuf/model.proto

test-app:
	@docker compose -p fetch-system-iac -f ./build/docker-compose.test.yml up --attach server --build --abort-on-container-exit

build-app:
	docker build -t fetch-system-initializer:latest -f ./build/Dockerfile .
