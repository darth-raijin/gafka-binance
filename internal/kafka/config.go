package kafka

type KafkaConfig struct {
	Brokers          []string       `yaml:"brokers"`
	Topics           []string       `yaml:"topics"`
	ConsumerSettings ConsumerConfig `yaml:"consumer"`
	ProducerSettings ProducerConfig `yaml:"producer"`
	MaxPollRecords   int            `yaml:"max_poll_records"`
	SessionTimeoutMs int            `yaml:"session_timeout_ms"`
	RequestTimeoutMs int            `yaml:"request_timeout_ms"`
}

type ConsumerConfig struct {
	GroupID          string `yaml:"group_id"`
	AutoOffsetReset  string `yaml:"auto_offset_reset"`
	EnableAutoCommit bool   `yaml:"enable_auto_commit"`
}

type ProducerConfig struct {
	Acks           string `yaml:"acks"`
	Retries        int    `yaml:"retries"`
	RetryBackoffMs int    `yaml:"retry_backoff_ms"`
}
