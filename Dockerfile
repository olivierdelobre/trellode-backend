# STEP 1: build
FROM golang:1.21 as builder

# setup the working directory
WORKDIR /api

# install dependencies
COPY go.*  /api/
RUN go mod download

# add source code
COPY cmd cmd/
COPY internal internal/
COPY docs docs/

# build the source
RUN CGO_ENABLED=0 GOOS=linux GOFLAGS="-ldflags=-s -ldflags=-w" go build -o server ./cmd/api/

# STEP 2: app
FROM golang:1.21-bookworm

ENV TZ=Europe/Zurich

# add ca-certificates in case you need them
RUN apt-get update && apt-get install ca-certificates jq -y && rm -rf /var/cache/apk/*

RUN groupadd trellode && \
    useradd -r --uid 1001 -g trellode trellode

# set working directory
RUN mkdir -p /home/trellode/data
RUN echo "test" > /home/trellode/data/test.out
WORKDIR /home/trellode

# copy the binary from builder
COPY --from=builder /api/server /home/trellode/server
COPY assets/i18n/ /home/trellode/i18n
COPY docs docs/

# Ownership so that these folders can be written when running in K8S
RUN chgrp -R 0 /home/trellode && chmod -R g=u /home/trellode

USER 1001
CMD ["/home/trellode/server"]
