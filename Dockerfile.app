# Start from the base Go image
FROM golang

# Set /go/src/github.com/ubclaunchpad/rocket as the CWD
WORKDIR /go/src/github.com/ubclaunchpad/rocket

# Copy package source files to container
ADD . .

# Download and install dependency manager
RUN go get github.com/Masterminds/glide
RUN go install github.com/Masterminds/glide

# Install dependencies
RUN glide install

# Build Rocket
RUN go install github.com/ubclaunchpad/rocket

# Start Rocket
ENTRYPOINT [ "rocket" ]
