# sshi

Simple interactive SSH manager that lets you select and connect to hosts from your ~/.ssh/config or similar files.

## Usage

```bash
# build app
GOOS=windows GOARCH=amd64 go build -o sshi_windows_amd64.exe main.go
GOOS=darwin  GOARCH=amd64 go build -o sshi_macos_amd64 main.go
GOOS=darwin  GOARCH=arm64 go build -o sshi_macos_arm64 main.go
GOOS=linux   GOARCH=amd64 go build -o sshi_linux_amd64 main.go

# check example
./sshi_linux_amd64 config.example

# for use
sudo mv ./sshi_linux_amd64 /usr/local/bin
```
