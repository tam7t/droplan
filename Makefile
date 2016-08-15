test:
	GO15VENDOREXPERIMENT=1 go test . -cover

build:
	GO15VENDOREXPERIMENT=1 go build .

build-amd64:
	@env GO15VENDOREXPERIMENT=1 GOOS=linux GOARCH=amd64 go build -ldflags="-X main.appVersion=${DROPLAN_VERSION}" -o droplan

build-i386:
	@env GO15VENDOREXPERIMENT=1 GOOS=linux GOARCH=386 go build -ldflags="-X main.appVersion=${DROPLAN_VERSION}" -o droplan_i386

release: build-amd64 build-i386
	@zip droplan_${DROPLAN_VERSION}_linux_amd64.zip droplan
	@tar -cvzf droplan_${DROPLAN_VERSION}_linux_amd64.tar.gz droplan
	@rm droplan

	@mv droplan_i386 droplan
	@zip droplan_${DROPLAN_VERSION}_linux_386.zip droplan
	@tar -cvzf droplan_${DROPLAN_VERSION}_linux_386.tar.gz droplan
	@rm droplan

docker: build-amd64
	@docker build -t tam7t/droplan .


clean:
	@rm -f droplan
	@rm -rf droplan_*.zip
	@rm -rf droplan_*.tar.gz
