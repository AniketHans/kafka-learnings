package main

import (
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	KAFKA_SERVER = "localhost:9092"
	KAFKA_TOPIC  = "english-words"
	GROUP_ID = "Grp1"
)
func main() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":KAFKA_SERVER,
		"group.id":GROUP_ID,
		"auto.offset.reset":    "smallest",
	})
	if err != nil {
	panic(err)
	}
	defer c.Close()
	topic := KAFKA_TOPIC
 	c.SubscribeTopics([]string{topic}, nil)
	for {
		ev := c.Poll(100)
		switch e:= ev.(type){
		case *kafka.Message:
				c.Commit()
				msgProcess(e)
		case kafka.PartitionEOF:
			fmt.Printf("%% Reached %v\n", e)
		case kafka.Error:
			fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
		default:
		}
	}


}

func msgProcess(msg *kafka.Message){
	fmt.Println("Haha!!! Got the message", string(msg.Value))

}
