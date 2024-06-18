# syntax=docker/dockerfile:1

FROM golang:1.22

WORKDIR /

COPY . ./
RUN go mod download

RUN ./dev.sh --test

RUN CGO_ENABLED=0 GOOS=linux go build -o /emrs

# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can (optionally) document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
EXPOSE 8080

# Run

ENTRYPOINT ["/emrs", "-addr", "0.0.0.0:8080"]
