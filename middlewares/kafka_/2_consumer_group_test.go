package kafka_

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"log"
)

// https://www.cnblogs.com/payapa/p/15401357.html

// 自定义消费者处理程序
type MyConsumerHandler struct{}

// Setup 在每个消费者协程启动前调用，用以准备初始化工作
func (h MyConsumerHandler) Setup(session sarama.ConsumerGroupSession) error {
	fmt.Println("Setup")
	return nil
}

// Cleanup 在每个消费者协程退出后调用，用以清理资源
func (h MyConsumerHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	fmt.Println("Cleanup")
	return nil
}

// ConsumeClaim 消费消息的实际逻辑，在每个分配到的 claim 上执行
func (h MyConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		fmt.Printf("Message claimed: key = %s, value = %s, topic = %s, partition = %d, offset = %d\n",
			string(message.Key), string(message.Value), message.Topic, message.Partition, message.Offset)
		session.MarkMessage(message, "")
	}
	return nil
}

// 消费者组
func SaramaConsumerGroup(addr []string) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.V3_4_0_0                                                                             // 指定kafka 版本
	config.Consumer.Offsets.Initial = sarama.OffsetOldest                                                        // 未找到组消费位移的时候从哪边开始消费
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()} // 设置rebalance 策略
	group, err := sarama.NewConsumerGroup(addr, "my-group", config)
	if err != nil {
		panic(err)
	}
	defer func() { _ = group.Close() }()

	// Track errors
	go func() {
		for err := range group.Errors() {
			log.Println(err)
		}
	}()
	fmt.Println("Consumed start")
	// Iterate over consumer sessions.
	ctx := context.Background()
	for {
		topics := []string{"my_topic"}
		handler := MyConsumerHandler{}

		// `Consume` should be called inside an infinite loop, when a
		// server-side rebalance happens, the consumer session will need to be
		// recreated to get the new claims
		err := group.Consume(ctx, topics, handler)
		if err != nil {
			panic(err)
		}
	}
}
