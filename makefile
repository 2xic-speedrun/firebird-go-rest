
build_windows:
	GOOS=windows GOARCH=amd64 go build -o ./dist/

build:
	go build -o ./dist/
	
