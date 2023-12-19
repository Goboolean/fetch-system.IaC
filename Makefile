proto-generate:
	@protoc --go_out=. ./api/protobuf/model.proto

test-app:
	docker compose -p fetch-system-iac -f ./deploy/docker-compose.test.yml up --attach server --build --abort-on-container-exit
	docker compose -p fetch-system-iac -f ./deploy/docker-compose.test.yml down --remove-orphans

build-app:
	docker build -t fetch-system-initializer:latest -f ./deploy/Dockerfile.job .

generate-sqlc:
	sqlc generate -f ./api/sql/sqlc.yml

generate-wire:
	wire cmd/wire/wire_setup.go