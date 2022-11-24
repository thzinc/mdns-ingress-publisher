.PHONY: build

build:
	GOOS=darwin GOARCH=amd64 go build -o artifacts/darwin-amd64/mdns-ingress-publisher cmd/mdns-ingress-publisher/main.go
	GOOS=darwin GOARCH=arm64 go build -o artifacts/darwin-arm64/mdns-ingress-publisher cmd/mdns-ingress-publisher/main.go
	GOOS=linux GOARCH=amd64 go build -o artifacts/linux-amd64/mdns-ingress-publisher cmd/mdns-ingress-publisher/main.go
	GOOS=linux GOARCH=arm go build -o artifacts/linux-arm/mdns-ingress-publisher cmd/mdns-ingress-publisher/main.go
	GOOS=windows GOARCH=amd64 go build -o artifacts/windows-amd64/mdns-ingress-publisher.exe cmd/mdns-ingress-publisher/main.go
	GOOS=windows GOARCH=arm go build -o artifacts/windows-arm/mdns-ingress-publisher.exe cmd/mdns-ingress-publisher/main.go

microk8s-systemd: microk8s-systemd-clean
	mkdir -p artifacts/microk8s-systemd
	cp -r build/microk8s-systemd/* artifacts/microk8s-systemd/

microk8s-systemd-clean:
	rm -rf artifacts/microk8s-systemd