package service

import (
	"App/internal/ai"
	"App/internal/models"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type AiService struct {
	QuestionRepository
	TestRepository
	*ai.OpenAIClient
}

func NewAiService(questionRepository QuestionRepository, testRepository TestRepository) *AiService {
	return &AiService{
		questionRepository, testRepository, ai.New(viper.GetString("open_ai_url")),
	}
}

func (s *AiService) CreateQuestionsFromGPT(userID, testID int, promptParams ai.PromptParams) ([]models.QuestionWithAnswersWithOutIsCorrect, error) {
	_, err := s.TestRepository.Get(testID, userID)

	if err != nil {
		log.Err(err).Send()
		return nil, err
	}

	questions, err := s.OpenAIClient.CreateQuestionFromGPT(promptParams.TitleTheme, promptParams.CountQuestions)

	if err != nil {
		log.Err(err).Send()
		return nil, err
	}

	questionsWithIDs, err := s.QuestionRepository.CreateManyQuestions(testID, questions)

	if err != nil {
		return nil, err
	}

	return questionsWithIDs, nil
}
