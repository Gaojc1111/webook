.PHONY: docker
docker:
	@rm webook || true
	@GOOS=linux GOARCH=arm go build -o webook .
	@docker rmi -f hbzhtd/webook
	@docker build -t hbzhtd/webook:v0.0.1 .

.PHONY: mock
mock:
	@mockgen -source=./internal/service/user.go -package=mocksvc -destination=./internal/service/mock/user.mock.go
	@mockgen -source=./internal/service/code.go -package=mocksvc -destination=./internal/service/mock/code.mock.go

	@mockgen -source=./internal/repository/user.go -package=mocksvc -destination=./internal/repository/mock/user.mock.go
	@mockgen -source=./internal/repository/code.go -package=mocksvc -destination=./internal/repository/mock/code.mock.go

	@mockgen -source=./internal/repository/dao/user.go -package=dao_mocksvc -destination=./internal/repository/dao/mock/user.mock.go
	@mockgen -source=./internal/repository/cache/user.go -package=cache_mocksvc -destination=./internal/repository/cache/mock/user.mock.go
	@mockgen -source=./internal/repository/cache/code.go -package=cache_mocksvc -destination=./internal/repository/cache/mock/code.mock.go

	@mockgen -package=redis_mock -destination=./internal/repository/cache/redis_mock/cmd.mock.go github.com/redis/go-redis/v9 Cmdable

	@go mod tidy
