# Check Args
if [ $# -lt 1 ]; then
   echo "Please Input Build Args : sh build.sh {{tag}}"
   exit 1
fi

BUILD_TAG=$1
echo BUild Tag: $BUILD_TAG

BUILD_PATH=$(pwd)
echo Build Path: $BUILD_PATH

GOLANG_VERSION="latest"
echo Build Golang Version : $GOLANG_VERSION

docker run -it --rm \
	-v $BUILD_PATH:/go/src/github.com/ming-go/lab/get-container-id \
	-w /go/src/github.com/ming-go/lab/get-container-id \
	-e CGO_ENABLED=0 \
	-e GOOS=linux \
	-e GOARCH=amd64 \
	golang:$GOLANG_VERSION \
	go build -v -a -installsuffix cgo \
	-o build/bin/main .

cd ""$BUILD_PATH""

docker build -t iwdmb/get-container-id:$BUILD_TAG .

docker run -it --rm iwdmb/get-container-id:$BUILD_TAG -httpPort 8080

rm -f $BUILD_PATH/build/bin/main
