VERSION = snapshot

build:
	go build -v -buildvcs=false -ldflags="-s -w" -o build/archonpdf ./cmd/

dockerbuild:
	docker run --rm \
		-v "$$PWD":/usr/src/myapp \
		-w /usr/src/myapp \
		golang:latest \
		go build -v -buildvcs=false -ldflags="-s -w" -o build/archonpdf ./cmd/

clean:
	rm build/*

install:
	docker build -t archonpdf:$(VERSION) .

export: install
	docker save archonpdf -o archonpdf-$(VERSION).tar
