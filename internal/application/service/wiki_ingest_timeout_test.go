package service

import (
	"context"
	"testing"
	"time"

	"github.com/Tencent/WeKnora/internal/config"
	"github.com/Tencent/WeKnora/internal/types"
	"github.com/alicebob/miniredis/v2"
	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWikiClaimTimeouts(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name     string
		cfg      *config.Config
		stale    time.Duration
		identity time.Duration
	}{
		{"defaults", nil, 90 * time.Minute, 2 * time.Hour},
		{"configured", &config.Config{
			KnowledgeBase: &config.KnowledgeBaseConfig{WikiIngestTimeout: 4 * time.Hour},
		}, 4*time.Hour + 30*time.Minute, 5 * time.Hour},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.stale, wikiClaimStaleAfter(tc.cfg))
			assert.Equal(t, tc.identity, wikiIdentityClaimTTL(tc.cfg))
		})
	}
}

func TestEnqueueWikiIngestTriggerConfiguredTimeout(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{
		KnowledgeBase: &config.KnowledgeBaseConfig{WikiIngestTimeout: 4 * time.Hour},
	}
	queue := &wikiGuardTaskQueue{}
	require.NoError(t, enqueueWikiIngestTrigger(context.Background(), cfg, queue, 1, "kb-1"))
	require.Len(t, queue.tasks, 1)
	task := queue.tasks[0]
	assert.Equal(t, types.TypeWikiIngest, task.Type())
	rdb := miniredis.RunT(t)
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: rdb.Addr()})
	t.Cleanup(func() { _ = client.Close() })
	info, err := client.Enqueue(task)
	require.NoError(t, err)
	assert.Equal(t, 4*time.Hour, info.Timeout)
}
