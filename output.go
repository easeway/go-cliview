package cliview

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

type StylerFunc func(class, text string, data interface{}) string
type FormatterFunc func(class string, data interface{}, formatter FormatterFunc) string

type Output struct {
	Padding   int
	Writer    io.Writer
	Styler    StylerFunc
	Formatter FormatterFunc
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

func (o *Output) Styling(class, text string, data interface{}, styler StylerFunc) string {
	if styler == nil {
		styler = o.Styler
	}
	if styler != nil {
		return styler(class, text, data)
	}
	return text
}

func defaultFormatter(class string, data interface{}, formatter FormatterFunc) string {
	if data == nil {
		return ""
	}
	switch data.(type) {
	case float32, float64:
		return fmt.Sprintf("%g", data.(float64))
	}
	return fmt.Sprintf("%v", data)
}

func (o *Output) Format(class string, data interface{}) string {
	if o.Formatter != nil {
		return o.Formatter(class, data, defaultFormatter)
	} else {
		return defaultFormatter(class, data, nil)
	}
}
