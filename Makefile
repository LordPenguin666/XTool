dev-gen:
	go build -o xtool-dev *.go

dev-gen-arm:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o xtool-dev-arm *.go

dev-clean:
	rm xtool

gen:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o xtool *.go