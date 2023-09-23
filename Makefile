proto-generate:
	protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
    ./api/model/data.proto