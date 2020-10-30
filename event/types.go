package event

// Event Событие
type Event struct {
	// SourceID Идентификатор процесса создавшего событие
	SourceID string

	// TargetID Идентификатор процесса назначения события, если пусто - все процессы
	TargetID string

	// Act Действие которое необходимо выполнить, либо произошедная смена состояния
	Act Operation

	// Err Событие ошибки
	Err error
}

// Operation type
type Operation string
