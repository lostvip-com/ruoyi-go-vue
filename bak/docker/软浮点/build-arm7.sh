docker build -t cgo-arm .

#方式1 纯go CGO_ENABLED=0
#docker run --rm -v "$PWD":/workspace \
#  -e CGO_ENABLED=0 \
#  -e GOOS=linux \
#  -e GOARCH=arm \
#  -e GOARM=7 \
#  golang:1.24-bullseye \
#  go build -ldflags '-w -s' -o app-arm7 ./cmd

#方式2 CGO_ENABLED=1 软浮点
docker run --rm -v "$PWD":/workspace \
  cgo-arm \
  go build -ldflags '-w -s' -o app-arm7 ./cmd



