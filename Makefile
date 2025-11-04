clear:
	clear

.PHONY: build
build: clear
	go build -o at-plan main.go && \
	sudo mv at-plan /usr/local/bin/

amd-runner-build:
	GOOS=linux GOARCH=amd64 go build -o at-plan main.go

arm-runner-build:
	GOOS=linux GOARCH=arm64 go build -o at-plan main.go

run: build
	at-plan plan \
	-g github-token \
	-f master \
	-u https://atlantis.dev.leedonggyu.com \
	-t leedonggyu-1234 \
	-r zkfmapf123/atlantis-fargate \
	-c atlantis.yaml \ 
	-s slackt-webhook-url \
	-l channel 

test:
	go test -v ./...