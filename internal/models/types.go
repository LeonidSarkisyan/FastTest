package models

const (
	Choose = "choose"
	Group_ = "group"
	Range_ = "range"

	CountGroups        = 3
	CountAnswerInGroup = 1
	CountRanges        = 3
)

func GetTextFromType(type_ string) string {
	switch type_ {
	case Group_:
		return "Установите соответствие между группой и вариантами ответов"
	case Range_:
		return "Расставьте пункты по порядку"
	default:
		return ""
	}
}
