package main
import (
	"github.com/aws/aws-sdk-go/aws"
	_ "github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"fmt"
	"os"
)
func main(){
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config : aws.Config{
			Region: aws.String("ap-northeast-2"),
		},
		SharedConfigState: session.SharedConfigEnable,
		Profile: "jay",
	}))
	svc := dynamodb.New(sess)
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Year"),
				AttributeType: aws.String("N"),
			},
			{
				AttributeName: aws.String("Title"),
				AttributeType: aws.String("S"),
			},
		},
		BillingMode: aws.String(dynamodb.BillingModePayPerRequest),
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Year"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("Title"),
				KeyType:       aws.String("RANGE"),
			},
		},
		/*ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},*/
		TableName: aws.String("testtab"),
	}
	_, err := svc.CreateTable(input)
	if err != nil {
		fmt.Println("Got error calling CreateTable : ")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("Created the table", "testtab")
}
