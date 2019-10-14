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

docker build -t get-container-id:1.0.0 .

docker run -it --rm get-container-id:1.0.0 -httpPort 6666

rm -f $BUILD_PATH/build/bin/main
