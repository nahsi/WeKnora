package service

import (
	"testing"

	"github.com/Tencent/WeKnora/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestWikiSuggestionFromPage_Locales(t *testing.T) {
	tests := []struct {
		locale  string
		concept string
		about   string
	}{
		{"zh-CN", "什么是RAG？", "介绍一下RAG"},
		{"zh-TW", "什麼是RAG？", "介紹一下RAG"},
		{"en-US", "What is RAG?", "Tell me about RAG"},
		{"ru-RU", "Что такое RAG?", "Расскажи о RAG"},
		{"ko-KR", "RAG란 무엇인가요?", "RAG에 대해 알려주세요"},
		{"ja-JP", "RAGとは？", "RAGについて教えて"},
		{"en-GB", "What is RAG?", "Tell me about RAG"},
		{"en-AU", "What is RAG?", "Tell me about RAG"},
		{"en", "What is RAG?", "Tell me about RAG"},
		{"ru", "Что такое RAG?", "Расскажи о RAG"},
		{"ru-BY", "Что такое RAG?", "Расскажи о RAG"},
		{"zh", "什么是RAG？", "介绍一下RAG"},
		{"ko", "RAG란 무엇인가요?", "RAG에 대해 알려주세요"},
		{"ja", "RAGとは？", "RAGについて教えて"},
		{"fr-FR", "RAG", "RAG"},
		{"english", "RAG", "RAG"},
		{"", "RAG", "RAG"},
	}
	for _, tt := range tests {
		t.Run(tt.locale, func(t *testing.T) {
			for _, pageType := range []string{
				types.WikiPageTypeConcept,
				types.WikiPageTypeEntity,
				types.WikiPageTypeSummary,
				types.WikiPageTypeSynthesis,
				types.WikiPageTypeComparison,
				types.WikiPageTypeIndex,
				"unknown",
			} {
				t.Run(pageType, func(t *testing.T) {
					want := "RAG"
					switch pageType {
					case types.WikiPageTypeConcept:
						want = tt.concept
					case types.WikiPageTypeEntity, types.WikiPageTypeSummary:
						want = tt.about
					}
					page := &types.WikiPage{Title: " \tRAG\n ", PageType: pageType}
					assert.Equal(t, want, wikiSuggestionFromPage(page, tt.locale))
				})
			}
		})
	}
}

func TestWikiSuggestionFromPage_EmptyTitle(t *testing.T) {
	tests := []struct {
		name string
		page *types.WikiPage
	}{
		{"nil page", nil},
		{"empty title", &types.WikiPage{PageType: types.WikiPageTypeConcept}},
		{"blank title", &types.WikiPage{Title: " \t\n ", PageType: types.WikiPageTypeConcept}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Empty(t, wikiSuggestionFromPage(tt.page, "ru-RU"))
		})
	}
}
