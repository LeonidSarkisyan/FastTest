package responses

func NewListResponse[T any](list []T) []T {
	if len(list) == 0 {
		list = make([]T, 0)
	}
	return list
}
