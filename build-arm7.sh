#docker build -t cgo-arm7 .
docker run --rm -v "$PWD":/workspace cgo-arm7 go build -o myapp-arm7 ./cmd
