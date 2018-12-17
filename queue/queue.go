package queue

// Queue provides interface for message queue operations
type Queue interface {
	Close()
	MakeCDMsgChan() (<-chan []byte, error)
}
