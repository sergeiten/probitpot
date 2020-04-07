build:
	env GOOS=windows GOARCH=386 go build -o bin/win32/probitpot.exe cmd/probitpot/main.go
	env GOOS=windows GOARCH=amd64 go build -o bin/win64/probitpot.exe cmd/probitpot/main.go
	env GOOS=darwin GOARCH=amd64 go build -o bin/darwin/probitpot cmd/probitpot/main.go
	env GOOS=linux GOARCH=386 go build -o bin/linux/probitpot cmd/probitpot/main.go
	env GOOS=linux GOARCH=amd64 go build -o bin/linux/probitpot cmd/probitpot/main.go
