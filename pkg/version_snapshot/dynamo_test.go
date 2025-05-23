package version_snapshot

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func Test(t *testing.T) {
	// Setup local DynamoDB
	dynamoManager, err := prepareTestEnv()
	if err != nil {
		t.Fatalf("Failed to setup local DynamoDB: %v", err)
	}

	// Test WriteVersionSnapshot
	snapshot := &VersionSnapshot{
		ID: ID{
			AssetID:      "test-asset-id",
			ResourceName: "test-resource-name",
			ResourceType: "test-resource-type",
			Location:     "test-location",
		},
		SnapshotSpec: []byte("test"),
	}
	err = dynamoManager.WriteVersionSnapshot(snapshot)
	if err != nil {
		t.Fatalf("Failed to write version snapshot: %v", err)
	}

	// Test GetVersionSnapshot
	id := ID{
		AssetID:      "test-asset-id",
		ResourceName: "test-resource-name",
		ResourceType: "test-resource-type",
		Location:     "test-location",
	}
	snapshot, err = dynamoManager.GetVersionSnapshot(id)
	if err != nil {
		t.Fatalf("Failed to get version snapshot: %v", err)
	}
	if snapshot == nil {
		t.Fatalf("Expected snapshot to be not nil")
	}
	if string(snapshot.SnapshotSpec) != "test" {
		t.Fatalf("Expected snapshot spec to be 'test', got '%s'", string(snapshot.SnapshotSpec))
	}

	// Delete the test table
	if err := cleanupTestEnv(dynamoManager); err != nil {
		t.Fatalf("Failed to cleanup local DynamoDB: %v", err)
	}
}

func prepareTestEnv() (*VersionSnapshotDynamoManager, error) {
	dynamoManager := NewVersionSnapshotDynamoManager("us-west-2", "http://localhost:8000")
	// Create a local DynamoDB table for testing
	_, err := dynamoManager.dynamoDbClient.CreateTable(context.TODO(),
		&dynamodb.CreateTableInput{
			TableName: aws.String(TableName),
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("partitionKey"),
					KeyType:       types.KeyTypeHash, // Partition key
				},
				{
					AttributeName: aws.String("sortKey"),
					KeyType:       types.KeyTypeRange, // Sort key
				},
			},
			AttributeDefinitions: []types.AttributeDefinition{
				{
					AttributeName: aws.String("partitionKey"),
					AttributeType: types.ScalarAttributeTypeS,
				},
				{
					AttributeName: aws.String("sortKey"),
					AttributeType: types.ScalarAttributeTypeS,
				},
			},
			ProvisionedThroughput: &types.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(5),
				WriteCapacityUnits: aws.Int64(5),
			},
		})
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}
	// Verify that the table was created successfully
	describeOutput, err := dynamoManager.dynamoDbClient.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(TableName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to describe table: %v", err)
	}
	if describeOutput.Table.TableStatus != types.TableStatusActive {
		return nil, fmt.Errorf("expected table status to be ACTIVE, got %s", describeOutput.Table.TableStatus)
	}
	return dynamoManager, nil
}

func cleanupTestEnv(dynamoManager *VersionSnapshotDynamoManager) error {
	_, err := dynamoManager.dynamoDbClient.DeleteTable(context.TODO(), &dynamodb.DeleteTableInput{
		TableName: aws.String(TableName),
	})
	if err != nil {
		return fmt.Errorf("failed to delete table: %v", err)
	}
	return nil
}
