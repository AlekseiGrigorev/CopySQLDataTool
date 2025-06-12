package app

type RowsProcessorInterface interface {
	Write(buffer []string, data []any) error
	GetProcessedMsg() string
}
