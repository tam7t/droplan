test:
	GO15VENDOREXPERIMENT=1 go test . -cover

build:
	mkdir -p build
	GO15VENDOREXPERIMENT=1 GOOS=linux GOARCH=amd64 go build -o build/dolan_linux_amd64 .
	GO15VENDOREXPERIMENT=1 GOOS=linux GOARCH=386 go build -o build/dolan_linux_386 .
	GO15VENDOREXPERIMENT=1 GOOS=freebsd GOARCH=amd64 go build -o build/dolan_freebsd_amd64 .
	GO15VENDOREXPERIMENT=1 GOOS=freebsd GOARCH=386 go build -o build/dolan_freebsd_386 .

build_local:
	GO15VENDOREXPERIMENT=1 go build .

clean:
	rm -rf build
