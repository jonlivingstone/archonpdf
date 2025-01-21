VERSION = snapshot

build:	
	go build -v -buildvcs=false -ldflags="-s -w" -o mergepdf ./cmd/

dockerbuild:
	docker run --rm \
		-v "$$PWD":/usr/src/myapp \
		-w /usr/src/myapp \
		golang:latest \
		go build -v -buildvcs=false -ldflags="-s -w" -o mergepdf ./cmd/

clean:
	rm mergepdf

install:
	docker build -t mergepdf:$(VERSION) .

export: install
	docker save mergepdf -o mergepdf-$(VERSION).tar
