include .env
export $(shell sed 's/=.*//' .env)

PROJECT_NAME=ntp
BINARY_NAME=${PROJECT_NAME}

# Run

parse:
	cd cmd/cli && go run .

serve-teleblog:
	npx tailwindcss build -i tailwind.css -o cmd/teleblog/httpapi/public/style.css
	cd cmd/teleblog \
	&& go generate ./... \
	&& go run . serve

live-teleblog:
	templ generate --watch --proxy="http://localhost:8090" --cmd="make serve-teleblog"

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
	GOARCH=amd64 GOOS=darwin go build -o ./cmd/teleblog/${BINARY_NAME}-teleblog-darwin ./cmd/teleblog

clean-mac:
	go clean
	rm ${BINARY_NAME}-darwin

build-cli-linux:
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-cli-linux ./cmd/cli

build-teleblog-linux:
	npx tailwindcss build -i tailwind.css -o cmd/teleblog/httpapi/public/style.css
	make templ
	GOARCH=amd64 GOOS=linux go build -o ./cmd/teleblog/${BINARY_NAME}-teleblog-linux ./cmd/teleblog

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
	scp ./infra/pocketbase.service root@${SERVER_IP}:/lib/systemd/system/pocketbase.service
	ssh root@${SERVER_IP} "apt update \
	&& mkdir -p /root/ntp \
	&& systemctl enable pocketbase \
	&& systemctl daemon-reload"
	scp ./cmd/teleblog/app.env.example root@${SERVER_IP}:/root/ntp/app.env

# Deploy

deploy:
	make build-teleblog-linux
	ssh root@${SERVER_IP} "systemctl stop pocketbase"
	scp ./cmd/teleblog/${BINARY_NAME}-teleblog-linux root@${SERVER_IP}:/root/ntp/${BINARY_NAME}-teleblog-linux
	ssh root@${SERVER_IP} "systemctl start pocketbase"