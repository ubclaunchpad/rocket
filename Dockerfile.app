################
# BINARY BUILD #
################

FROM golang:alpine AS build
ENV BUILD_HOME=/go/src/github.com/ubclaunchpad/rocket

# Mount source code.
ADD . ${BUILD_HOME}
WORKDIR ${BUILD_HOME}

# Install dependencies if not already available.
RUN if [ ! -d "vendor" ]; then \
    apk add --update --no-cache git; \
    go get -u github.com/golang/dep/cmd/dep; \
    dep ensure; \
    fi

# Build binary.
RUN go build -o /bin/rocket .

#################
#  IMAGE BUILD  #
#################

# Start from a fresh container
FROM alpine
LABEL maintainer "UBC Launchpad team@ubclaunchpad.com"

# Copy just the Rocket binary
COPY --from=build /bin/rocket /usr/local/bin

# Start Rocket by default
ENTRYPOINT [ "rocket" ]
