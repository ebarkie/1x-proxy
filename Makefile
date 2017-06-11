NAME=1x-proxy

GO=$(shell which go)

.PHONY: default build clean mipsle mips64

default: clean build
	@echo "Native build complete"
erl erp: clean mips64
	@echo "EdgeRouter Lite/Pro build complete"
erx erxsfp: clean mipsle
	@echo "EdgeRouter X [SFP] build complete"

mips64: export GOARCH=mips64
mips64: export GOOS=linux

mipsle: export GOARCH=mipsle
mipsle: export GOOS=linux

build mipsle mips64:
	@echo "Building 802.1x proxy.."
	@$(GO) generate
	@$(GO) build -o $(NAME)

clean:
	@$(GO) clean
	@rm -f $(NAME)
