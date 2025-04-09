VERSION=0.1.0
.PHONY: build
build:
	@GOOS=windows GOARCH=amd64 go build -o ./bin/TextStash-$(VERSION)/ts.exe -v ./main.go
	@cp ./assets/* ./bin/TextStash-$(VERSION)

.PHONY: bundle
bundle:
	@zip -r -j ./bin/TextStash-$(VERSION).zip ./bin/TextStash-$(VERSION)