FROM golang:latest

ARG TESTNET

RUN mkdir /app
ADD . /app
WORKDIR /app
# set up directory for IPC VOLUME
RUN mkdir /ipc
VOLUME /ipc

RUN go mod tidy
RUN go build -o main ./cmd/main/
CMD ["/app/main", "--env", "TESTNET=$TESTNET"]