module github.com/AwesomeLogos/bimi-logos

go 1.22.0

replace github.com/AwesomeLogos/bimi-explorer/generated => ./generated

require (
	github.com/AwesomeLogos/bimi-explorer/generated v0.0.0-00010101000000-000000000000
	github.com/aymerick/raymond v2.0.2+incompatible
	golang.org/x/net v0.10.0
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)

require (
	github.com/jackc/pgx/v5 v5.5.5
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
