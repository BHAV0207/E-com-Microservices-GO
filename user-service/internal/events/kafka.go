package events

import (
	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	Writer *kafka.Writer
	//kafka.Writer is the object responsible for sending messages to a Kafka topic.
}

// NewKafkaProducer creates a new KafkaProducer with the given broker and topic.
func KafkaWriter(brokerURL, topic string) *KafkaProducer {
	return &KafkaProducer{
		Writer: &kafka.Writer{
			Addr:         kafka.TCP(brokerURL), // Connect to the Kafka broker
			Topic:        topic,                // The topic to publish to
			Balancer:     &kafka.LeastBytes{},  // Choose the partition with least data
			RequiredAcks: kafka.RequireAll,     // Wait for all replicas to confirm
		},
	}
}

type KafkaConsumer struct {
	Reader *kafka.Reader
}

func KafkaReader(brokerURL, topic, groupID string) *KafkaConsumer {
	return &KafkaConsumer{
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{brokerURL},
			Topic:   topic,
			GroupID: groupID,
		}),
	}
}

// | Symbol | Meaning                                                                      | Analogy                                                      |
// | ------ | ---------------------------------------------------------------------------- | ------------------------------------------------------------ |
// | `&`    | “Get the **address** of” a variable                                          | Think of it as writing someone's **home address** on paper   |
// | `*`    | “Access the **value stored at** that address” (or define a **pointer type**) | Think of it as **going to the house** and meeting the person |

/*
x := 10       // x holds the value 10
p := &x       // p holds the address of x (like a pointer to x)
fmt.Println(p)  // prints something like 0xc0000140a8 (the memory address)
fmt.Println(*p) // prints 10 (value at that address)

📘 Analogy:
x is the actual data (a value stored in memory).
p is the pointer, which tells you where that data lives.
*p means “go to that address and read the value”.*/
