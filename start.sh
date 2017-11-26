docker build \
    -t rocket-app:latest \
    /go/src/github.com/ubclaunchpad/rocket

docker stop rocket-app

docker rm rocket-app

docker run \
    --name rocket-app \
    --env-file /go/src/github.com/ubclaunchpad/rocket/.app.env \
    --network rocket-net \
    -p 80:80 \
    -d \
    rocket-app
