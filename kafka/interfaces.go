package kafka

// Producer defines the producer operations
type Producer interface {
	SendMessage(topic, msg string, value []byte) error
}

// Consumer defines the consumer operations
type Consumer interface {
	Listen(topic string, received chan []byte)
}

// IngestMessage defines the video ingest message structure
type IngestMessage struct {
	Ingest    string `json:"ingest"`
	IngestID  string `json:"ingestId"`
	Timestamp string `json:"timestamp"`
	Payload   string `json:"payload"`
	Metadata  string `json:"metadata"`
}

// OutputMessage defines the output message structure
type OutputMessage struct {
	ThingID   string
	Type      string
	Timestamp string
	Value     string
	Metadata  string
}
