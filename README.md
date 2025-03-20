# go_thought

A proxy server for Large Language Models (LLMs) with plans for extended functionality.

## Usage

just set env variable `OPENAI_BASE_URL`

`OPENAI_BASE_URL=http://localhost:8080`

```python
from langchain.chat_models.base import init_chat_model

llm = init_chat_model(model="gpt-4o-mini", model_provider="openai")

llm.invoke({})
```



## Why go_thought?

I created go_thought to solve common challenges when working with LLM applications

`Zero-code integration`: Add powerful features to your LLM pipeline without modifying your application code
`Visibility`: Gain insights into your LLM usage patterns, performance, and costs
`Optimization`: Improve performance and reduce costs with intelligent caching and request management
`Flexibility`: Work with multiple LLM providers through a single, consistent interface

go_thought acts as a transparent layer between your application and LLM providers, adding valuable functionality without disrupting your existing workflow.

## Features

- Acts as a proxy server for LLM API calls
- Seamless integration with existing LLM applications
- Designed for extensibility with robust plugin architecture

## Roadmap

Future enhancements planned for go_thought include:
- Request/response logging
- Rate limiting
- Caching mechanisms
- Custom preprocessing of prompts
- Response transformation capabilities
- Support for multiple LLM providers
- Performance analytics
- Prompt analytics
- Cost tracking and optimization

## Getting Started

(Installation and usage instructions will be added as the project develops)

## Contributing

Contributions, suggestions, and feature requests are welcome as the project evolves.

