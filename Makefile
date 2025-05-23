.PHONY: all test start-dynamodb check-java clean

JAVA_VERSION := 21
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
    echo $$! > dynamodb_local.pid

# Run tests
test: start-dynamodb
	@echo "Running tests..."
	@go test ./... -v

# Clean up
clean:
	@echo "Stopping DynamoDB Local..."
	@if [ -f dynamodb_local.pid ]; then \
		PID=$$(cat dynamodb_local.pid); \
		if ps -p $$PID > /dev/null 2>&1; then \
			kill $$PID; \
			echo "DynamoDB Local stopped."; \
		else \
			echo "No process found with PID $$PID."; \
		fi; \
		rm -f dynamodb_local.pid; \
	else \
		echo "No PID file found. DynamoDB Local may not be running."; \
	fi