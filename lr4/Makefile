BIN_DIR=bin
CMD_DIR=cmd
KEYS_DIR=keys

build:
	mkdir -p $(BIN_DIR) #
	make aes rsa dsa

keys-dir:
	mkdir -p $(KEYS_DIR)

RSA_PRIVATE=$(KEYS_DIR)/rsa
RSA_PUBLIC=$(KEYS_DIR)/rsa.pub
rsa-keygen: keys-dir
	openssl genpkey -algorithm RSA -out $(RSA_PRIVATE) -pkeyopt rsa_keygen_bits:2048
	openssl rsa -pubout -in $(RSA_PRIVATE) -out $(RSA_PUBLIC)

DSA_PRIVATE = $(KEYS_DIR)/dsa
DSA_PUBLIC = $(KEYS_DIR)/dsa.pub
dsa-keygen: keys-dir 
	openssl genpkey -algorithm Ed25519 -out $(DSA_PRIVATE)
	openssl pkey -in $(DSA_PRIVATE) -pubout -out $(DSA_PUBLIC)

aes:
	go build -o $(BIN_DIR)/aes.exe -v $(CMD_DIR)/aes/main.go

rsa:
	go build -o $(BIN_DIR)/rsa.exe -v $(CMD_DIR)/rsa/main.go

dsa:
	go build -o $(BIN_DIR)/dsa.exe -v $(CMD_DIR)/dsa/main.go

all: build dsa-keygen rsa-keygen

clean:
	rm -rf $(BIN_DIR)
	rm -rf $(KEYS_DIR)


coverage:
	go test -v -coverprofile cover.out ./...
	go tool cover -html cover.out -o cover.html
	rm cover.out