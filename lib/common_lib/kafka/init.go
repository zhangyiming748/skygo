package kafka

import (
	"skygo_detection/service"

	"fmt"
)

// consume file程序使用
var kafkaProduce *ProducerEngine

func InitKafkaProduce() error {
	config := service.LoadConfig().Kafka
	brokers := config.Brokers
	version := config.Version
	if engine, err := NewProduceEngine(brokers, version); err != nil {
		fmt.Println(err)
		return err
	} else {
		kafkaProduce = engine
		return nil
	}
}

func GetKafkaProduce() *ProducerEngine {
	if kafkaProduce == nil {
		InitKafkaProduce()
	}
	return kafkaProduce
}

// -----------------------------------------------------
// consume kafka程序使用的
func GetKafkaConsumerCK() *ConsumerEngine {
	config := service.LoadConfig().Kafka
	brokers := config.Brokers
	version := config.Version
	group := config.ConsumerGroupPrefix
	engine := NewConsumerEngine(brokers, version, group, "")
	return engine
}

var kafkaProduceCK *ProducerEngine

func InitKafkaProduceCK() error {
	config := service.LoadConfig().Kafka
	brokers := config.Brokers
	version := config.Version
	if engine, err := NewProduceEngine(brokers, version); err != nil {
		fmt.Println(err)
		return err
	} else {
		kafkaProduceCK = engine
		return nil
	}
}

func GetKafkaProduceCK() *ProducerEngine {
	if kafkaProduceCK == nil {
		InitKafkaProduceCK()
	}
	return kafkaProduceCK
}
