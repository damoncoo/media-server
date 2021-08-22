
build:
	@go build -v -a -ldflags '-s -w' \
	-trimpath \
	-o media-server server/*.go

subtitle:
  @subify dl