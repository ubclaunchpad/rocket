docker stop rocket-db

docker rm rocket-db

docker run \
    --name rocket-db \
    --env-file /go/src/github.com/ubclaunchpad/rocket/.db.env \
    --network rocket-net \
    -p 5432:5432 \
    -v pgdata:/var/lib/postgresql/data \
    -d \
    postgres:9.6
