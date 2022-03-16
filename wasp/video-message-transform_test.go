package wasp_test

import (
	"encoding/json"
	"testing"

	"github.com/digicatapult/wasp-simple-h264-payload-parser/wasp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransformVideoMessage(t *testing.T) {
	ingest := &wasp.IngestMessage{
		ThingID:   "883d818a-e115-4b99-b264-bb2033de21f1",
		Type:      "h264",
		Ingest:    "ingest-rtmp",
		IngestID:  "APP/STREAM",
		Payload:   "bytes",
		Timestamp: "2021-03-23 00:00:00", // iso date example
		Metadata:  map[string]interface{}{},
	}

	msgBytes, err := json.Marshal(ingest)
	require.NoError(t, err)

	outputBytes, err := wasp.TransformVideoMessages(msgBytes)
	require.NoError(t, err)

	output := &wasp.OutputMessage{}

	err = json.Unmarshal(outputBytes, output)
	require.NoError(t, err)

	assert.Equal(t, ingest.ThingID, output.ThingID)
	assert.Equal(t, ingest.Type, output.Type)
	assert.Equal(t, ingest.Payload, output.Value)
	assert.Equal(t, ingest.Timestamp, output.Timestamp)
	assert.Equal(t, ingest.Metadata, output.Metadata)
}
