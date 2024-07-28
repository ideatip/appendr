module go.ideatip.dev.appendr/example

go 1.22.0

replace (
	ideatip.dev.appendr => ../core
	ideatip.dev.appendr/nats => ../nats_appender
)

require (
	ideatip.dev.appendr v0.0.0-00010101000000-000000000000
	ideatip.dev.appendr/nats v0.0.0-00010101000000-000000000000
)

require (
	github.com/klauspost/compress v1.17.2 // indirect
	github.com/nats-io/nats.go v1.36.0 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.18.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
)
