### Compile stage
FROM golang:1.21-alpine3.18 AS build-env
RUN apk add --no-cache build-base make git

ADD . /dockerbuild
WORKDIR /dockerbuild

# timezone data for alpine builds
ENV GOEXPERIMENT=loopvar
RUN GIT_VERSION=$(git describe --tags --long --always) && \
    go build -tags timetzdata -o /zookeeper ./cmd/zookeeper

### Run stage
FROM alpine:3.18

RUN apk add --no-cache --update dumb-init ca-certificates
ENTRYPOINT ["dumb-init", "--"]

WORKDIR /
RUN mkdir -p data/zookeeper
COPY --from=build-env /zookeeper /

# small things to make golang binaries work well under alpine
ENV GODEBUG=netdns=go
ENV TZ=Etc/UTC

EXPOSE 9092 9192

CMD ["/zookeeper", "run"]
