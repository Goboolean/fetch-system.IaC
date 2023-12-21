test-app:
	docker compose -p fetch-system-iac -f ./deploy/docker-compose.test.yml up --attach server --build --abort-on-container-exit
	docker compose -p fetch-system-iac -f ./deploy/docker-compose.test.yml down --remove-orphans

build-retriever-app:
	docker build -t fetch-system-retriever:latest -f ./deploy/Dockerfile.prepare .

build-preparer-app:
	docker build -t fetch-system-preparer:latest -f ./deploy/Dockerfile.retrieve .

generate-proto:
	@protoc --go_out=. ./api/protobuf/model.proto

generate-sqlc:
	@sqlc generate -f ./api/sql/sqlc.yml

generate-wire:
	@wire cmd/wire/wire_setup.go