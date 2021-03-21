cd "$(dirname "$0")"
DIST="server-build"
rm -r ../$DIST/*
cp -r templates ../$DIST
cp -r static ../$DIST
go build -o ../server-build/main.bin