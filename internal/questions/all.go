package questions

import (
	"App/internal/models"
	"encoding/json"
	"errors"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var (
	NotFoundType = errors.New("неизвестный тип вопроса")
)

const (
	Group_ = "group"
	Range  = "range"
)

func AsBytes[t any](d t) []byte {
	data, err := json.Marshal(d)

	if err != nil {
		log.Err(err).Send()
	}

	return data
}

func CheckJSONStructure(type_ string, data []byte) (err error) {
	switch type_ {
	case Group_:
		_, err = UnMarshalGroupData(data)
	case Range:
		_, err = UnMarshalRangeData(data)
	default:
		return NotFoundType
	}

	return err
}

func UnMarshalData(type_ string, data []byte) (d any, err error) {
	switch type_ {
	case Group_:
		d, err = UnMarshalGroupData(data)
	case Range:
		d, err = UnMarshalRangeData(data)
	default:
		return nil, nil
	}

	return
}

func CreateDataForType(type_ string) []byte {
	var data any

	switch type_ {
	case Group_:
		data = NewGroupData()
	case Range:
		data = NewRangeData()
	}

	return AsBytes(data)
}

func Check(q models.QuestionWithAnswers, index int) (err error) {

	switch qData := q.Data.(type) {
	case GroupData:
		err = qData.IsValid(index)
	case RangeData:
		err = qData.IsValid(index)
	default:
		err = IsValidDefault(q, index)
	}

	return err
}

func HideData(questions []models.QuestionWithAnswers) (err error) {
	for i, q := range questions {
		switch q.Type {
		case Group_:
			var data GroupData

			err = mapstructure.Decode(q.Data, &data)

			data.HideData()

			questions[i].Data = data
		default:
			HideDataDefault(&questions[i])
		}

		if err != nil {
			log.Err(err).Send()
			return err
		}
	}

	return nil
}

func Score(questions, questionsFromUser []models.QuestionWithAnswers) (int, error) {
	questionMap := questionSliceToMap(questions)

	var score int

	for _, qFromUser := range questionsFromUser {
		q, ok := questionMap[qFromUser.ID]

		if !ok {
			log.Error().Msg("вопрос не был найден при оценивании результатов")
			continue
		}

		switch qFromUser.Type {
		case Group_:
			log.Info().Any("DATA", qFromUser.Data).Send()

			var data GroupData

			if err := mapstructure.Decode(qFromUser.Data, &data); err != nil {
				return 0, errors.New("некорректная структура")
			}

			var qData GroupData

			if err := mapstructure.Decode(q.Data, &qData); err != nil {
				return 0, errors.New("некорректная структура")
			}

			score += data.Scores(qData)

			log.Info().Int("score +++", data.Scores(qData)).Send()
		default:
			score += ScoreDefault(q, qFromUser)
		}
	}

	return score, nil
}

func questionSliceToMap(questions []models.QuestionWithAnswers) map[int]models.QuestionForMap {
	questionMap := make(map[int]models.QuestionForMap, len(questions))

	for _, q := range questions {
		answerMap := make(map[int]models.Answer, len(q.Answers))
		for _, a := range q.Answers {
			answerMap[a.ID] = a
		}
		questionMap[q.ID] = models.QuestionForMap{
			Question:   models.Question{ID: q.ID, Text: q.Text, Data: q.Data, Type: q.Type},
			AnswersMap: answerMap,
		}
	}

	return questionMap
}
