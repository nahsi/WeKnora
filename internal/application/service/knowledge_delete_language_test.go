package service

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Tencent/WeKnora/internal/types"
	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/require"
)

func TestDeleteTaskPreservesWikiRetractLanguage(t *testing.T) {
	for _, tt := range []struct {
		name     string
		language string
	}{
		{name: "request language", language: "ru-RU"},
		{name: "legacy payload"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f := newDocumentWriteFixture(t)
			f.kbs.values["kb"].IndexingStrategy.WikiEnabled = true
			f.svc.wikiRepo = &moveWikiPageRepo{}
			f.svc.wikiService = &moveWikiMetaService{}
			f.svc.taskPendingRepo = &moveWikiPendingRepo{}
			q := &autoTagRecordingQueue{}
			f.svc.task = q
			payload := map[string]interface{}{
				"tenant_id": 7, "knowledge_base_id": "kb", "knowledge_ids": []string{"doc"},
			}
			if tt.language != "" {
				payload["language"] = tt.language
			}
			raw, err := json.Marshal(payload)
			require.NoError(t, err)
			ctx := context.Background()
			require.NoError(t, f.svc.ProcessKnowledgeListDelete(ctx, asynq.NewTask(types.TypeKnowledgeListDelete, raw)))
			require.NotNil(t, q.task)
			var retract WikiRetractPayload
			require.NoError(t, json.Unmarshal(q.task.Payload(), &retract))
			want := tt.language
			if want == "" {
				want = types.DefaultLanguage()
			}
			require.Equal(t, want, retract.Language)
		})
	}
}

func TestDeleteTagPreservesLanguage(t *testing.T) {
	f := newDocumentWriteFixture(t)
	f.kbs.values["kb"].Type = types.KnowledgeBaseTypeDocument
	require.NoError(t, f.repo.SetKnowledgeTags(f.ctx, "doc", []string{"tag"}))
	q := &autoTagRecordingQueue{}
	svc := &knowledgeTagService{
		repo: f.svc.tagRepo, kbService: f.kbs, tenantRepo: f.tenants,
		knowledgeRepo: f.repo, task: q,
	}
	ctx := context.WithValue(f.ctx, types.LanguageContextKey, "ru-RU")
	require.NoError(t, svc.DeleteTag(ctx, "tag", true, true, nil))
	require.NotNil(t, q.task)
	var payload map[string]interface{}
	require.NoError(t, json.Unmarshal(q.task.Payload(), &payload))
	require.Equal(t, "ru-RU", payload["language"])
}
