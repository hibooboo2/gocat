# build stage
FROM golang:alpine AS build-env
ENV D=/go/src/github.com/hibooboo2/gocat/containerips
# If dep ever gets decent enough to use, try `dep ensure --vendor-only` from here: https://medium.com/travis-on-docker/triple-stage-docker-builds-with-go-and-angular-1b7d2006cb88
ADD . $D
RUN cd $D && go build -o containerips && cp containerips /tmp/

# final stage
FROM golang:alpine
EXPOSE 42335
RUN apk add --no-cache ca-certificates curl
WORKDIR /app
COPY --from=build-env /tmp/containerips /app/containerips
ENTRYPOINT ["./containerips"]
