all: windows osx64 linux

windows:
	GOOS=windows GOARCH=386 go build -o bin/arachne-386.exe main.go server.go graph.go subject.go

osx64:
	GOOS=darwin GOARCH=amd64 go build -o bin/arachne-amd64-darwin main.go server.go graph.go subject.go


linux: 
	GOOS=linux GOARCH=amd64 go build -o bin/arachne-amd64-linux main.go server.go graph.go subject.go

