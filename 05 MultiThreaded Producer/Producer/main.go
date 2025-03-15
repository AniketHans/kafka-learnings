package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	KAFKA_SERVER = "localhost:9092"
	KAFKA_TOPIC  = "csv-data"
)

var wg sync.WaitGroup

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
	for _, fileName := range([]string{"../DataFiles/file1.csv","../DataFiles/file2.csv"}){
		wg.Add(1)
		go sendDataToBroker(p,topic,fileName)
	}
	wg.Wait()
	fmt.Println("Data sent to broker")
}

func sendDataToBroker(p *kafka.Producer, topic string, fileName string){
	file, err := os.Open(fileName)
	if err!=nil{
		fmt.Println("Error in reading file")
		panic(err)
	}
	defer file.Close()
	defer wg.Done()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err!=nil{
		fmt.Println("Error in reading records")
		panic(err)
	}
	for _, eachRecord := range records{
		err := p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: int32(kafka.PartitionAny) }, Value: []byte(fmt.Sprintf("%v",eachRecord)),
		},nil)

		if err !=nil{
			fmt.Println("Error in sending records")
			panic(err)
		}
	}

}