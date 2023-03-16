all: linux windows android

linux:
	GOOS=linux GOARCH=amd64 go build -o myprogram-linux main.go

windows:
	CGO_ENABLED=1 CC=x86_64-w64-mingw32-cc GOOS=windows GOARCH=amd64 go build -o myprogram-windows.exe main.go

android:
	GOOS=android GOARCH=arm go build -o myprogram-android main.go

clean:
	rm -f myprogram-linux myprogram-windows.exe myprogram-android
