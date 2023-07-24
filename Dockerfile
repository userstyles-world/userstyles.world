FROM docker.io/library/golang:1.20.6-alpine3.18
WORKDIR /src
RUN apk add --no-cache build-base sassc
COPY go.sum go.mod ./
RUN  go mod download && go mod tidy
COPY . /src
RUN set -x \
 && mkdir -p web/static/css \
 && for f in web/scss/*.scss; do sassc --style nested --sourcemap=inline -l "$f" "web/static/css/$(basename $f .scss).css";done \
 && go build -v -o bin/userstyles-fonts cmd/userstyles-fonts/main.go \
 && go build -v -o bin/userstyles-ts cmd/userstyles-ts/main.go \
 && go build -v -o bin/userstyles -tags "fts5" cmd/userstyles-world/main.go

FROM docker.io/library/alpine:3.18
RUN apk add --no-cache vips-tools
WORKDIR /data
ENV DATA_DIR=/data
ENV DB=userstyles.db
COPY --from=0 /src/bin/ /usr/local/bin/
COPY docker-entrypoint.sh /usr/local/bin/
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["userstyles"]
