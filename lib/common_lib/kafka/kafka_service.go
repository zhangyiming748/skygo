package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/wvanbergen/kafka/consumergroup"
	"github.com/wvanbergen/kazoo-go"
	"strings"
	"sync"
	"time"
)

var (
	defaultProducer *ProducerEngine
	kafkaMu         sync.Mutex
)

func NewConsumerEngine(brokers []string, version, consumerGroup, zookeeperNodes string) *ConsumerEngine {
	engine := &ConsumerEngine{
		interruptChan:  make(chan int),
		closerChan:     make(chan int),
		brokers:        brokers,
		version:        version,
		consumerGroup:  consumerGroup,
		zookeeperNodes: zookeeperNodes,
	}
	return engine
}

type ConsumerEngine struct {
	interruptChan  chan int // 用于通知消费者，收到interrupt信号
	closerChan     chan int // 用于协程之间告知“消费者协程”全部平滑停止完毕
	brokers        []string
	version        string
	consumerGroup  string
	zookeeperNodes string
}

type InterfaceConsumer interface {
	DeliveryMsg(string, string, int64)
}

func (e *ConsumerEngine) StartConsumer(kafkaConsumer InterfaceConsumer, topic string) error {
	version, err := sarama.ParseKafkaVersion(e.version)
	if err != nil {
		return err
	}
	if version.IsAtLeast(sarama.V0_9_0_0) {
		e.startClusterConsumer(kafkaConsumer, topic, e.interruptChan, e.closerChan)
	} else {
		e.startWvanbergenConsumer(kafkaConsumer, topic, e.interruptChan, e.closerChan)
	}

	return nil
}

func (e *ConsumerEngine) Interrupt() {
	close(e.interruptChan)
}

func (e *ConsumerEngine) GetInterrupt() chan int {
	return e.interruptChan
}

func (e *ConsumerEngine) GetCloserChan() chan int {
	return e.closerChan
}

// 消费者组支持
// kafka consumer <= 0.82, 使用Wvanbergen
func (e *ConsumerEngine) startWvanbergenConsumer(kafkaConsumer InterfaceConsumer, topic string, InterruptChan, CloserChan chan int) {
	conf := consumergroup.NewConfig()
	conf.Offsets.ProcessingTimeout = 10 * time.Second
	conf.Offsets.Initial = sarama.OffsetOldest

	var zookeeperNodes []string
	zookeeperNodes, conf.Zookeeper.Chroot = kazoo.ParseConnectionString(e.zookeeperNodes)

	kafkaTopics := strings.Split(topic, ",")

	consumerG, consumerErr := consumergroup.JoinConsumerGroup(fmt.Sprintf("%s_%s", e.consumerGroup, topic), kafkaTopics, zookeeperNodes, conf)
	if consumerErr != nil {
		fmt.Println("consumer start failed:" + fmt.Sprint(consumerErr))
	}

	defer consumerG.Close()

	go func() {
		<-InterruptChan
		fmt.Println("consumer get the interrupt signal")

		if err := consumerG.Close(); err != nil {
			fmt.Println("fail to close the consumer:" + fmt.Sprint(err))
		} else {
			close(CloserChan)
		}
	}()

	go func() {
		for err := range consumerG.Errors() {
			fmt.Println("consumer get Errors:", fmt.Sprint(err))
		}
	}()

	eventCount := 0
	offsets := make(map[string]map[int32]int64)

	for message := range consumerG.Messages() {
		if offsets[message.Topic] == nil {
			offsets[message.Topic] = make(map[int32]int64)
		}

		eventCount++

		if offsets[message.Topic][message.Partition] != 0 && offsets[message.Topic][message.Partition] != message.Offset-1 {
			s := fmt.Sprintf("Unexpected offset on %s:%d. Expected %d, found %d, diff %d.\n", message.Topic, message.Partition, offsets[message.Topic][message.Partition]+1, message.Offset, message.Offset-offsets[message.Topic][message.Partition]+1)
			fmt.Println(s)
		}

		// Simulate processing time
		kafkaConsumer.DeliveryMsg(message.Topic, string(message.Value), 0)

		//消费记录日志
		info := fmt.Sprintf("Consume message topic-partition-offset：  %s %d %d", message.Topic, message.Partition, message.Offset) //message.Value
		fmt.Println(info)

		offsets[message.Topic][message.Partition] = message.Offset
		consumerG.CommitUpto(message)
	}
}

// 消费者组支持
// kafka consumer >= 0.9 使用 sarama-cluster
func (e *ConsumerEngine) startClusterConsumer(kafkaConsumer InterfaceConsumer, topic string, InterruptChan, CloserChan chan int) {
	conf := cluster.NewConfig()
	conf.Group.Mode = cluster.ConsumerModePartitions
	conf.Consumer.Offsets.Initial = sarama.OffsetOldest
	// 新下载的sarama包，里面跟老版本相比有些修改，导致sarama包里面没有对如下CommitInterval做了配置，默认0s，sarama-cluster包会报错
	// 所以这里我们要设定好默认值，默认1s
	conf.Consumer.Offsets.CommitInterval = 1 * time.Second

	kafkaTopics := strings.Split(topic, ",")
	consumerG, consumerErr := cluster.NewConsumer(e.brokers, fmt.Sprintf("%s_%s", e.consumerGroup, topic), kafkaTopics, conf)
	if consumerErr != nil {
		fmt.Println("consumer start failed:" + fmt.Sprint(consumerErr))
	}

	defer consumerG.Close()

	go func() {
		<-InterruptChan
		if err := consumerG.Close(); err != nil {
			fmt.Println("fail to close the consumer:" + fmt.Sprint(err))
		} else {
			close(CloserChan)
		}
	}()

	go func() {
		for err := range consumerG.Errors() {
			fmt.Println("consumer get Errors:", fmt.Sprint(err))
		}
	}()

	for part := range consumerG.Partitions() {
		go func(pc cluster.PartitionConsumer) {
			for message := range pc.Messages() {
				kafkaConsumer.DeliveryMsg(message.Topic, string(message.Value), 0)
				consumerG.MarkOffset(message, "")
			}
		}(part)
	}
}

// ------------------------------------------------------------
func NewSyncProducer(brokers []string, version string) (sarama.SyncProducer, error) {
	return getSyncProducer(brokers, version)
}

// ------------------------------------------------------------
func NewProduceEngine(brokers []string, version string) (*ProducerEngine, error) {
	var err error
	producer := new(ProducerEngine)
	// --------------- 同步生产---------------------------
	if producer.syncProducer, err = getSyncProducer(brokers, version); err != nil {
		return nil, err
	}

	// --------------- 异步生产， 并发量大时，必须采用这种方式 ---------------------------
	if producer.asyncProducer, err = getAsyncProducer(brokers, version); err != nil {
		return nil, err
	}

	return producer, nil
}

type ProducerEngine struct {
	syncProducer  sarama.SyncProducer
	asyncProducer sarama.AsyncProducer
}

// 生产消息， 同步
func (k *ProducerEngine) Produce(strBody, topicName string) (partition int32, offset int64, err error) {
	message := getProducerMessage(strBody, topicName)
	return k.syncProducer.SendMessage(message)
}

// 生产消息, 同步， 按照消息key分partition
func (k *ProducerEngine) ProduceSequential(strBody, topicName string, key string) (partition int32, offset int64, err error) {
	message := getProducerMessage(strBody, topicName, key)
	return k.syncProducer.SendMessage(message)
}

// 生产消息, 异步
func (k *ProducerEngine) AsyncProducer(strBody, topicName string) {
	message := getProducerMessage(strBody, topicName)
	k.asyncProducer.Input() <- message
}

// 生产消息, 异步, 按照消息key分partition
func (k *ProducerEngine) AsyncProducerSequential(strBody, topicName string, key string) {
	message := getProducerMessage(strBody, topicName, key)
	k.asyncProducer.Input() <- message
}

// 异步生产的success管道
func (k *ProducerEngine) GetAsyncProducerSuccessChan() <-chan *sarama.ProducerMessage {
	return k.asyncProducer.Successes()
}

// 异步生产的success管道
func (k *ProducerEngine) GetAsyncProducerErrorsChan() <-chan *sarama.ProducerError {
	k.asyncProducer.Close()
	return k.asyncProducer.Errors()
}

// 关闭消费
func (k *ProducerEngine) Close() error {
	if err := k.asyncProducer.Close(); err != nil {
		return err
	}
	if err := k.asyncProducer.Close(); err != nil {
		return err
	}
	return nil
}

// 同步生产
func getSyncProducer(brokers []string, version string) (sarama.SyncProducer, error) {
	kafkaVersion, err := sarama.ParseKafkaVersion(version)
	if err != nil {
		return nil, err
	}

	// For the data collector, we are looking for strong consistency semantics.
	// Because we don't change the flush settings, sarama will try to produce messages
	// as fast as possible to keep latency low.
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Version = kafkaVersion

	if syncProducer, err := sarama.NewSyncProducer(brokers, config); err != nil {
		return nil, err
	} else {
		return syncProducer, nil
	}
}

// 异步生产
func getAsyncProducer(brokers []string, version string) (sarama.AsyncProducer, error) {
	kafkaVersion, err := sarama.ParseKafkaVersion(version)
	if err != nil {
		return nil, err
	}

	// For the access log, we are looking for AP semantics, with high throughput.
	// By creating batches of compressed messages, we reduce network I/O at a cost of more latency
	config := sarama.NewConfig()

	config.Producer.RequiredAcks = sarama.WaitForLocal       //Only wait for the leader to ack
	config.Producer.Compression = sarama.CompressionSnappy   //Compress messages
	config.Producer.Flush.Frequency = 500 * time.Millisecond //Flush batches every 500ms
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Version = kafkaVersion

	if asyncProducer, err := sarama.NewAsyncProducer(brokers, config); err != nil {
		return nil, err
	} else {
		// 异步消息生产的使用者要自己去获取响应结果
		go func() {
			for {
				select {
				case <-asyncProducer.Successes():
				case fail := <-asyncProducer.Errors():
					if fail != nil {
						fmt.Println(fail.Error())
					}
				}
			}
		}()
		return asyncProducer, nil
	}
}

// 生产消息构建
func getProducerMessage(strBody, topicName string, key ...string) *sarama.ProducerMessage {
	if len(key) > 0 {
		return &sarama.ProducerMessage{
			Topic: topicName,
			Value: sarama.ByteEncoder(strBody),
			Key:   sarama.ByteEncoder(strings.Join(key, "")),
		}
	} else {
		return &sarama.ProducerMessage{
			Topic: topicName,
			Value: sarama.ByteEncoder(strBody),
		}
	}
}

// -------------------------------------------------
// 检验连通性
func CheckConnection(brokers []string, version string) error {
	config := sarama.NewConfig()
	config.Net.DialTimeout = time.Second * 5
	config.Net.ReadTimeout = time.Second * 5
	config.Admin.Timeout = time.Second * 5
	config.Admin.Retry.Max = 1
	config.Metadata.Timeout = time.Second * 5
	config.Metadata.Retry.Max = 1

	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		return err
	}
	client.Close()
	return nil
}
