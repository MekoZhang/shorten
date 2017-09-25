version = 1.2.0

dep:
	go get -d ./...

test: build
	go test -v ./...

vet:
	go list ./... | grep -v "./vendor*" | xargs go vet

fmt:
	find . -type f -name "*.go" | grep -v "./vendor*" | xargs gofmt -s -w

build: dep vet fmt
	go build -ldflags="-X github.com/zhangxd1989/shorten/conf.Version=$(version)" -o shorten main.go

clean:
	rm -f shorten

start:
	nohup ./shorten -c config.conf 1> shorten.out 2> shorten.err &