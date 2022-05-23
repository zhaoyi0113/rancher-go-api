build:
	go build -o dist/api

buildimage:
	docker build --platform=linux/amd64 -t zhaoyi0113/rancher-go-api .

publishimage:
	docker push zhaoyi0113/rancher-go-api

