all: linux windows

dev:
	wails dev

linux:
	GOOS=linux GOARCH=amd64 wails build

windows:
	GOOS=windows GOARCH=amd64 wails build --clean 
