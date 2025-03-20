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

## Getting Started

(Installation and usage instructions will be added as the project develops)

## Contributing

Contributions, suggestions, and feature requests are welcome as the project evolves.

## License

MIT License

Copyright (c) 2025

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.