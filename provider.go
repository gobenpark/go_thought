//go:generate mockgen -source=./provider.go -destination=./mock_provider.go -package=main
package main

import (
	"context"

	"github.com/gobenpark/gothought/tool"
)

type Provider interface {
	Generate(ctx context.Context, tools map[string]tool.Tool, messages []Message) (*Message, string, error)
}

type StreamingCapable interface {
	GenerateStreaming(ctx context.Context, messages []Message, callback func(Message) error) error
}
