set -e
rm -fr ./out/
mkdir -p ./out/
cd ./sources/

# Building
GOOS=linux GOARCH=amd64 go build -o ../out/linux-x64/ .
GOOS=darwin GOARCH=amd64  go build -o ../out/darwin-x64/ .
GOOS=darwin GOARCH=arm64  go build -o ../out/darwin-arm64/ .
GOOS=windows GOARCH=amd64 go build -o ../out/windows-x64/ .
