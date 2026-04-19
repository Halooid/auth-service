module github.com/halooid/backend/auth-service

go 1.25.0

require (
	github.com/google/uuid v1.6.0
	github.com/lib/pq v1.12.3
	google.golang.org/grpc v1.64.0
	google.golang.org/protobuf v1.34.1
)

require (
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240318140521-94a12d6c2237 // indirect
)

replace github.com/halooid/backend/go-shared => ../go-shared
