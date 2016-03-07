test:
	GO15VENDOREXPERIMENT=1 go test . -cover

build:
	GO15VENDOREXPERIMENT=1 go build .

release:
	@env GO15VENDOREXPERIMENT=1 GOOS=linux GOARCH=amd64 go build -ldflags="-X main.appVersion=${DROPLAN_VERSION}" -o droplan
	@zip droplan_${DROPLAN_VERSION}_linux_amd64.zip droplan
	@rm droplan

	@env GO15VENDOREXPERIMENT=1 GOOS=linux GOARCH=386 go build -ldflags="-X main.appVersion=${DROPLAN_VERSION}" -o droplan
	@zip droplan_${DROPLAN_VERSION}_linux_386.zip droplan
	@rm droplan

clean:
	@rm -f droplan
	@rm -rf droplan_*.zip
