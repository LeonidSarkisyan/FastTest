package questions

import (
	"encoding/json"
	"fmt"
)

type RangeData struct {
	Ranges []RangeAnswer `json:"ranges"`
}

type RangeAnswer struct {
	Text  string `json:"text"`
	Index int    `json:"index"`
}

func NewRangeData() RangeData {
	return RangeData{
		Ranges: []RangeAnswer{
			{Index: 0}, {Index: 1},
		},
	}
}

func UnMarshalRangeData(data []byte) (r RangeData, err error) {
	if err = json.Unmarshal(data, &r); err != nil {
		return RangeData{}, err
	}

	return r, nil
}

func (rd *RangeData) IsValid(index int) error {
	if len(rd.Ranges) < 2 {
		return fmt.Errorf("у вопроса с номером %d должно быть хотя бы 2 пункта", index)
	}

	for _, r := range rd.Ranges {
		if len(r.Text) == 0 {
			return fmt.Errorf("у вопроса с номером %d один из пунктов не имеет текста", index)
		}
	}

	return nil
}

func (rd *RangeData) Scores(rdFromUser RangeData) int {
	if len(rd.Ranges) != len(rdFromUser.Ranges) {
		return 0
	}

	for i := 0; i < len(rd.Ranges); i++ {
		if rd.Ranges[i].Index != rdFromUser.Ranges[i].Index {
			return 0
		}
	}

	return 1
}
