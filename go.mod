module github.com/joostvdg/proglog

go 1.17

require github.com/gorilla/mux v1.8.0

require (
	github.com/stretchr/testify v1.7.0
	google.golang.org/protobuf v1.27.1
)

require github.com/tysonmote/gommap v0.0.1

require (
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013
	google.golang.org/grpc v1.33.2
)

require (
	github.com/casbin/casbin v1.9.1
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
)

require (
	go.opencensus.io v0.22.6
	go.uber.org/zap v1.10.0
)

require (
	github.com/hashicorp/serf v0.9.7
	github.com/travisjeffery/go-dynaport v1.0.0
)

require (
	github.com/hashicorp/raft v1.3.3
	github.com/hashicorp/raft-boltdb v0.0.0-00010101000000-000000000000
)

replace github.com/hashicorp/raft-boltdb => github.com/travisjeffery/raft-boltdb v1.0.0
