package config

import (
	"testing"
	"time"
)

func TestWikiIngestTimeout(t *testing.T) {
	cases := []struct {
		name     string
		cfg      *Config
		expected time.Duration
	}{
		{"nil config", nil, 60 * time.Minute},
		{"nil knowledge base", &Config{}, 60 * time.Minute},
		{"zero", &Config{KnowledgeBase: &KnowledgeBaseConfig{}}, 60 * time.Minute},
		{"negative", &Config{KnowledgeBase: &KnowledgeBaseConfig{WikiIngestTimeout: -time.Minute}}, 60 * time.Minute},
		{"set", &Config{KnowledgeBase: &KnowledgeBaseConfig{WikiIngestTimeout: 15 * time.Minute}}, 15 * time.Minute},
		{"no upper bound", &Config{
			KnowledgeBase: &KnowledgeBaseConfig{WikiIngestTimeout: 72 * time.Hour},
		}, 72 * time.Hour},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := WikiIngestTimeout(tc.cfg); got != tc.expected {
				t.Fatalf("WikiIngestTimeout() = %v, want %v", got, tc.expected)
			}
		})
	}
}

func TestApplyKnowledgeBaseEnvOverrides_WikiIngestTimeout(t *testing.T) {
	cases := []struct {
		name          string
		knowledgeBase *KnowledgeBaseConfig
		env           string
		expected      time.Duration
	}{
		{"nil knowledge base", nil, "", 60 * time.Minute},
		{"zero defaults", &KnowledgeBaseConfig{}, "", 60 * time.Minute},
		{"negative defaults", &KnowledgeBaseConfig{WikiIngestTimeout: -time.Minute}, "", 60 * time.Minute},
		{"preserves config", &KnowledgeBaseConfig{WikiIngestTimeout: 15 * time.Minute}, "", 15 * time.Minute},
		{"environment overrides config", &KnowledgeBaseConfig{
			WikiIngestTimeout: 15 * time.Minute,
		}, "25m", 25 * time.Minute},
		{"blank ignored", &KnowledgeBaseConfig{WikiIngestTimeout: 15 * time.Minute}, "  ", 15 * time.Minute},
		{"invalid ignored", &KnowledgeBaseConfig{WikiIngestTimeout: 15 * time.Minute}, "invalid", 15 * time.Minute},
		{"zero environment ignored", &KnowledgeBaseConfig{
			WikiIngestTimeout: 15 * time.Minute,
		}, "0s", 15 * time.Minute},
		{"negative environment ignored", &KnowledgeBaseConfig{
			WikiIngestTimeout: 15 * time.Minute,
		}, "-1m", 15 * time.Minute},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("WEKNORA_WIKI_INGEST_TIMEOUT", tc.env)
			cfg := &Config{KnowledgeBase: tc.knowledgeBase}

			applyKnowledgeBaseEnvOverrides(cfg)

			if got := cfg.KnowledgeBase.WikiIngestTimeout; got != tc.expected {
				t.Fatalf("wiki_ingest_timeout = %v, want %v", got, tc.expected)
			}
		})
	}
}
