all: cfg genan

.PHONY: cfg genan
cfg:
	go build ./cmd/cfg

genan:
	go build ./cmd/genan
