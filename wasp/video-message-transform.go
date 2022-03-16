package wasp

// IngestMessage defines the video ingest message structure
type IngestMessage struct {
	Ingest    string                 `json:"ingest"`
	IngestID  string                 `json:"ingestId"`
	Timestamp string                 `json:"timestamp"`
	Payload   string                 `json:"payload"`
	Metadata  map[string]interface{} `json:"metadata"`
	ThingID   string                 `json:"thingId"`
	Type      string                 `json:"type"`
}

// OutputMessage defines the output message structure
type OutputMessage struct {
	ThingID   string
	Type      string
	Timestamp string
	Value     string
	Metadata  string
}

// TransformVideoMessages processes the messages as part of the message relay
func TransformVideoMessages(msg []byte) ([]byte, error) {
	// Transform messages here
	// Unmarshal
	// Modify values
	// Create new message
	// Marshal
	// return
	return msg, nil
}
