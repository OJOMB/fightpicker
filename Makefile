gen-dtos:
	oapi-codegen -generate types -o internal/server/dtos/dtos.go -package dtos ./api/swagger.yaml

gen-sql:
	sqlc generate

test:
	go test ./... -v -covmode=atomic