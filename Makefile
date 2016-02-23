test:
	GO15VENDOREXPERIMENT=1 go test . -cover

build:
	GO15VENDOREXPERIMENT=1 go build .

release:
	@env GO15VENDOREXPERIMENT=1 GOOS=linux GOARCH=amd64 go build -ldflags="-X main.appVersion=${DOLAN_VERSION}" -o dolan
	@zip dolan_${DOLAN_VERSION}_linux_amd64.zip dolan
	@rm dolan

	@env GO15VENDOREXPERIMENT=1 GOOS=linux GOARCH=386 go build -ldflags="-X main.appVersion=${DOLAN_VERSION}" -o dolan
	@zip dolan_${DOLAN_VERSION}_linux_386.zip dolan
	@rm dolan

	@env GO15VENDOREXPERIMENT=1 GOOS=freebsd GOARCH=amd64 go build -ldflags="-X main.appVersion=${DOLAN_VERSION}" -o dolan
	@zip dolan_${DOLAN_VERSION}_freebsd_amd64.zip dolan
	@rm dolan

	@env GO15VENDOREXPERIMENT=1 GOOS=freebsd GOARCH=386 go build -ldflags="-X main.appVersion=${DOLAN_VERSION}" -o dolan
	@zip dolan_${DOLAN_VERSION}_freebsd_386.zip dolan
	@rm dolan

clean:
	@rm -f dolan
	@rm -rf dolan_*.zip
