#docker build -t cgo-arm7 .

#方式1 纯go CGO_ENABLED=0
#docker run --rm -v "$PWD":/workspace \
#  -e CGO_ENABLED=0 \
#  -e GOOS=linux \
#  -e GOARCH=arm \
#  -e GOARM=7 \
#  golang:1.24-bullseye \
#  go build -ldflags '-w -s' -o app-arm7 ./cmd

#方式2 CGO_ENABLED=1
docker run --rm -v "$PWD":/workspace \
  -e CGO_ENABLED=1 \
  -e CC=arm-linux-gnueabihf-gcc \
  -e CGO_LDFLAGS="-static -s -w" \
  cgo-arm7 \
  go build -a -installsuffix cgo -ldflags '-w -s' -o app-arm7 ./cmd


