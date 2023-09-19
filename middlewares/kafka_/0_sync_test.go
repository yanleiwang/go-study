package kafka_

import (
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"os"
	"os/signal"
	"sync"
	"testing"
	"time"
)

var Address = []string{"127.0.0.1:9092"}

const (
	Topic = "my-topic"
)

func TestSync(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		syncProducer(Address)
	}()

	go func() {
		defer wg.Done()
		SaramaConsumer(Address)
	}()
	wg.Wait()
}

// 同步消息模式:  消息send 需要等待ack
func syncProducer(address []string) {
	// 配置
	config := sarama.NewConfig()
	// 属性设置
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Timeout = 5 * time.Second
	// 创建生成者
	p, err := sarama.NewSyncProducer(address, config)
	// 判断错误
	if err != nil {
		log.Fatalln(err)
	}
	// 最后关闭生产者
	defer func() {
		if err := p.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	// 消息
	srcValue := "sync: this is a message. index=%d"
	// 循环发消息
	for i := 0; i < 10; i++ {
		// 格式化消息
		value := fmt.Sprintf(srcValue, i)
		// 创建消息
		msg := &sarama.ProducerMessage{
			Topic: Topic,
			Value: sarama.StringEncoder(value),
		}
		// 发送消息
		part, offset, err := p.SendMessage(msg)
		if err != nil {
			log.Printf("send message(%s) err=%s \n", value, err)
		} else {
			log.Printf("%s, 发送成功，partition=%d, offset=%d \n", value, part, offset)
		}
		// 每隔两秒发送一个消息
		time.Sleep(1 * time.Second)
	}
}

// 消费者
func SaramaConsumer(address []string) {
	consumer, err := sarama.NewConsumer(address, nil)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition(Topic, 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	consumed := 0

	timer := time.NewTimer(time.Second * 3)
Loop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("Consumed message offset %d\n", msg.Offset)
			consumed++
		case <-timer.C:
			break Loop
		}
	}

	log.Printf("Consumed: %d\n", consumed)
}
