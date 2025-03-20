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

## Overview

go_thought is a specialized proxy server designed to enhance interactions with Large Language Models. It serves as an intermediary layer that can add various capabilities to LLM requests and responses.

## Features

- Acts as a proxy server for LLM API calls
- Designed for extensibility with new features planned

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

## Getting Started

(Installation and usage instructions will be added as the project develops)

## Contributing

Contributions, suggestions, and feature requests are welcome as the project evolves.

