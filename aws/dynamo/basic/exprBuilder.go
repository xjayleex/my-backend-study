package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"os"
)

func CheckError (err error, msg string) {
	if err != nil {
		fmt.Println()
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func main ()  {
	// 자격 증
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String("ap-northeast-2"),
		},
		SharedConfigState: session.SharedConfigEnable,
		Profile:           "jay",
	}))
	svc := dynamodb.New(sess)
	// 테이블 항목의 minRating 및 연도에 대한 변수 생성
	minRating := 3.0
	year := 2015
	// 검색할 항목을 필터링할 연도를 정의하는 expression 작성
	filt := expression.Name("Year").Equal(expression.Value(year))
	// 검색된 항목의 연도와 rating을 가져오도록 하는 프로젝션 생성
	proj := expression.NamesList(expression.Name("Title"),
		expression.Name("Year"), expression.Name("Rating"))
	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	CheckError(err,"Got error building expression:")

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames: 		expr.Names(),
		ExpressionAttributeValues: 		expr.Values(),
		FilterExpression: 				expr.Filter(),
		ProjectionExpression: 			expr.Projection(),
		TableName: 						aws.String(TableName),
	}
	// Scan 입력으로 작성한 파라미터에 대해 Scan 호출
	result, err := svc.Scan(params)
	CheckError(err,"Query API call failed:")

	numItems := 0

	for _, i := range result.Items {
		item := Item{}

		err = dynamodbattribute.UnmarshalMap(i, &item)
		CheckError(err, "Got error unmarshalling:")

		if item.Rating > minRating {

			numItems++

			fmt.Println("Title: ", item.Title)
			fmt.Println("Rating:", item.Rating)
			fmt.Println()
		}
	}

	fmt.Println("Found", numItems, "movie(s) with a rating above", minRating, "in", year)
}