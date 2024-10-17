.PHONY: test cover build

# Переменные
BUILD_DIR := build

# Юнит тесты и покрытие кода
test:
	go test -race -count 1 ./...

cover:
	go test -short -race -count 1 -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

# Сборка
build:
	mkdir -p $(BUILD_DIR)
	rm -rf $(BUILD_DIR)/*
	#GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o $(BUILD_DIR)/sing-box-linux-amd64 main.go
	#GOOS=linux GOARCH=386 CGO_ENABLED=1 go build -o $(BUILD_DIR)/sing-box-linux-i386 main.go
	#GOOS=linux GOARCH=arm64 CGO_ENABLED=1 go build -o $(BUILD_DIR)/sing-box-linux-arm64 /main.go
	#GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -o $(BUILD_DIR)/sing-box-windows-amd64.exe /main.go
	#GOOS=windows GOARCH=386 CGO_ENABLED=1 go build -o $(BUILD_DIR)/sing-box-windows-i386.exe main.go
	#GOOS=windows GOARCH=arm64 CGO_ENABLED=1 go build -o $(BUILD_DIR)/sing-box-windows-arm64.exe main.go
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -o $(BUILD_DIR)/sing-box-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -o $(BUILD_DIR)/sing-box-darwin-arm64 main.go