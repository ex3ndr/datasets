set -e
rm -fr ./out/
mkdir -p ./out/
cd ./sources/

# Building
GOOS=linux GOARCH=amd64 go build -o ../out/datasets-linux-amd64 .
GOOS=darwin GOARCH=amd64  go build -o ../out/datasets-darwin-amd64 .
GOOS=darwin GOARCH=arm64  go build -o ../out/datasets-darwin-arm64 .
GOOS=windows GOARCH=amd64 go build -o ../out/datasets-windows.exe .
