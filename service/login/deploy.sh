# kill the old one container
workdir=/home/trance/Documents/docker/lemonstreet/service
service=login
target="$workdir"/"$service"
echo "starting build service..."
cd cmd && CGO_ENABLED=0 go build && cd ..
echo "build success, depolying..."
cp -r cmd/cmd configs -t "$target"
echo "execute run.sh"
cd "$target" && ./run.sh
 
