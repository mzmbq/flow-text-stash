.PHONY: build
build:
	@GOOS=windows GOARCH=amd64 go build -o ./bin/textstash/ts.exe -v ./main.go
	@cp ./assets/* ./bin/textstash
