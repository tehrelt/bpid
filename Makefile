build:
	make caesar
	make feistel

caesar:
	go build -o bin/caesar.exe -v cmd/caesar/main.go

feistel:
	go build -o bin/feistel.exe -v cmd/feistel/main.go

coverage:
	go test -v -coverprofile cover.out ./...
	go tool cover -html cover.out -o cover.html
	rm cover.out