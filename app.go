package main

import (
	"context"
	"fmt"
	"log"

	"github.com/PullRequestInc/go-gpt3"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
)

type Config struct {
	TelegramToken string `mapstructure:"tgToken"`
	GptToken      string `mapstructure:"gptToken"`
	Preamble      string `mapstructure:"preamble"`
}

func LoadConfig(path string) (c Config, err error) {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(path)     // path to look for the config file in

	viper.AutomaticEnv() // read in environment variables that match

	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {
		return
	}

	err = viper.Unmarshal(&c) // Unmarshal config
	return
}

func sendChatGPT(apiKey, sendText string) string {
	ctx := context.Background()

	// Create a client
	client := gpt3.NewClient(apiKey)
	var response string

	err := client.CompletionStreamWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt:      []string{sendText},
		MaxTokens:   gpt3.IntPtr(100),
		Temperature: gpt3.Float32Ptr(0),
	}, func(res *gpt3.CompletionResponse) {
		response += res.Choices[0].Text
	})
	if err != nil {
		log.Println(err)
		return "ChatGPT is not available at the moment, please try again later."
	}
	return response
}

func main() {
	// We first check the userPrompt for validity, then we can assign it to the gptPrompt
	var userPrompt string
	var gptPrompt string

	config, err := LoadConfig(".")
	apiKey := config.GptToken

	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	bot, err := tgbotapi.NewBotAPI(config.TelegramToken)
	if err != nil {
		panic(err)
	}

	// be able to debug our logs in runtime
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
}
