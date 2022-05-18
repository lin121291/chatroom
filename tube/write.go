package tube

import (
	"fmt"
	"io"
)

type CommandWriter struct {
	writer io.Writer
}

func NewCommandWriter(writer io.Writer) *CommandWriter {
	return &CommandWriter{
		writer: writer,
	}
}

func (w *CommandWriter) writeString(msg string) error {
	_, err := w.writer.Write([]byte(msg))

	return err
}

//note : 空接口是接口类型的特殊形式，空接口没有任何方法。从实现的角度看，任何值都满足这个接口的需求
//。因此空接口类型可以保存任何值，也可以从空接口中取出原值
func (w *CommandWriter) Write(command interface{}) error {
	var err error

	switch v := command.(type) {
	case SendCommand:
		err = w.writeString(fmt.Sprintf("SEND %v\n", v.Message))
	case MessageCommand:
		err = w.writeString(fmt.Sprintf("MESSAGE %v %v\n", v.Name, v.Message))
	case NameCommand:
		err = w.writeString(fmt.Sprintf("NAME %v\n", v.Name))
	default:
		err = UnknownCommand
	}

	return err
}
