package handler

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Tencent/WeKnora/internal/types"
	"github.com/stretchr/testify/require"
)

func TestEnqueueKnowledgeListDeletePreservesLanguage(t *testing.T) {
	q := &documentDeleteEnqueuer{}
	h := &KnowledgeHandler{asynqClient: q}
	ctx := context.WithValue(context.Background(), types.LanguageContextKey, "ru-RU")
	_, err := h.enqueueKnowledgeListDelete(ctx, 7, "kb", []string{"doc"})
	require.NoError(t, err)
	var payload map[string]interface{}
	require.NoError(t, json.Unmarshal(q.task.Payload(), &payload))
	require.Equal(t, "ru-RU", payload["language"])
}
