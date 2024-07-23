PROJECT_NAME= fetch-system-iac
test-package:
	docker compose -p $(PROJECT_NAME) -f ./deploy/docker-compose.test.yml build
	docker compose -p $(PROJECT_NAME) -f ./deploy/docker-compose.test.yml \
	run server go test $(TEST_DIR)
	docker compose -p $(PROJECT_NAME) -f ./deploy/docker-compose.test.yml down --remove-orphans

test-app:
	docker compose -p $(PROJECT_NAME) -f ./deploy/docker-compose.test.yml up --attach server --build --abort-on-container-exit
	docker compose -p $(PROJECT_NAME) -f ./deploy/docker-compose.test.yml down --remove-orphans

build-db-initer-app:
	docker build -t fetch-system-db-initer:latest -f ./deploy/Dockerfile.dbiniter .

build-preparer-app:
	docker build -t fetch-system-preparer:latest -f ./deploy/Dockerfile.preparer .

generate-proto:
	@protocf --go_out=. ./api/protobuf/model.proto

generate-sqlc:
	@sqlc generate -f ./api/sql/sqlc.yml

generate-wire:
	@wire cmd/wire/wire_setup.go