test-app:
	docker compose -p fetch-system-iac -f ./deploy/docker-compose.test.yml up --attach server --build --abort-on-container-exit
	docker compose -p fetch-system-iac -f ./deploy/docker-compose.test.yml down --remove-orphans

build-db-initer-app:
	docker build -t registry.mulmuri.dev/fetch-system-db-initer:latest -f ./deploy/Dockerfile.dbiniter .
	docker push registry.mulmuri.dev/fetch-system-db-initer:latest
build-preparer-app:
	docker build -t registry.mulmuri.dev/fetch-system-preparer:latest -f ./deploy/Dockerfile.preparer .
	docker push registry.mulmuri.dev/fetch-system-preparer:latest
	helm upgrade fetch-system ~/lab -n goboolean

generate-proto:
	@protocf --go_out=. ./api/protobuf/model.proto

generate-sqlc:
	@sqlc generate -f ./api/sql/sqlc.yml

generate-wire:
	@wire cmd/wire/wire_setup.go