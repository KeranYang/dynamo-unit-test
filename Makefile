.PHONY: all test start-dynamodb check-java clean

JAVA_VERSION := 11
JAVA_PACKAGE := openjdk-$(JAVA_VERSION)-jdk

# Check for Java installation
check-java:
	@echo "Checking for Java installation..."
	@if ! command -v java &> /dev/null; then \
		echo "Java not found. Installing..."; \
		sudo apt update && sudo apt install -y $(JAVA_PACKAGE); \
	else \
		echo "Java is already installed."; \
	fi

# Run the local DynamoDB, redirecting output to /dev/null and using a subshell
start-dynamodb: check-java
	@echo "Starting DynamoDB Local..."
	@java -Djava.library.path=./DynamoDBLocal_lib -jar DynamoDBLocal.jar -sharedDb > /dev/null 2>&1 & \

# Run tests
test: start-dynamodb
	@echo "Running tests..."
	@go test ./... -v

# Clean up
clean:
	@echo "Stopping DynamoDB Local..."
	@pkill -f DynamoDBLocal.jar || true