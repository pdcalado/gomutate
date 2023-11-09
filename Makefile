test:
	go run main.go -type Acme testdata > testdata/mutations.go
	go run testdata/*