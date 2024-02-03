git pull
tag=$(git describe --abbrev=0 --tags)
if [[ $tag == "" ]]
then
    tag="build"
fi
docker buildx build -t danbai225/go-rustdesk-server:"$tag" --platform=linux/amd64 .
docker buildx build -t danbai225/go-rustdesk-server:latest --platform=linux/amd64 .
docker push danbai225/go-rustdesk-server:"$tag"
docker push danbai225/go-rustdesk-server:latest
