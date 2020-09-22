package main
import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"fmt"
	"os"
	"strconv"
)

func main() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String("ap-northeast-2"),
		},
		SharedConfigState: session.SharedConfigEnable,
		Profile:           "jay",
	}))
	svc := dynamodb.New(sess)

	item := &Item {
		Year: 2015,
		Title: "Interstella",
		Plot: "n 2067,[5] crop blights and dust storms threaten humanity's survival. Corn is the last viable crop. The world also regresses into a post-truth society where younger generations are taught false history, including the faking of the Apollo moon missions. Widowed engineer and former NASA pilot Joseph Cooper is now a farmer. Living with him are his father-in-law, Donald; his 15-year-old son, Tom Cooper, and 10-year-old daughter, Murphy \"Murph\" Cooper. After a dust storm, strange dust patterns inexplicably appear on Murphy's bedroom floor; she attributes the anomaly to a ghost. Cooper eventually deduces the patterns were caused by gravity variations and that they represent geographic coordinates in binary code. Cooper follows the coordinates to a secret NASA facility headed by Professor John Brand, Cooper's former supervisor. Professor Brand says gravitational anomalies have happened elsewhere. Forty-eight years earlier, unknown beings positioned a wormhole near Saturn, opening a path to a distant galaxy with twelve potentially habitable worlds located near a black hole named Gargantua. Twelve volunteers traveled through the wormhole to individually survey the planets. Astronauts Miller, Edmunds, and Mann reported positive results. Based on their data, Professor Brand conceived two plans to ensure humanity's survival. Plan A involves developing a gravitational propulsion theory to propel colonies into space, while Plan B involves launching the Endurance spacecraft carrying 5,000 frozen human embryos to colonize a habitable planet.\n\nCooper is recruited to pilot the Endurance. The crew includes scientists Dr. Amelia Brand (Professor Brand's daughter), Dr. Romilly, Dr. Doyle, and robots TARS and CASE. Before leaving, Cooper gives a distraught Murphy his wristwatch to compare their relative time for when he returns. After traversing the wormhole, Romilly studies the black hole while Cooper, Doyle, and Brand descend in a landing craft to investigate Miller's planet, an ocean world. After finding wreckage from Miller's ship, a gigantic tidal wave kills Doyle and delays the lander's departure. Due to the proximity of the black hole, time is severely dilated: as a result, 23 years have elapsed for Romilly on Endurance by the time Cooper and Brand return.",
		Rating: 9.5,
	}
	attr, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		fmt.Println("Got error mashalling new movie item:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	input := &dynamodb.PutItemInput {
		Item: attr,
		TableName: aws.String(TableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println("Got error calling PutItem : ")
		fmt.Println(err.Error)
		os.Exit(1)
	}
	year := strconv.Itoa(item.Year)

	fmt.Println("Successfully added '" + item.Title + "' (" + year + ") to table " + TableName)
}