package drivers

type BaseDriver interface {
	Send(record map[string]interface{}) (string, error)
	Close() error
}
