.PHONY: build

build: build_windows build_current
	echo "Done"

build_windows:
	GOOS=windows GOARCH=amd64 go build -o ./dist/

build_current:
	go build -o ./dist/
