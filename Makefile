dev-gen:
	go build -o xtool-dev *.go

clean:
	rm xtool*

released:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o xtool *.go
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o xtool-arm *.go