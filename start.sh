docker stop rocket-app

docker build \
    -t rocket-app:latest \
    /go/src/github.com/ubclaunchpad/rocket

docker run \
    --name rocket-app \
    --env-file /go/src/github.com/ubclaunchpad/rocket/.app.env \
    --network rocket-net \
    -p 80:80 \
    -p 5432:5432 \
    -d \
    rocket-app

