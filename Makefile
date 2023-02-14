dev-gen:
	go build -o xtool *.go

dev-clean:
	rm xtool

gen:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o xtool *.go