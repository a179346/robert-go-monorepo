package es_bulkrequester

type BulkItemEvent struct {
	onSuccess func()
	onError   func()
}

func NewBulkItemEvent(onSuccess func(), onError func()) *BulkItemEvent {
	return &BulkItemEvent{
		onSuccess: onSuccess,
		onError:   onError,
	}
}
