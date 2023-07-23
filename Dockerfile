FROM docker.io/library/golang:1.20.6-alpine3.18
WORKDIR /src
RUN apk add --no-cache build-base sassc vips-tools
COPY go.sum go.mod ./
RUN  go mod download && go mod tidy
COPY . /src
RUN go build -v -o userstyles -tags "fts5" cmd/userstyles-world/main.go

FROM docker.io/library/alpine:3.18
RUN apk add --no-cache vips-tools
WORKDIR /data
ENV DATA_DIR=/data
ENV DB=userstyles.db
COPY --from=0 /src/userstyles /usr/local/bin/userstyles
COPY docker-entrypoint.sh /usr/local/bin/
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["userstyles"]
