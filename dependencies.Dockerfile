FROM golang:1.13.7
# Add the module files and download dependencies.
ENV GO111MODULE=on
COPY ./go.mod /go/src/github.com/elections/go.mod
COPY ./go.sum /go/src/github.com/elections/go.sum
WORKDIR /go/src/github.com/elections
RUN go mod download
# Add the shared packages.
COPY ./cache /go/src/github.com/elections/source/cache
COPY ./databuilder /go/src/github.com/elections/source/databuilder
COPY ./donations /go/src/github.com/elections/source/donations
COPY ./dynamo /go/src/github.com/elections/source/dynamo
COPY ./idhash /go/src/github.com/elections/source/idhash
COPY ./indexing /go/src/github.com/elections/source/indexing
COPY ./persist /go/src/github.com/elections/source/persist
COPY ./protobuf /go/src/github.com/elections/source/protobuf
COPY ./server /go/src/github.com/elections/source/server
COPY ./svc/proto /go/src/github.com/elections/source/svc/proto
COPY ./util /go/src/github.com/elections/source/util

