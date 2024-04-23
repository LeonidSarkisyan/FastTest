package questions

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"math/rand/v2"
	"sort"
)

type GroupData struct {
	Answers []Answer `json:"answers"`
	Groups  []Group  `json:"groups"`
}

type Answer struct {
	Text       string `json:"text"`
	GroupIndex int    `json:"group_index" mapstructure:"group_index"`
}

type Group struct {
	Title string `json:"title"`
}

func NewGroupData() GroupData {
	return GroupData{
		Groups: make([]Group, 2),
		Answers: []Answer{
			{GroupIndex: 0}, {GroupIndex: 1},
		},
	}
}

func UnMarshalGroupData(data []byte) (d GroupData, err error) {
	if err = json.Unmarshal(data, &d); err != nil {
		return GroupData{}, err
	}

	sort.Slice(d.Answers, func(i, j int) bool {
		return d.Answers[i].GroupIndex < d.Answers[j].GroupIndex
	})

	return d, nil
}

func (gd *GroupData) IsValid(index int) error {
	if len(gd.Groups) < 2 {
		return fmt.Errorf("у вопроса с номером %d должно быть хотя бы 2 группы", index)
	}

	for _, g := range gd.Groups {
		if len(g.Title) == 0 {
			return fmt.Errorf("у вопроса с номером %d группа не имеет текста, проверьте тест", index)
		}
	}

	if len(gd.Answers) < 1 {
		return fmt.Errorf("у вопроса с номером %d должно быть хотя бы один вариант ответа", index)
	}

	for _, a := range gd.Answers {
		if len(a.Text) == 0 {
			return fmt.Errorf("у вопроса с номером %d вариант ответа не имеет текста, проверьте тест", index)
		}
	}

	return nil
}

func (gd *GroupData) Scores(gdFromUser GroupData) int {
	if len(gd.Answers) != len(gdFromUser.Answers) {
		return 0
	}

	sort.Slice(gd.Answers, func(i, j int) bool {
		return gd.Answers[i].Text < gd.Answers[j].Text
	})

	sort.Slice(gdFromUser.Answers, func(i, j int) bool {
		return gdFromUser.Answers[i].Text < gdFromUser.Answers[j].Text
	})

	log.Info().Any("answers", gd.Answers).Any("answers from users", gdFromUser.Answers).Send()

	for i := 0; i < len(gd.Answers); i++ {
		if gd.Answers[i].GroupIndex != gdFromUser.Answers[i].GroupIndex+1 {
			return 0
		}
	}

	return 1
}

func (gd *GroupData) HideData() GroupData {
	gdh := GroupData{
		Answers: gd.Answers,
		Groups:  gd.Groups,
	}

	for i := 0; i < len(gdh.Answers); i++ {
		gdh.Answers[i].GroupIndex = 0
	}

	rand.Shuffle(len(gdh.Answers), func(i, j int) {
		gdh.Answers[i], gdh.Answers[j] = gdh.Answers[j], gdh.Answers[i]
	})

	return gdh
}
