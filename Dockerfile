# STEP 1 build the app
FROM golang:alpine as builder
WORKDIR $GOPATH/src/github.com/matrjoschka/matrjoschka-server/

RUN mkdir -p /opt/matrjoschka
RUN apk update && apk add git

COPY . ./
RUN go get -u github.com/Masterminds/glide
RUN glide install

ENV GIN_MODE=release

RUN go build -o /opt/matrjoschka/server .

# STEP 2 build a small image
FROM alpine
COPY --from=builder /opt/matrjoschka /opt/matrjoschka
WORKDIR /opt/matrjoschka

RUN apk update && apk add ca-certificates

ENV GIN_MODE=release
EXPOSE 9287

ENTRYPOINT ["./server"]