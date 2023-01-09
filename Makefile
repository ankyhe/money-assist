all:
	go build -mod=vendor cmd/main.go
	mv main money-assist
test:
	go test -mod=vendor ./...
zip:
	rm -rf release.zip release
	mkdir release
	cp -f money-assist release
	cp -f config.txt release
	zip -r release.zip release
clean:
	rm -rf main money-assist release release.zip

