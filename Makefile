SERVER_PATH="cmd/server/main.go"
PROTO_RELATIVE_PATH="api"

local:
	ENV=local go run cmd/server/main.go

development:
	ENV=development go run cmd/server/main.go

staging:
	ENV=staging go run cmd/server/main.go

production:
	ENV=production go run cmd/server/main.go