package app

import (
	"context"
	"log"
	"os"
	"rag/chat"
	"rag/config"
	"rag/ingest"
	"rag/llm"
	"rag/rag"
	"rag/vector"
	"rag/vector/pgvector"
	"sync"
)

func Run(parent context.Context, cfg config.Config) error {

	logger := log.New(os.Stderr, "[rag] ", log.LstdFlags)

	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	client := llm.New(cfg)

	embedder := llm.NewEmbedder(cfg)

	store, err := openStore(ctx, cfg)
	if err != nil {
		logger.Printf("vector store disable: %v", err)
	}

	var wg sync.WaitGroup
	if store != nil {
		wg.Go(func() {
			opts := ingest.Options{
				SourceDir:    cfg.IngestDir,
				ProcessedDir: cfg.ProcessDir,
			}

			if err := ingest.Watch(ctx, opts, embedder, store, logger); err != nil && ctx.Err() == nil {
				logger.Printf("watcher stopped: %v", err)
			}
		})

		logger.Printf("watching %s for new documents", cfg.IngestDir)
	}

	if store != nil {
		defer store.Close()
		logger.Printf("vector store ready")
	}

	var retriever *rag.Retriever
	if store != nil {
		retriever = rag.New(embedder, store, rag.Options{
			TopK:     5,
			Rewriter: rag.NewRewriter(client),
		})
	}

	replErr := chat.RunREPL(ctx, client, retriever, chat.Options{
		SystemPromptFile: cfg.SystemPromptFile,
	})

	cancel()
	wg.Wait()

	return replErr
}

func openStore(ctx context.Context, cfg config.Config) (vector.Store, error) {

	if cfg.DatabaseURL == "" {
		return nil, nil
	}

	s, err := pgvector.New(ctx, pgvector.Options{
		DSN:          cfg.DatabaseURL,
		EmbeddingDim: cfg.EmbeddingDim,
	})

	if err != nil {
		return nil, err
	}

	return s, nil

}
