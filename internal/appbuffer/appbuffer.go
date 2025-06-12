package appbuffer

import "strings"

type AppBuffer struct {
	buffer []string
}

func (buff *AppBuffer) AppendStr(str string) *AppBuffer {
	buff.buffer = append(buff.buffer, str)
	return buff
}

func (buff *AppBuffer) GetBuffer() []string {
	return buff.buffer
}

func (buff *AppBuffer) Clear() {
	buff.buffer = make([]string, 0)
}

func (buff *AppBuffer) Join(sep string) string {
	return strings.Join(buff.buffer, sep)
}

func (buff *AppBuffer) Len() int {
	return len(buff.buffer)
}
