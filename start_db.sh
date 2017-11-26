docker run \
    --name rocket-db \
    --env-file /go/src/github.com/ubclaunchpad/rocket/.env \
    -p 5432:5432 -d postgres:9.6