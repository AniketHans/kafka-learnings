package main

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	KAFKA_SERVER = "localhost:9092"
	KAFKA_TOPIC  = "english-words"
)

func main() {
	fmt.Println("Creating Kafka Producer....")
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": KAFKA_SERVER,
	})
	if err!=nil{
		panic(err)
	}
	defer p.Close()

	topic := KAFKA_TOPIC
	fmt.Println("Sending data to brokers....")
	for _, val := range([]string{"Helloyoyo","HIIIIyyo","Kaise HO yoyo"}){
		err := p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: int32(kafka.PartitionAny) }, Value: []byte(val),
		},nil)

		if err !=nil{
			panic(err)
		}
	}
	fmt.Println("Data sent to broker")
}