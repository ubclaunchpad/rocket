# Start from the base Go image
FROM golang

# Copy package source files to container
ADD . /go/src/github.com/ubclaunchpad/rocket

# Download and install dependency manager
RUN go get github.com/Masterminds/glide
RUN go install github.com/Masterminds/glide

# Install dependencies
RUN cd /go/src/github.com/ubclaunchpad/rocket && \
    glide install

# Build Rocket
RUN go install github.com/ubclaunchpad/rocket

# Run
ENTRYPOINT [ "/go/bin/rocket" ]
