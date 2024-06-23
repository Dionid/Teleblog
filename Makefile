include .env
export $(shell sed 's/=.*//' .env)

PROJECT_NAME=ntp
BINARY_NAME=${PROJECT_NAME}

# Run

parse:
	cd cmd/cli && go run .

serve-saas:
	npx tailwindcss build -i tailwind.css -o cmd/saas/httph/public/style.css
	cd cmd/saas \
	&& go generate ./... \
	&& go run . serve

live-saas:
	templ generate --watch --proxy="http://localhost:8090" --cmd="make serve-saas"

# Generate

templ:
	cd cmd/saas \
	&& go generate ./...

# Build

build-cli:
	make build-cli-mac && make build-cli-linux

build-cli-mac:
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-cli-darwin ./cmd/cli

build-saas-mac:
	npx tailwindcss build -i tailwind.css -o cmd/saas/httph/public/style.css
	make templ
	GOARCH=amd64 GOOS=darwin go build -o ./cmd/saas/${BINARY_NAME}-saas-darwin ./cmd/saas

clean-mac:
	go clean
	rm ${BINARY_NAME}-darwin

build-cli-linux:
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-cli-linux ./cmd/cli

build-saas-linux:
	npx tailwindcss build -i tailwind.css -o cmd/saas/httph/public/style.css
	make templ
	GOARCH=amd64 GOOS=linux go build -o ./cmd/saas/${BINARY_NAME}-saas-linux ./cmd/saas

clean:
	go clean
	rm ${BINARY_NAME}-cli-darwin
	rm ${BINARY_NAME}-cli-linux
	rm ${BINARY_NAME}-saas-darwin
	rm ${BINARY_NAME}-saas-linux

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
	scp ./cmd/saas/app.env.example root@${SERVER_IP}:/root/ntp/app.env

# Deploy

deploy:
	make build-saas-linux
	ssh root@${SERVER_IP} "systemctl stop pocketbase"
	scp ./cmd/saas/${BINARY_NAME}-saas-linux root@${SERVER_IP}:/root/ntp/${BINARY_NAME}-saas-linux
	ssh root@${SERVER_IP} "systemctl start pocketbase"