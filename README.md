# go_thought

A lightweight, intuitive library for building LLM-powered applications in Go.

## What is go_thought?

go_thought provides a simple, fluent API for interacting with Large Language Models. Unlike more complex frameworks, go_thought focuses on minimizing boilerplate code while maintaining flexibility.

```go
cli := NewClient(
    WithApiKey(os.Getenv("OPENAI_API_KEY")),
    WithModel("gpt-4o-mini"),
)

response, err := cli.State(OPENAI).
    SystemPrompt("You are a helpful assistant.").
    HumanPrompt("Tell me about Go programming.").
    Q(context.Background())
```

## Why go_thought?

While solutions like langchain-go offer comprehensive features, they often require significant configuration and understanding of complex abstractions. go_thought aims to solve common challenges in LLM application development with:

- **Minimal Setup**: Get started with just a few lines of code
- **Fluent API**: Intuitive chain-style syntax for building prompts
- **Extensible Design**: Support for multiple LLM providers through a unified interface
- **Type Safety**: Leverage Go's type system for reliable code

## Features

- Simple, chainable API for prompt construction
- Support for different message roles (system, human, AI)
- State pattern design for easy provider switching
- Functional options for flexible configuration

## Installation

```bash
go get github.com/gobenpark/go_thought
```

## Quickstart

```go
package main

import (
    "context"
    "fmt"
    "os"
    
    "github.com/gobenpark/go_thought"
)

func main() {
    // Initialize with your API key
    cli := NewClient(
        WithApiKey(os.Getenv("OPENAI_API_KEY")),
        WithModel("gpt-4o"),
    )
    
    // Build your prompt chain and execute
    response, err := cli.State(OPENAI).
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

## Supported LLM Providers

- OpenAI (ChatGPT, GPT-4)
- More coming soon!

## Roadmap

Future plans for go_thought include:

- Response parsing and structured output
- Streaming support for real-time responses
- Context management for multi-turn conversations
- Additional LLM providers (Anthropic, Gemini, etc.)
- Helper functions for common LLM tasks
- Middleware support for request/response processing
- Caching mechanisms
- Prompt templates

## Contributing

Contributions, suggestions, and feature requests are welcome! Feel free to open issues or submit pull requests as the project evolves.
