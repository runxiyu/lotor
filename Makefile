lotor: schema.go *.go
	go build -o lotor

schema.go: schema.bare
	go run git.sr.ht/~runxiyu/lotor/bareish/baregen -p main schema.bare schema.go
