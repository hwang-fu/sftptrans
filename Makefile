.PHONY: all build frontend backend clean run

all: build

frontend:
	cd web/sftptrans && npm ci && npm run build -- --configuration=production --output-path=../../internal/server/static

backend: frontend
	go build -o sftptrans ./cmd/sftptrans

build: backend

clean:
	rm -rf sftptrans internal/server/static

# ./sftptrans -h <host> -u <username> -pass <password>
