Windows:
	GOOS=windows	GOARCH=amd64	go build -o target/windows-amd64-neonx.exe -ldflags "-X main.tag=$(version)" main.go

MacOS:
	GOOS=darwin	GOARCH=amd64	go build -o target/mac-amd64-neonx -ldflags "-X main.tag=$(version)" main.go

Linux:
	GOOS=linux	GOARCH=amd64	go build -o target/linux-amd64-neonx -ldflags "-X main.tag=$(version)" main.go

all: clean MacOS Linux Windows

clean:
	rm -fr target

build: all
	tar cvf neonx-cli.tar target

install:
	go build -o /usr/local/bin/neonx -ldflags "-X main.tag=$(version)" main.go
