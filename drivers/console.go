package drivers

import (
	"encoding/json"
	"fmt"
)

type ConsoleDriver struct{}

func NewConsoleDriver() *ConsoleDriver {
	return &ConsoleDriver{}
}

func (d *ConsoleDriver) Send(record map[string]interface{}) (string, error) {
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return "", err
	}
	fmt.Println(string(data))
	return "console-log", nil
}

func (d *ConsoleDriver) Close() error {
	return nil
}
