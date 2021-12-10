_DEPS= $(wildcard **/*.proto)
_OUT=gen

all: golang typescript

go: $(_DEPS)
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
	cd .. && protoc --go_out=.. protocol/**/*.proto protocol/**/**/*.proto


package.json:
node_modules: package.json
	yarn install


typescript: $(_DEPS) node_modules
	cd gen && mkdir -p js
	cd .. && protoc --plugin=./web/node_modules/.bin/protoc-gen-ts_proto --ts_proto_out=protocol/gen/js protocol/**/*.proto protocol/**/**/*.proto