APP_NAME=go-service
BUILD_DIR=.

clean:
	@echo "Cleaning..."
	-rm rf $(BUILD_DIR)

swagger:
    swag init -g main.go -o docs

run: build
	@echo "Running $(APP_NAME)..."
	./$(BUILD_DIR)/$(APP_NAME)

test:
	@echo "Running tests..."
	go test ./... -v -cover
build:
	@echo "Building $(APP_NAME)"
	go build main.go -o $(BUILD_DIR)/$(APP_NAME) 
all: clean build

.PHONY all clean build