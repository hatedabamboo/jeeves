# jeeves

Your own personal (miniature) butler. In CLI.

## Overview

This is a command-line utility to query ChatGPT. As a command-line tool, the
main scope of this utility is quick questions that require fast answers.

## Installation

```shell
go build -o jeeves main.go
sudo mv jeeves /usr/local/bin/
```
## Usage

```shell
~ $ jeeves What is the answer to the Ultimate Question of Life, the Universe, and Everything?

The answer to the Ultimate Question of Life, the Universe, and Everything is famously given as the number **42** in Douglas Adams' science fiction series, "The Hitchhiker's Guide to the Galaxy." However, the actual Ultimate Question itself remains unknown, leading to much humor and speculation within the narrative.

```

## Configuration

jeeves supports configuration via the environment variables:

- `OPENAI_API_KEY` your OpenAI API key; obtain one from [OpenAI Platform](https://platform.openai.com/docs/overview) (mandatory).
- `JEEVES_OPENAI_MODEL` OpenAI model of choice; default is `gpt-4o-mini`, which allows communicate with ChatGPT for free (optional).
- `JEEVES_LOG_LEVEL` log level; default is `info`; supported values are: `info`, `debug` (optional).
- `JEEVES_CUSTOM_PROMPT` additional custom message that will be added before each request; default is `""` (optional)
