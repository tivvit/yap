test: $(wildcard **/*.go)
	go test -v -count=1 ./...
yap: $(wildcard **/*.go)
	go build .
install: yap
	cp yap /usr/local/bin/