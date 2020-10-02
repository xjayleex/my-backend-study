package main
import (
	"fmt"
	"github.com/Shopify/sarama"
)
type Producer struct {
	sProducer  sarama.SyncProducer
	asProducer sarama.AsyncProducer
}

func main() {
	config :=  sarama.NewConfig()
	config.Producer.Partitioner =
		sarama.NewRoundRobinPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	p, err := sarama.NewSyncProducer([]string{
		"master:9092",
		"node2:9092",
		"node4:9092" },config)
	if err != nil {
		panic(err)
	}

	part, offset, err := p.SendMessage(&sarama.ProducerMessage{
		Topic: "simple",
		Value: sarama.StringEncoder("Hello"),
	})

	if err != nil {
		panic(err)
	}
	fmt.Printf("%d/%d\n", part, offset)
}