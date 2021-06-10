FROM node:14-alpine3.13 AS frontend
COPY src/frontend /go/src/saltgram/frontend
WORKDIR /go/src/saltgram/frontend
RUN npm install -g @vue/cli
RUN npm install
RUN npm run build

FROM golang:1.16-alpine3.13 AS dependencies
COPY --from=frontend /go/src/saltgram/frontend /go/src/saltgram/frontend
COPY src/go.mod /go/src/saltgram
COPY src/go.sum /go/src/saltgram
COPY src/data /go/src/saltgram/data
COPY src/internal /go/src/saltgram/internal
COPY src/protos /go/src/saltgram/protos
COPY src/pki /go/src/saltgram/pki
COPY ./wait-for-postgres.sh /go/src/saltgram
WORKDIR /go/src/saltgram
RUN go get -d -v ./...
