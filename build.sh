git pull
tag=$(git describe --abbrev=0 --tags)
if [[ $tag == "" ]]
then
    tag="build"
fi
docker build -t danbai225/go-rustdesk-server:"$tag" .
docker build -t danbai225/go-rustdesk-server:latest .
docker push danbai225/go-rustdesk-server:"$tag"
docker push danbai225/go-rustdesk-server:latest
