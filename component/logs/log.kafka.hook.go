package logs

import (
	"errors"
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"github.com/zhanghuizong/bitgame/service/config"
	"time"
)

type KafkaHook struct {
	// Id of the hook
	id string

	// Log levels allowed
	levels []logrus.Level

	// Log entry formatter
	formatter logrus.Formatter

	// sarama.AsyncProducer
	producer sarama.AsyncProducer
}

// initKafkaHook Create a new KafkaHook.
func initKafkaHook() (*KafkaHook, error) {
	brokers := config.GetKafkaBrokers()
	if len(brokers) <= 0 {
		return nil, errors.New("kafka 集群列表为空")
	}

	// kafka 异常
	sarama.PanicHandler = func(data interface{}) {
		logrus.Errorf("kafka 处理过程异常. %v", data)
	}

	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForLocal
	kafkaConfig.Producer.Compression = sarama.CompressionSnappy
	kafkaConfig.Producer.Flush.Frequency = 500 * time.Millisecond

	producer, err := sarama.NewAsyncProducer(brokers, kafkaConfig)
	if err != nil {
		return nil, err
	}

	go func() {
		for err := range producer.Errors() {
			logrus.Errorf("kafka 发生异常. %v", err)
		}
	}()

	hook := &KafkaHook{
		"kafka.hook",
		[]logrus.Level{
			logrus.TraceLevel,
			logrus.DebugLevel,
			logrus.InfoLevel,
			logrus.WarnLevel,
			logrus.ErrorLevel,
		},
		getLogstashFormatter(),
		producer,
	}

	return hook, nil
}

func (hook *KafkaHook) Id() string {
	return hook.id
}

func (hook *KafkaHook) Levels() []logrus.Level {
	return hook.levels
}

func (hook *KafkaHook) Fire(entry *logrus.Entry) error {
	var topics = config.GetKafkaTopics()
	if len(topics) <= 0 {
		return errors.New("kafka topic 没有配置")
	}

	// 追加日志时间
	entry.Data["date"] = time.Now().Format("2006-01-02 15:04:05.000")

	// Format before writing
	b, err := hook.formatter.Format(entry)
	if err != nil {
		return err
	}

	value := sarama.ByteEncoder(b)
	for _, topic := range topics {
		hook.producer.Input() <- &sarama.ProducerMessage{
			Key:   sarama.StringEncoder(topic),
			Topic: topic,
			Value: value,
		}
	}

	return nil
}
