package roberthttp_extended

type JsonResponse[T interface{}] struct {
	Data T `json:"data"`
}
