package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal().Msgf("ошибка при подключении env - файла: %s", err.Error())
	}

	url := "https://api.openai.com/v1/chat/completions"

	client := New(url)

	messages := []map[string]string{
		{
			"role":    "user",
			"content": "Привет!",
		},
	}

	result, err := client.CreateRequest(messages)

	if err != nil {
		log.Err(err).Send()
		return
	}

	log.Info().Str("result", result).Send()
}

type OpenAIClient struct {
	url string
}

func (c *OpenAIClient) CreateRequest(messages []map[string]string) (string, error) {
	resp, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY")).
		SetBody(map[string]interface{}{
			"model":       "gpt-3.5-turbo",
			"messages":    messages,
			"temperature": 0.7,
		}).
		Post(c.url)

	if err != nil {
		log.Err(err).Send()
		return "", err
	}

	log.Info().Int("status", resp.StatusCode()).Send()
	log.Info().Any("result", resp.Result()).Send()

	return resp.String(), nil
}

func New(url string) *OpenAIClient {
	return &OpenAIClient{
		url,
	}
}
