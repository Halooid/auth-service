module github.com/halooid/backend/auth-service

go 1.25.0

require (
	github.com/google/uuid v1.6.0
	github.com/halooid/backend/go-shared v0.0.0-00010101000000-000000000000
	github.com/jackc/pgx/v5 v5.9.2
	google.golang.org/grpc v1.64.0
	google.golang.org/protobuf v1.34.1
)

require (
	github.com/MicahParks/keyfunc/v2 v2.1.0 // indirect
	github.com/Nerzal/gocloak/v13 v13.9.0 // indirect
	github.com/go-resty/resty/v2 v2.7.0 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/segmentio/ksuid v1.0.4 // indirect
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240318140521-94a12d6c2237 // indirect
)

replace github.com/halooid/backend/go-shared => ../go-shared
