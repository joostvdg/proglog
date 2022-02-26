FROM golang:1.17 AS build
WORKDIR /go/src/proglog
ARG TARGETARCH
ARG TARGETOS
ENV GRPC_HEALTH_PROBE_VERSION=v0.3.2
ENV HEALTH_PROBE_URL="https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-${TARGETOS}-${TARGETARCH}"
RUN wget "$HEALTH_PROBE_URL"
RUN mv grpc_health_probe-${TARGETOS}-${TARGETARCH} /go/bin/grpc_health_probe
RUN chmod +x /go/bin/grpc_health_probe
WORKDIR /go/src/proglog
COPY go.* ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 go build -o /go/bin/proglog ./cmd/proglog

FROM alpine:3
ARG TARGETARCH
RUN apk --no-cache add ca-certificates
ENTRYPOINT ["/bin/proglog"]
COPY --from=build /go/bin/grpc_health_probe /bin/grpc_health_probe
COPY --from=build /go/bin/proglog /bin/proglog