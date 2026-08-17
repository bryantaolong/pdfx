.PHONY: all linux win clean

VERSION := 0.2.1
GO := go
BIN_DIR := bin

all: linux win

linux:
	@mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build -trimpath -ldflags "-s -w" -o $(BIN_DIR)/pdfx-$(VERSION)-linux-amd64 .

win:
	@mkdir -p $(BIN_DIR)
	GOOS=windows GOARCH=amd64 $(GO) build -trimpath -ldflags "-s -w" -o $(BIN_DIR)/pdfx-$(VERSION)-windows-amd64.exe .

clean:
	rm -rf $(BIN_DIR)/*
