package cliview

import (
	"bytes"
	"fmt"
	"io"
)

const (
	AlignLeft   = 0
	AlignMiddle = 1
	AlignRight  = 2
)

type Column struct {
	Title     string // title to be displayed
	Field     string // field name for retrieving data
	Width     int    // column width, >0 fixed, =0 auto, <0 percentage
	MaxWidth  int    // maximum column width
	Align     int
	Fetcher   func(column Column, row map[string]interface{}) interface{}
	Formatter FormatterFunc
}

type Table struct {
	Output

	Columns  []Column // Column definitions
	MaxWidth int      // maximum table width

	columns []Column // actuall columns
}

func (tv *Table) Print(data []map[string]interface{}) {
	// calculate column width
	tv.columns = make([]Column, len(tv.Columns))
	fixedWidth := 0
	for i := 0; i < len(tv.Columns); i++ {
		if tv.Columns[i].Width > 0 {
			tv.columns[i].Width = tv.Columns[i].Width
			fixedWidth += tv.Columns[i].Width
		} else if tv.Columns[i].Width == 0 {
			width := len(tv.Columns[i].Title)
			for _, v := range data {
				valLen := len(tv.formatCell("table:row:", tv.Columns[i], v))
				if valLen > width {
					width = valLen
				}
				if tv.Columns[i].MaxWidth > 0 && width > tv.Columns[i].MaxWidth {
					width = tv.Columns[i].MaxWidth
				}
			}
			tv.columns[i].Width = width
			fixedWidth += width
		} else {
			tv.columns[i].Width = tv.Columns[i].Width
		}
	}
	restWidth := tv.MaxWidth - fixedWidth - len(tv.Columns) - 1

	// calculate final width
	width := 1
	for i := 0; i < len(tv.columns); i++ {
		if tv.columns[i].Width < 0 {
			tv.columns[i].Width = restWidth * -tv.columns[i].Width / 100
		}
		if tv.columns[i].Width > 0 {
			tv.columns[i].Title = tv.Columns[i].Title
			tv.columns[i].Field = tv.Columns[i].Field
			tv.columns[i].Align = tv.Columns[i].Align
			width += tv.columns[i].Width
		}
		width++
	}

	if width <= 2 {
		return
	}

	// print head
	row := tv.startPrintRow()
	for _, c := range tv.columns {
		row.column("head", c.Title, &c, c.Title)
	}
	row.top()
	row.end()

	// print rows
	for _, d := range data {
		row = tv.startPrintRow()
		for i, c := range tv.columns {
			val := d[c.Field]
			if c.Width > 0 {
				valStr := tv.formatCell("table:row:", tv.Columns[i], d)
				row.column("row", valStr, &c, val)
			} else {
				row.column("row", "", &c, val)
			}
		}
		row.end()
	}
}

func (tv *Table) formatCell(classPrefix string, col Column, row map[string]interface{}) string {
	var val interface{}
	if col.Fetcher != nil {
		val = col.Fetcher(col, row)
	} else {
		val = row[col.Field]
	}
	class := classPrefix + col.Field
	if col.Formatter != nil {
		return col.Formatter(class, val, func(class string, data interface{}, formatter FormatterFunc) string {
			return tv.Format(class, data)
		})
	}
	return tv.Format(class, val)
}

func ellipsis(text string, width int) string {
	if len(text) <= width {
		return text
	}
	if width <= 3 {
		text = text[0:1]
		for i := 1; i < width; i++ {
			text += "."
		}
	} else {
		text = text[0:width-3] + "..."
	}
	return text
}

func wrapLen(text string, width int, align int) string {
	if i := len(text); i < width {
		buf := new(bytes.Buffer)
		for ; i < width; i++ {
			buf.WriteString(" ")
		}
		pads := buf.String()
		switch align {
		case AlignRight:
			return pads + text
		case AlignMiddle:
			left := (width - len(text)) / 2
			return pads[0:left] + text + pads[left:len(pads)]
		default:
			return text + pads
		}
	}
	return text
}

type printRow struct {
	sepBuf, lineBuf *bytes.Buffer
	view            *Table
	writer          io.Writer
}

func (tv *Table) startPrintRow() *printRow {
	row := &printRow{
		sepBuf:  tv.PaddingBuffer(),
		lineBuf: tv.PaddingBuffer(),
		view:    tv,
		writer:  tv.Out(),
	}
	row.sepBuf.WriteString("+")
	row.lineBuf.WriteString("|")
	return row
}

func (row *printRow) column(class, text string, c *Column, data interface{}) {
	for i := 0; i < c.Width; i++ {
		row.sepBuf.WriteString("-")
	}
	if c.Width > 0 {
		row.lineBuf.WriteString(row.view.Styling("table:"+class+":"+c.Field, wrapLen(ellipsis(text, c.Width), c.Width, c.Align), data))
	}
	row.sepBuf.WriteString("+")
	row.lineBuf.WriteString("|")
}

func (row *printRow) top() {
	fmt.Fprintln(row.writer, row.sepBuf.String())
}

func (row *printRow) end() {
	fmt.Fprintln(row.writer, row.lineBuf.String())
	fmt.Fprintln(row.writer, row.sepBuf.String())
}
