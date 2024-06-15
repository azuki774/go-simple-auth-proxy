FROM golang:1.21.0 as builder
LABEL description="go-simple-auth-proxy"

COPY main.go /src/main.go
COPY Makefile /src/Makefile
COPY cmd/ /src/cmd/
COPY internal/ /src/internal/
COPY go.mod /src/go.mod
WORKDIR /src
RUN go mod tidy && mkdir bin && make bin

FROM gcr.io/distroless/static-debian12
COPY --from=builder /src/bin/go-simple-auth-proxy /usr/local/bin/go-simple-auth-proxy
COPY build/default.toml /etc/default.toml
# Overwritting by volume in order to use your settings.
ENTRYPOINT ["/usr/local/bin/go-simple-auth-proxy", "start"]
CMD ["-c", "/etc/default.toml"]
