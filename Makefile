clear:
	clear

.PHONY: build
build: clear
	go build -o at-plan main.go
	sudo mv at-plan /usr/local/bin/

run: build
	at-plan plan \
	-g github-token \
	-f master \
	-u https://atlantis.dev.leedonggyu.com \
	-t leedonggyu-1234 \
	-r zkfmapf123/atlantis-fargate \
	-c atlantis.yaml

test:
	go test -v ./...