package logger

import (
	"context"
	"log/slog"
	"os"
)

type SlogAdapter struct {
	log *slog.Logger
}

func NewSlogAdapter(l *slog.Logger) *SlogAdapter {
	return &SlogAdapter{log: l}
}

func (s *SlogAdapter) Info(ctx context.Context, msg string, args ...any) {
	s.log.InfoContext(ctx, msg, args...)
}

func (s *SlogAdapter) Error(ctx context.Context, msg string, args ...any) {
	s.log.ErrorContext(ctx, msg, args...)
}

func (s *SlogAdapter) Fatal(ctx context.Context, msg string, args ...any) {
	s.log.ErrorContext(ctx, msg, args...)
	os.Exit(1)
}

func (s *SlogAdapter) Warn(ctx context.Context, msg string, args ...any) {
	s.log.WarnContext(ctx, msg, args...)
}

func (s *SlogAdapter) Debug(ctx context.Context, msg string, args ...any) {
	s.log.DebugContext(ctx, msg, args...)
}
