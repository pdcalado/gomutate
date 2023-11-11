test:
	go run main.go -type Acme ./testdata/acme.go > testdata/mutations.go
	go run testdata/*.go | diff - testdata/expected.txt