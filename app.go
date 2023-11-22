package main

import (
	"context"
	"fmt"
	"log"
	"strings"

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

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u) // get any new updates from the bot

	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if !strings.HasPrefix(update.Message.Text, "/topic") && !strings.HasPrefix(update.Message.Text, "/phrase") {
			continue
		}

		if strings.HasPrefix(update.Message.Text, "/topic") {
			userPrompt = strings.TrimPrefix(update.Message.Text, "/topic")
			gptPrompt = config.Preamble + "TOPIC: "
		} else if strings.HasPrefix(update.Message.Text, "/phrase") {
			userPrompt = strings.TrimPrefix(update.Message.Text, "/phrase")
			gptPrompt = config.Preamble + "PHRASE: "
		}

		if userPrompt != "" {
			gptPrompt += userPrompt
			response := sendChatGPT(apiKey, gptPrompt)
			update.Message.Text = response
		} else {
			update.Message.Text = "Please enter a valid prompt of /topic or /phrase."
		}

		// Now we send this to telegram
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		// Send the message
		_, err := bot.Send(msg)
		if err != nil {
			log.Println("ERROR: ", err)
		}
	}
}
