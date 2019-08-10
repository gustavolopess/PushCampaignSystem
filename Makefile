BASE_PACKAGE		:= github.com/gustavolopess/PushCampaignSystem
BASE_PACKAGE_SRC 	:= src/$(BASE_PACKAGE)

PUBLISHER_MAIN 		:= ./cmd/pub/
SUBSCRIBER_MAIN		:= ./cmd/sub/
READER_MAIN			:= ./cmd/reader/

PUBLISHER_BIN		:= publisher
SUBSCRIBER_BIN		:= subscriber
READER_BIN			:= reader

all: mod build run-publisher run-subscriber

mod:
	@cd $(BASE_PACKAGE_SRC) && go mod init $(BASE_PACKAGE)

deps:
	@cd $(BASE_PACKAGE_SRC) && go build ./...

build:
	@mkdir -p bin
	@cd $(BASE_PACKAGE_SRC) && go build -o $(PUBLISHER_BIN) $(PUBLISHER_MAIN)
	@mv $(BASE_PACKAGE_SRC)/$(PUBLISHER_BIN) bin/$(PUBLISHER_BIN)
	@cd $(BASE_PACKAGE_SRC) && go build -o $(SUBSCRIBER_BIN) $(SUBSCRIBER_MAIN)
	@mv $(BASE_PACKAGE_SRC)/$(SUBSCRIBER_BIN) bin/$(SUBSCRIBER_BIN)
	@cd $(BASE_PACKAGE_SRC) && go build -o $(READER_BIN) $(READER_MAIN)
	@mv $(BASE_PACKAGE_SRC)/$(READER_BIN) bin/$(READER_BIN)