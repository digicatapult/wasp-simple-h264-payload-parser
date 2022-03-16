package wasp_test

import (
	"encoding/json"
	"testing"

	"github.com/digicatapult/wasp-simple-h264-payload-parser/wasp"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransformVideoMessage(t *testing.T) {
	ingest := &wasp.IngestMessage{
		Ingest:    uuid.NewV4().String(),
		IngestID:  uuid.NewV4().String(),
		Timestamp: uuid.NewV4().String(),
		Payload:   uuid.NewV4().String(),
		ThingID:   uuid.NewV4().String(),
		Type:      uuid.NewV4().String(),
		Metadata: map[string]interface{}{
			uuid.NewV4().String(): uuid.NewV4().String(),
		},
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
	assert.Equal(t, ingest.ThingID, output.ThingID)
	assert.Equal(t, ingest.Metadata, output.Metadata)
}
