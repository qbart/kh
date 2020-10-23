build:
	mkdir -p bin/
	go build -o bin/kh main.go

install: build
	cp bin/kh ~/bin

