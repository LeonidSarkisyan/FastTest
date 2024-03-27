package ai

import (
	"App/internal/models"
	"App/pkg/systems"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
	"strings"
)

func main() {
	systems.SetupLogger()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal().Msgf("ошибка при подключении env - файла: %s", err.Error())
	}

	url := "https://api.proxyapi.ru/openai/v1/chat/completions"

	client := New(url)

	result, err := client.CreateQuestionFromGPT("численные методы", 5)

	log.Info().Any("ответы от chat gpt", result).Send()
}

type OpenAIClient struct {
	url string
}

func New(url string) *OpenAIClient {
	return &OpenAIClient{
		url,
	}
}

func (c *OpenAIClient) CreateQuestionFromGPT(themeName string, countQuestion int) ([]models.QuestionWithAnswersWithOutIsCorrect, error) {

	prompt := `
		напиши тест на русском языке на тему: "` + themeName + `" ответ сформируй в json, где questions - 
		массив вопросов, элементы которого являются объектами (вопросом), у которых есть 
		поля text и поле answers - массив ответов, элементы которого являются объектом (ответом)
		с полями text и is_correct - булево значение. Правильных ответов может быть несколько.

		количество вопросов - ` + strconv.Itoa(countQuestion) + `
		
		очень важно следи за структурой json, так как ты используешься в системе тестирования,
		не используй никаких специальных символов, которые могут помешать взять твой ответ-json.
		`

	messages := []Message{
		{
			Role:    "user",
			Content: prompt,
		},
	}

	messagesFromChatGPT, err := c.AskChatGPT(messages)

	if err != nil {
		log.Err(err).Send()
		return nil, err
	}

	text := messagesFromChatGPT[0].Content
	text = strings.ReplaceAll(strings.ReplaceAll(text, "\\n", ""), "\\t", "")
	text = strings.ReplaceAll(strings.ReplaceAll(text, "\n", ""), "\t", "")

	var questions models.QuestionsFromChatGPT

	if err := json.Unmarshal([]byte(text), &questions); err != nil {
		log.Err(err).Send()
		return nil, err
	}

	return questions.Questions, nil
}

func (c *OpenAIClient) AskChatGPT(messages []Message) ([]Message, error) {
	messageJSON := make([]map[string]string, len(messages))

	for i, message := range messages {
		messageJSON[i] = map[string]string{
			"role":    message.Role,
			"content": message.Content,
		}
	}

	data, err := c.CreateRequest(messageJSON)

	if err != nil {
		log.Err(err).Send()
		return nil, err
	}

	var messageFromGPT = make([]Message, len(data.Choices))

	for i, choice := range data.Choices {
		messageFromGPT[i] = choice.Message
	}

	return messageFromGPT, nil
}

func (c *OpenAIClient) CreateRequest(messages []map[string]string) (Data, error) {
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
		return Data{}, err
	}

	var data Data

	if err := json.Unmarshal([]byte(resp.String()), &data); err != nil {
		log.Err(err).Send()
		return Data{}, err
	}

	return data, nil
}
