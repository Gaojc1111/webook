.PHONY: docker
docker:
	@rm webook || true
	@GOOS=linux GOARCH=arm go build -o webook .
	@docker rmi -f hbzhtd/webook
	@docker build -t hbzhtd/webook:v0.0.1 .
