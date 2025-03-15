package main

import (
	"context"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	KAFKA_SERVER = "localhost:9092"
	KAFKA_TOPIC_1  = "hello-producer-1"
	KAFKA_TOPIC_2  = "hello-producer-2"
	TRANSACTION_ID = "Hello-Producer-Trans"
)


func main() {
	ctx := context.Background()
	fmt.Println("Creating Kafka Producer....")
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": KAFKA_SERVER,
		"transactional.id": TRANSACTION_ID,  // idempotence is automatically enabled incase of transaction id, transactional.id should be unique for each producer instance
	})
	if err!=nil{
		panic(err)
	}
	defer p.Close()

	p.InitTransactions(ctx)
	fmt.Println("Sending data to brokers....")
	fmt.Println("First transaction start..............")
	p.BeginTransaction()
	topic1:=KAFKA_TOPIC_1
	topic2:=KAFKA_TOPIC_2
	for i:=1;i<=5;i++{
		err := p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic1, Partition: int32(kafka.PartitionAny) }, Value: []byte(fmt.Sprintf("Simple-value-T1-%v",i)),
		},nil)

		if err !=nil{
			fmt.Println("Error in sending records")
			fmt.Println("Aborting First Transaction.....")
			p.AbortTransaction(ctx)
			p.Close()
			panic(err)
		}

		err = p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic2, Partition: int32(kafka.PartitionAny) }, Value: []byte(fmt.Sprintf("Simple-value-T1-%v",i)),
		},nil)

		if err !=nil{
			fmt.Println("Error in sending records")
			fmt.Println("Aborting First Transaction.....")
			p.AbortTransaction(ctx)
			p.Close()
			panic(err)
		}
	}
	fmt.Println("Committing First Transaction.....")
	p.CommitTransaction(ctx)


	fmt.Println("Second transaction start..............")
	p.BeginTransaction()
	for i:=1;i<=5;i++{
		err := p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic1, Partition: int32(kafka.PartitionAny) }, Value: []byte(fmt.Sprintf("Simple-value-T2-%v",i)),
		},nil)

		if err !=nil{
			fmt.Println("Error in sending records")
			fmt.Println("Aborting Second Transaction.....")
			p.AbortTransaction(ctx)
			p.Close()
			panic(err)
		}

		err = p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic2, Partition: int32(kafka.PartitionAny) }, Value: []byte(fmt.Sprintf("Simple-value-T2-%v",i)),
		},nil)

		if err !=nil{
			fmt.Println("Error in sending records")
			fmt.Println("Aborting Second Transaction.....")
			p.AbortTransaction(ctx)
			p.Close()
			panic(err)
		}
	}
	fmt.Println("Aborting Second Transaction explicitly.....") /// Thus, the messages sent by this transaction will not be received.
	p.AbortTransaction(ctx)





	
}