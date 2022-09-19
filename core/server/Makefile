

gomod: go.mod
	go get -d ./...;
	bazel run //server:gazelle -- update-repos -from_file=server/go.mod