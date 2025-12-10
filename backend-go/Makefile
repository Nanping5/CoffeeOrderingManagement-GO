.PHONY: run build clean test docker-build docker-run

# 运行开发服务器
run:
	go run main.go

# 编译
build:
	go build -o coffee-ordering-backend main.go

# 清理
clean:
	rm -f coffee-ordering-backend
	go clean

# 测试
test:
	go test -v ./...

# 安装依赖
deps:
	go mod download
	go mod tidy

# Docker 构建
docker-build:
	docker build -t coffee-ordering-backend:latest .

# Docker 运行
docker-run:
	docker run -p 8080:8080 --env-file .env coffee-ordering-backend:latest

# 格式化代码
fmt:
	go fmt ./...

# 代码检查
lint:
	golangci-lint run

# 热重载（需要安装 air）
dev:
	air
