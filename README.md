# dynamo-unit-test

If you want to unit tests for your GoLang+DynamoDB application,
you can start a local DynamoDB instance and directly execute CRUD operations on it.
More about [DynamoDB Local](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/DynamoDBLocal.html).

dynamo-unit-test proposes a structured way of testing DynamoDB CRUD operations, it

* Defines the data struct to be irrelevant to DynamoDB.
* Defines interfaces for CRUD operations, which are also irrelevant to DynamoDB.
* Defines a DynamoDB implementation of the CRUD operations.
* Implements unit tests for DynamoDB implementation, which uses a local DynamoDB instance, without mocking 
  the DynamoDB client.
* Enables CI to run unit tests.

To run the tests, simply checkout the repository and run the following command:

```bash
make test
```

To clean up the local DynamoDB instance, run:

```bash
make clean
```