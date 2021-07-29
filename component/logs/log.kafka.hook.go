package logs

import (
	"errors"
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"github.com/zhanghuizong/bitgame/service/config"
	"strconv"
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

	// 指针深拷贝
	dupEntry := entry.Dup()
	dupEntry.Level = entry.Level
	dupEntry.Caller = entry.Caller
	dupEntry.Message = entry.Message
	dupEntry.Buffer = entry.Buffer

	// 保证字段类型一致, 否则 ELK 会发生冲突
	// uid 转 string
	uid, uidExists := dupEntry.Data["uid"]
	if uidExists {
		uidInt, isInt := uid.(int)
		if isInt {
			dupEntry.Data["uid"] = strconv.Itoa(uidInt)
		}
	}

	// rid 转 string
	rid, ridExists := dupEntry.Data["rid"]
	if ridExists {
		ridInt, isInt := rid.(int)
		if isInt {
			dupEntry.Data["rid"] = strconv.Itoa(ridInt)
		}
	}

	// pid 转 string
	pid, pidExists := dupEntry.Data["pid"]
	if pidExists {
		pidInt64, isInt64 := pid.(int64)
		if isInt64 {
			dupEntry.Data["pid"] = strconv.FormatInt(pidInt64, 10)
		}

		pidInt, isInt := pid.(int)
		if isInt {
			dupEntry.Data["pid"] = strconv.Itoa(pidInt)
		}
	}

	// 追加日志时间
	dupEntry.Data["date"] = time.Now()

	// Format before writing
	b, err := hook.formatter.Format(dupEntry)
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
