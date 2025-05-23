package version_snapshot

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const TableName = "version_snapshot"

type dynamoVersionSnapshotItem struct {
	VersionSnapshot `json:"versionSnapshot"`
	PartitionKey    string `json:"partitionKey"`
	SortKey         string `json:"sortKey"`
}

var _ VersionSnapshotManager = &VersionSnapshotDynamoManager{}

type VersionSnapshotDynamoManager struct {
	dynamoDbClient *dynamodb.Client
}

func NewVersionSnapshotDynamoManager(region, endpoint string) *VersionSnapshotDynamoManager {
	credentials := aws.NewCredentialsCache(
		// TODO - conditionally use this for local testing
		credentials.NewStaticCredentialsProvider("dummy", "dummy", ""),
	)
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           endpoint,
					SigningRegion: region,
				}, nil
			})),
		config.WithCredentialsProvider(credentials),
	)
	if err != nil {
		return nil
	}
	client := dynamodb.NewFromConfig(cfg)
	return &VersionSnapshotDynamoManager{
		dynamoDbClient: client,
	}
}

func (v *VersionSnapshotDynamoManager) GetVersionSnapshot(id ID) (*VersionSnapshot, error) {
	ctx := context.Background()
	var err error

	response, err := v.dynamoDbClient.GetItem(ctx, &dynamodb.GetItemInput{
		Key:       buildDynamoKeysAttr(id),
		TableName: aws.String(TableName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve [%v]: %v", id, err)
	}
	if len(response.Item) == 0 {
		// return empty document if item not found
		return nil, fmt.Errorf("entry not found: [%v]", id)
	}

	var targetDocument interface{}
	err = attributevalue.UnmarshalMap(response.Item, &targetDocument)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	// Cast the targetDocument to the expected type
	dynamoItem, err := convertDoc[dynamoVersionSnapshotItem](targetDocument)
	if err != nil {
		return nil, err
	}
	return &dynamoItem.VersionSnapshot, nil
}

// WriteVersionSnapshot writes a version snapshot to DynamoDB if it doesn't already exist.
func (v *VersionSnapshotDynamoManager) WriteVersionSnapshot(snapshot *VersionSnapshot) error {
	ctx := context.Background()
	dynamoItem := &dynamoVersionSnapshotItem{
		VersionSnapshot: *snapshot,
		PartitionKey:    getPartitionKeyFromVersionSnapshot(snapshot),
		SortKey:         getSortKeyFromVersionSnapshot(snapshot),
	}
	item, err := attributevalue.MarshalMapWithOptions(dynamoItem, func(options *attributevalue.EncoderOptions) {
		options.TagKey = "json"
	})
	if err != nil {
		return err
	}

	partitionKeyNotExists := expression.AttributeNotExists(expression.Name("partitionKey"))
	sortKeyNotExists := expression.AttributeNotExists(expression.Name("sortKey"))
	condition := partitionKeyNotExists.And(sortKeyNotExists)

	expr, err := expression.NewBuilder().WithCondition(condition).Build()
	if err != nil {
		return err
	}

	_, err = v.dynamoDbClient.PutItem(ctx, &dynamodb.PutItemInput{
		Item:                      item,
		TableName:                 aws.String(TableName),
		ConditionExpression:       expr.Condition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return err
	}
	return nil
}

func getPartitionKeyFromVersionSnapshot(v *VersionSnapshot) string {
	return v.AssetID
}

func getSortKeyFromVersionSnapshot(v *VersionSnapshot) string {
	return v.ResourceType + "^" + v.ResourceName + "^" + v.Location
}

func getPartitionKeyFromID(id ID) string {
	return id.AssetID
}

func getSortKeyFromID(id ID) string {
	return id.ResourceType + "^" + id.ResourceName + "^" + id.Location
}

// BuildDynamoKeysAttr builds DynamoDB keys for the given version snapshot.
func buildDynamoKeysAttr(id ID) map[string]types.AttributeValue {
	partitionKey, err := attributevalue.Marshal(getPartitionKeyFromID(id))
	if err != nil {
		panic(err)
	}
	sortKey, err := attributevalue.Marshal(getSortKeyFromID(id))
	if err != nil {
		panic(err)
	}
	return map[string]types.AttributeValue{"partitionKey": partitionKey, "sortKey": sortKey}
}

// convertDoc converts a single document of type interface{} to type T.
func convertDoc[T any](doc interface{}) (T, error) {
	// Marshal the document. Since we are starting with interface{},
	// the underlying type must be compatible with json.Marshal.
	var emptyDoc T
	marshal, err := json.Marshal(doc)
	if err != nil {
		return emptyDoc, fmt.Errorf("failed to marshal document: %w", err)
	}
	// Prepare a variable of type T to hold the unmarshalled data.
	var bd T
	// Unmarshal the JSON back into the specific type T.
	err = json.Unmarshal(marshal, &bd)
	if err != nil {
		return emptyDoc, fmt.Errorf("failed to unmarshal document into type T: %w", err)
	}
	// Return the type T
	return bd, nil
}
