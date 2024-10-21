lotor: schema.go *.go
	go build -o lotor

schema.go: schema.bare
	go run lotor/bareish/baregen schema.bare schema.go
