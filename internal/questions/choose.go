package questions

import (
	"App/internal/models"
	"fmt"
)

func IsValidDefault(q models.QuestionWithAnswers, n int) error {
	if len(q.Answers) < 2 {
		return fmt.Errorf("у вопроса с номером %d меньше двух вариантов ответа, проверьте тест", n)
	}

	var isCorrect bool

	for _, a := range q.Answers {
		if len(a.Text) == 0 {
			return fmt.Errorf("у варианта ответа под вопросом с номером %d нет текста, проверьте тест", n)
		}

		if a.IsCorrect {
			isCorrect = a.IsCorrect
		}
	}

	if !isCorrect {
		return fmt.Errorf("у вопроса с номером %d нет хотя бы одного правильного ответа, проверьте тест", n)
	}

	return nil
}

func HideDataDefault(q *models.QuestionWithAnswers) {
	var countRight int
	for j, a := range q.Answers {
		if a.IsCorrect {
			countRight++
		}
		q.Answers[j].IsCorrect = false
	}
	if countRight >= 2 {
		q.Type = "checkbox"
	} else {
		q.Type = "radio"
	}
}

func ScoreDefault(q models.QuestionForMap, qFromUser models.QuestionWithAnswers) int {
	for _, aFromUser := range qFromUser.Answers {
		a, ok := q.AnswersMap[aFromUser.ID]

		if !ok {
			return 0
		}

		if a.IsCorrect != aFromUser.IsCorrect {
			return 0
		}
	}

	return 1
}
