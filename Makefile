-include .env
export $(shell sed 's/=.*//' .env)

PROJECT_NAME=teleblog
BINARY_NAME=${PROJECT_NAME}

# Run

parse:
	cd cmd/cli && go run .

serve-teleblog:
	npx tailwindcss build -i tailwind.css -o cmd/teleblog/httpapi/public/style.css
	cd cmd/teleblog \
	&& go generate ./... \
	&& go run . serve

upload-history:
	cd cmd/teleblog \
	&& go generate ./... \
	&& go run . upload-history

# Cmds

extract-tags:
	cd cmd/teleblog \
	&& go generate ./... \
	&& go run . extract-tags

# Generate

templ:
	cd cmd/teleblog \
	&& go generate ./...

# Build

build-cli:
	make build-cli-mac && make build-cli-linux

build-cli-mac:
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-cli-darwin ./cmd/cli

build-teleblog-mac:
	npx tailwindcss build -i tailwind.css -o cmd/teleblog/httpapi/public/style.css
	make templ
	GOARCH=amd64 GOOS=darwin go build -o ./cmd/teleblog/${BINARY_NAME}-darwin ./cmd/teleblog

clean-mac:
	go clean
	rm ${BINARY_NAME}-darwin

build-cli-linux:
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-cli-linux ./cmd/cli

build-teleblog-linux:
	npx tailwindcss build -i tailwind.css -o cmd/teleblog/httpapi/public/style.css
	make templ
	GOARCH=amd64 GOOS=linux go build -o ./cmd/teleblog/${BINARY_NAME}-linux ./cmd/teleblog

clean:
	go clean
	rm ${BINARY_NAME}-cli-darwin
	rm ${BINARY_NAME}-cli-linux
	rm ${BINARY_NAME}-teleblog-darwin
	rm ${BINARY_NAME}-teleblog-linux

# Setup

setup:
	npm i
	go install github.com/a-h/templ/cmd/templ@latest
	go mod tidy

setup-droplet:
	scp ./infra/teleblog.service root@${SERVER_IP}:/lib/systemd/system/teleblog.service
	ssh root@${SERVER_IP} "apt update \
	&& mkdir -p /root/teleblog \
	&& systemctl enable teleblog \
	&& systemctl daemon-reload"
	scp ./cmd/teleblog/app.env.example root@${SERVER_IP}:/root/teleblog/app.env
	scp ./infra/davidshekunts.ru root@${SERVER_IP}:/etc/nginx/sites-available/davidshekunts.ru
	scp ./infra/davidshekunts.ru root@${SERVER_IP}:/etc/nginx/sites-enabled/davidshekunts.ru

# Deploy

deploy:
	make build-teleblog-linux
	ssh root@${SERVER_IP} "systemctl stop teleblog"
	scp ./cmd/teleblog/${BINARY_NAME}-linux root@${SERVER_IP}:/root/teleblog/${BINARY_NAME}-linux
	ssh root@${SERVER_IP} "systemctl start teleblog"