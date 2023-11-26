proto-generate:
	protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
    ./api/model/data.proto

make test-app:
	@docker compose -f ./build/docker-compose.test.yml up --build --abort-on-container-exit
	@docker compose -f ./build/docker-compose.test.yml down