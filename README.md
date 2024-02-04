# go-telegram-gpt DEPRECATED: Engines within GPT-3 are DEPRECATED
# PROJECT NEEDS REFACTOR!

## An ELI5 Telegram Bot using OpenAI API

### Install

1. Clone repository
2. Create `config.yaml`. You will need Telegram token and GPT API token. And a bot on Telegram.

```yaml
tgToken: "xxxx"
gptToken: "xxxx"
preamble: "ELI5: "
```

3. Run `go mod tidy`
4. Run `go run main.go`

### Usage

A very simple Telegram bot that uses OpenAI API to generate ELI5 responses.

#### Commands

- `/start` - Start the bot
- `/topic <topic>` - Genrerate an ELI5 response to the given topic
- `/phrase <phrase>` - Generate an ELI5 response to the given phrase
