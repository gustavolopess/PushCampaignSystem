BASE_PATH       	:= $(shell pwd | sed 's/ /\\ /g')
BASE_PACKAGE_SRC 	:= $(BASE_PATH)/src/github.com/gustavolopess/PushCampaignSystem

PUBLISHER_MAIN 		:= $(BASE_PACKAGE_SRC)/cmd/pub/main.go
SUBSCRIBER_MAIN		:= $(BASE_PACKAGE_SRC)/cmd/sub/main.go

PUBLISHER_BIN		:= publisher
SUBSCRIBER_BIN		:= subscriber

all: build run-publisher run-subscriber

build:
	@mkdir -p bin
	@cd $(BASE_PACKAGE_SRC) && go build -o $(PUBLISHER_BIN) $(PUBLISHER_MAIN)
	@mv $(BASE_PACKAGE_SRC)/$(PUBLISHER_BIN) bin/$(PUBLISHER_BIN)
	@cd $(BASE_PACKAGE_SRC) && go build -o $(SUBSCRIBER_BIN) $(SUBSCRIBER_MAIN)
	@mv $(BASE_PACKAGE_SRC)/$(SUBSCRIBER_BIN) bin/$(SUBSCRIBER_BIN)