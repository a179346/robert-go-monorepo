package es_bulkwriter

type BulkWriteEvent struct {
	onSuccess func()
	onError   func()
}

func NewBulkWriteEvent(onSuccess func(), onError func()) *BulkWriteEvent {
	return &BulkWriteEvent{
		onSuccess: onSuccess,
		onError:   onError,
	}
}
