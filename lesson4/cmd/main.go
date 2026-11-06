package main

import (
	"os"

	"gostudy/lesson4/internal/storage"
)

func main() {
	runner := commandRunner{store: storage.FileStore{}}
	os.Exit(runner.run(os.Args[1:], os.Stdout, os.Stderr))
}
