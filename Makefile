lotor: schema.go *.go
	go build -o lotor

schema.go: schema.bare
	go-bare-gen -p main schema.bare schema.go
