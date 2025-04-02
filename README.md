<p align="center">
  <img src="./go_thought.png" alt="gothought Logo" width="300">
</p>

# gothought

A lightweight, intuitive library for building LLM-powered applications in Go.

## What is gothought?

gothought provides a simple, fluent API for interacting with Large Language Models. Unlike more complex frameworks, gothought focuses on minimizing boilerplate code while maintaining flexibility.

```go
cli := gothought.NewClient(
    gothought.WithApiKey(os.Getenv("OPENAI_API_KEY")),
    gothought.WithModel("gpt-4o-mini"),
)

response, err := cli.State(gothought.OPENAI).
    SystemPrompt("You are a helpful assistant.").
    HumanPrompt("Tell me about Go programming.").
    Q(context.Background())
```

## Why gothought?

While solutions like langchain-go offer comprehensive features, they often require significant configuration and understanding of complex abstractions. gothought aims to solve common challenges in LLM application development with:

- **Minimal Setup**: Get started with just a few lines of code
- **Fluent API**: Intuitive chain-style syntax for building prompts
- **Extensible Design**: Support for multiple LLM providers through a unified interface
- **Type Safety**: Leverage Go's type system for reliable code

## Features

- Simple, chainable API for prompt construction
- Support for different message roles (system, human, AI)
- State pattern design for easy provider switching
- Functional options for flexible configuration
- Streaming support for real-time responses
- Multiple LLM providers (OpenAI, Anthropic Claude)

## Installation

```bash
go get github.com/gobenpark/gothought
```

## Quickstart

```go
package main

import (
    "context"
    "fmt"
    "os"
    
    "github.com/gobenpark/gothought"
)

func main() {
    // Initialize with your API key
    cli := gothought.NewClient(
        gothought.WithApiKey(os.Getenv("OPENAI_API_KEY")),
        gothought.WithModel("gpt-4o"),
    )
    
    // Build your prompt chain and execute
    response, err := cli.State(gothought.OPENAI).
        SystemPrompt("You are a helpful coding assistant.").
        HumanPrompt("How do I read a file in Go?").
        Q(context.Background())
        
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    
    fmt.Println("Response:", response)
}
```

### Using Claude

```go
cli := gothought.NewClient(
    gothought.WithApiKey(os.Getenv("ANTHROPIC_API_KEY")),
    gothought.WithModel("claude-3-opus-20240229"),
)

response, err := cli.State(gothought.CLAUDE).
    SystemPrompt("You are a helpful assistant specialized in Go programming.").
    HumanPrompt("What are goroutines and how do they work?").
    Q(context.Background())
```

### Streaming Responses

```go
err := cli.State(gothought.OPENAI).
    SystemPrompt("You are a helpful assistant.").
    HumanPrompt("Tell me a story about space exploration.").
    QStream(context.Background(), func(msg gothought.ResponseMessage) error {
        fmt.Print(msg.Message) // Print message chunks as they arrive
        return nil
    })
```

## Supported LLM Providers

- OpenAI (ChatGPT, GPT-4)
- Anthropic Claude
- More coming soon!

## Roadmap

Future plans for gothought include:

- Context management for multi-turn conversations
- Additional LLM providers (Gemini, Cohere, etc.)
- Helper functions for common LLM tasks
- Middleware support for request/response processing
- Caching mechanisms
- Prompt templates
- Token counting and management
- Rate limiting and retry strategies

## Contributing

Contributions, suggestions, and feature requests are welcome! Feel free to open issues or submit pull requests as the project evolves.