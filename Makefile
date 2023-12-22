.PHONY: docker
docker:
	@rm gint || true
	@GOOS=linux GOARCH=arm go build -tags=k8s -o ./build-bin/gint .
	@docker rmi -f bswaterb/gint:0.0.1
	@docker build -t bswaterb/gint:0.0.1 .