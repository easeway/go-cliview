package cliview

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

type Output struct {
	Padding   int
	Writer    io.Writer
	Styler    func(class, text string, data interface{}) string
	Formatter func(class string, data interface{}) string
}

func PaddingBuffer(padding int) *bytes.Buffer {
	padBuf := new(bytes.Buffer)
	for i := 0; i < padding; i++ {
		padBuf.WriteString(" ")
	}
	return padBuf
}

func PaddingString(padding int) string {
	return PaddingBuffer(padding).String()
}

func (o *Output) PaddingBuffer() *bytes.Buffer {
	return PaddingBuffer(o.Padding)
}

func (o *Output) PaddingString() string {
	return o.PaddingBuffer().String()
}

func (o *Output) Out() io.Writer {
	if w := o.Writer; w != nil {
		return w
	}
	return os.Stdout
}

func (o *Output) Styling(class, text string, data interface{}) string {
	if o.Styler != nil {
		return o.Styler(class, text, data)
	}
	return text
}

func (o *Output) Format(class string, data interface{}) string {
	if o.Formatter != nil {
		return o.Formatter(class, data)
	} else if data == nil {
		return ""
	}
	return fmt.Sprintf("%v", data)
}
