BINARY_NAME = algorithms
TARGET_DIR =  bin # имя папки для бинарников
SRC_MAIN = main.go # имя файла с главной функцией
.DEFAULT_GOAL := run

.PHONY: build build-windows build-darwin build-linux run-darwin run-linux run-windows run clean test info

build:
	go build -o $(TARGET_DIR)/$(BINARY_NAME) $(SRC_MAIN)

build-windows:
	GOARCH=amd64 GOOS=windows go build -o $(TARGET_DIR)/$(BINARY_NAME)-windows.exe $(SRC_MAIN)

build-darwin:
	GOARCH=amd64 GOOS=darwin go build -o $(TARGET_DIR)/$(BINARY_NAME)-darwin $(SRC_MAIN)

build-linux:
	GOARCH=amd64 GOOS=linux go build -o $(TARGET_DIR)/$(BINARY_NAME)-linux $(SRC_MAIN)

run: build
	./$(TARGET_DIR)/$(BINARY_NAME)

run-windows: build-windows
	./$(TARGET_DIR)/$(BINARY_NAME)-windows.exe

run-darwin: build-darwin
	./$(TARGET_DIR)/$(BINARY_NAME)-darwin

run-linux: build-linux
	./$(TARGET_DIR)/$(BINARY_NAME)-linux

clean:
	rm -rf $(TARGET_DIR)

test:
	go test -race ./... -v

info:
	@echo "== Photographer's  optimization app =="
	@echo ""
	@echo "Go Version:"
	@go version
	@echo ""
	@echo "Module Info:"
	@echo " Name: $(shell go list -m)"
	@echo ""
	@echo "Repository status"
	@echo "Branch: $(shell git branch --show-current 2>nul)"
	@echo "Changes: "
	@if "$(shell git status --short 2>nul)" == "" (\
    		echo "no changes"\
    	) else (\
      		git status --short 2>nul\
    	)
	@echo "Legend:"
	@powershell -Command "Write-Host 'R' -ForegroundColor Green -NoNewline; Write-Host 'M' -ForegroundColor Red -NoNewline; Write-Host ' - renamed'"
	@powershell -Command "Write-Host 'M' -ForegroundColor Red -NoNewline; Write-Host ' - modified'"
	@powershell -Command "Write-Host '??' -ForegroundColor Red -NoNewline; Write-Host ' - untracked'"
	@powershell -Command "Write-Host 'D' -ForegroundColor Red -NoNewline; Write-Host ' - deleted'"
	@echo ""
	@echo " Commands:"
	@echo " make build - build project for current platform"
	@echo " make build-windows - build project for windows"
	@echo " make build-darwin - build project for darwin"
	@echo " make build-linux - build project for linux"
	@echo " make run - build for current platform and run"
	@echo " make run-darwin - build project for darwin and run"
	@echo " make run-linux - build project for linux and run"
	@echo " make run-windows - build project for windows and run"
	@echo " make clean - clean build"
	@echo " make info - this info"
	@echo " make - default command: run"
	@echo ""
