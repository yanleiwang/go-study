package kafka_

import (
	"github.com/IBM/sarama"
	"log"
	"os"
	"os/signal"
	"sync"
	"testing"
	"time"
)

func TestAsync(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		asyncProducer(Address)
	}()

	go func() {
		defer wg.Done()
		SaramaConsumer(Address)
	}()
	wg.Wait()
}

// 异步生产者，  消息send 不需要等待ack
func asyncProducer(address []string) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	// 默认就是enabled
	//config.Producer.Return.Errors = true
	producer, err := sarama.NewAsyncProducer(address, config)
	if err != nil {
		log.Fatalln(err)
	}

	// Trap SIGINT to trigger a graceful shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	var wg sync.WaitGroup
	var (
		enqueued, successes, producerErrors int
	)
	wg.Add(2)
	go func() {
		defer wg.Done()
		for range producer.Successes() {
			successes++
		}
	}()

	go func() {
		defer wg.Done()
		for err := range producer.Errors() {
			log.Println(err)
			producerErrors++
		}
	}()

	timer := time.NewTimer(time.Second * 3)
Loop:
	for {
		msg := &sarama.ProducerMessage{
			Topic: "my-topic",
			Value: sarama.StringEncoder("testing 123"),
		}
		select {
		case producer.Input() <- msg:
			enqueued++
		case <-timer.C:
			producer.AsyncClose() // Trigger a shutdown of the producer.
			break Loop
		}
		time.Sleep(time.Millisecond * 100)
	}
	wg.Wait()
	log.Printf("Successfully produced: %d; errors: %d\n", successes, producerErrors)

}
