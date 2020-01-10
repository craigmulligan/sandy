.PHONY: install 
install: go.sum
	go get -v -t


.PHONY: build 
build: install
	go build -o sandy 

.PHONY: test
test: build
	go test	

.PHONY: dist
dist:
	env GOOS=linux GOARCH=amd64 go build -o ./dist/sandy_linux_amd64
