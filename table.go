package cliview

import (
	"bytes"
	"fmt"
	"io"
	"unicode/utf8"
)

const (
	AlignLeft   = 0
	AlignMiddle = 1
	AlignRight  = 2
)

var (
	runesBorderFull = []rune{
		'\u2554', //c9 top: LT
		'\u2550', //cd top: T
		'\u2564', //d1 top: TS
		'\u2557', //bb top: RT
		'\u2551', //ba row: L
		'\u2502', //b3 row: S
		'\u2551', //ba row: R
		'\u255a', //c8 bot: LB
		'\u2550', //cd bot: B
		'\u2567', //cf bot: BS
		'\u255d', //bc bot: RB
		'\u255f', //c7 row-splitter: LS
		'\u2500', //c4 row-splitter: C
		'\u253c', //c5 row-splitter: CS
		'\u2562', //b6 row-splitter: RS
		'\u2551', //ba head: L
		'\u2502', //b3 head: S
		'\u2551', //ba head: R
		'\u2560', //cc head-splitter: LS
		'\u2550', //cd head-splitter: C
		'\u256a', //d8 head-splitter: CS
		'\u2563', //b9 head-splitter: RS
	}

	BorderFull = string(runesBorderFull)

	runesBorderCompact = []rune{
		'\u2554', //c9 top: LT
		'\u2550', //cd top: T
		'\u2564', //d1 top: TS
		'\u2557', //bb top: RT
		'\u2551', //ba row: L
		'\u2502', //b3 row: S
		'\u2551', //ba row: R
		'\u255a', //c8 bot: LB
		'\u2550', //cd bot: B
		'\u2567', //cf bot: BS
		'\u255d', //bc bot: RB
		'\u0000', //00 row-splitter: LS
		'\u2500', //c4 row-splitter: C
		'\u253c', //c5 row-splitter: CS
		'\u2562', //b6 row-splitter: RS
		'\u2551', //ba head: L
		'\u2502', //b3 head: S
		'\u2551', //ba head: R
		'\u2560', //cc head-splitter: LS
		'\u2550', //cd head-splitter: C
		'\u256a', //d8 head-splitter: CS
		'\u2563', //b9 head-splitter: RS
	}

	BorderCompact = string(runesBorderCompact)
)

type Column struct {
	Title     string // title to be displayed
	Field     string // field name for retrieving data
	Width     int    // column width, >0 fixed, =0 auto, <0 percentage
	MaxWidth  int    // maximum column width
	Align     int
	Fetcher   func(column Column, row map[string]interface{}) interface{}
	Formatter FormatterFunc
	Styler    StylerFunc
}

type Table struct {
	Output

	Border string // border characters
	// top: LT T TS RT
	// row: L S R
	// bot: LB B BS RB
	// row-splitter: LS C CS RS
	// head: LH SH RH
	// head-splitter: LSH CH CSH RSH
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
			width := textWidth(tv.Columns[i].Title)
			for _, v := range data {
				valLen := textWidth(tv.formatCell("table:row:", tv.Columns[i], v))
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

	chars := charsInString(tv.Border)
	if len(chars) < 11 {
		chars = runesBorderFull
	}

	// print head
	rowSepOff := 11
	headRowOff := 15
	headSepOff := 18
	if len(chars) < 15 {
		rowSepOff = -1
		headSepOff = -1
		headRowOff = 4
	} else if len(chars) < 18 {
		headRowOff = 4
		headSepOff = rowSepOff
	} else if len(chars) < 22 {
		headSepOff = rowSepOff
	}

	if rowSepOff >= 0 && chars[rowSepOff] == 0 {
		rowSepOff = -1
	}
	if headSepOff >= 0 && chars[headSepOff] == 0 {
		headSepOff = -1
	}

	row := tv.startPrintRow(chars, 0, headRowOff)
	for i, c := range tv.columns {
		row.column("head", c.Title, i, c.Title)
	}
	row.end()

	// print rows
	sepOff := headSepOff
	for _, d := range data {
		row = tv.startPrintRow(chars, sepOff, 4)
		sepOff = rowSepOff
		for i, c := range tv.columns {
			val := d[c.Field]
			if c.Width > 0 {
				valStr := tv.formatCell("table:row:", tv.Columns[i], d)
				row.column("row", valStr, i, val)
			} else {
				row.column("row", "", i, val)
			}
		}
		row.end()
	}

	row = tv.startPrintRow(chars, 7, -1)
	for i, c := range tv.columns {
		row.column("", c.Title, i, c.Title)
	}
	row.end()
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
	if textWidth(text) <= width {
		return text
	}
	chars := charsInString(text)
	if width <= 3 {
		text = string(chars[0:1])
		for i := 1; i < width; i++ {
			text += "."
		}
	} else {
		text = string(chars[0:width-3]) + "..."
	}
	return text
}

func wrapLen(text string, width int, align int) string {
	if i := textWidth(text); i < width {
		buf := new(bytes.Buffer)
		for ; i < width; i++ {
			buf.WriteString(" ")
		}
		pads := buf.String()
		switch align {
		case AlignRight:
			return pads + text
		case AlignMiddle:
			left := (width - textWidth(text)) / 2
			return pads[0:left] + text + pads[left:len(pads)]
		default:
			return text + pads
		}
	}
	return text
}

type printRow struct {
	bufSep, bufRow *bytes.Buffer
	view           *Table
	border         []rune
	offSep, offRow int
	writer         io.Writer
}

func (tv *Table) startPrintRow(chars []rune, offSep, offRow int) *printRow {
	row := &printRow{
		bufSep: tv.PaddingBuffer(),
		bufRow: tv.PaddingBuffer(),
		view:   tv,
		border: chars,
		offSep: offSep,
		offRow: offRow,
		writer: tv.Out(),
	}
	return row
}

func (row *printRow) column(class, text string, col int, data interface{}) {
	addSep := 2
	addRow := 1
	if col == 0 {
		addSep = 0
		addRow = 0
	}
	c := &row.view.columns[col]

	if row.offSep >= 0 {
		row.bufSep.WriteRune(row.border[row.offSep+addSep])
		if c.Width > 0 {
			for i := 0; i < c.Width; i++ {
				row.bufSep.WriteRune(row.border[row.offSep+1])
			}
		}
	}
	if row.offRow >= 0 {
		row.bufRow.WriteRune(row.border[row.offRow+addRow])
		if c.Width > 0 {
			row.bufRow.WriteString(row.view.Styling("table:"+class+":"+c.Field, wrapLen(ellipsis(text, c.Width), c.Width, c.Align), data, row.view.Columns[col].Styler))
		}
	}
}

func (row *printRow) end() {
	if row.offSep >= 0 {
		row.bufSep.WriteRune(row.border[row.offSep+3])
		fmt.Fprintln(row.writer, row.bufSep.String())
	}
	if row.offRow >= 0 {
		row.bufRow.WriteRune(row.border[row.offRow+2])
		fmt.Fprintln(row.writer, row.bufRow.String())
	}
}

func charsInString(text string) []rune {
	chars := make([]rune, 0)
	for str := text; len(str) > 0; {
		if char, size := utf8.DecodeRuneInString(str); size <= 0 {
			break
		} else {
			chars = append(chars, char)
			str = str[size:]
		}
	}
	return chars
}

func textWidth(text string) int {
	return utf8.RuneCountInString(text)
}
