BIN_DIR=bin
CMD_DIR=cmd

build:
	mkdir -p $(BIN_DIR) #
	go build -o $(BIN_DIR)/stego.exe -v $(CMD_DIR)/stego/main.go

coverage:
	go test -v -coverprofile cover.out ./...
	go tool cover -html cover.out -o cover.html
	rm cover.out

BIN=./bin/stego.exe
KEY=reltreltreltrelt
SRC_IMG=./test/src.bmp
TARGET=./test/main
EMBED=./test/embedded.bmp
EXTRACTED=./test/extracted

encode:
	$(BIN) \
		-mode e \
		-key $(KEY) \
		-in $(TARGET) \
		-img $(SRC_IMG) \
		-out $(EMBED)

decode:
	$(BIN) \
		-mode d \
		-key $(KEY) \
		-img $(EMBED) \
		-out $(EXTRACTED)

testing: encode decode

clean:
	rm $(EMBED) $(EXTRACTED)