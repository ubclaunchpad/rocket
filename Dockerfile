# Start from the base Go image
FROM golang

# Set /go/src/github.com/ubclaunchpad/rocket as the CWD
WORKDIR /go/src/github.com/ubclaunchpad/rocket

# Install Postgres client to check when the DB is ready for use
RUN apt-get update
RUN apt-get install -f -y postgresql-client

# Copy package source files to container
ADD . .

# Download and install dependency manager
RUN go get github.com/Masterminds/glide
RUN go install github.com/Masterminds/glide

# Install dependencies
RUN glide install

# Build Rocket
RUN go install github.com/ubclaunchpad/rocket
