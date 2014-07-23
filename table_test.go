package cliview

import (
	"bytes"
	"strings"
	"testing"
)

func TestTablePrintWidthPercentage(t *testing.T) {
	buf := new(bytes.Buffer)
	tv := &Table{
		Output: Output{Writer: buf},
		Columns: []Column{
			Column{Title: "Column1", Field: "col1", Width: 10, Align: AlignRight},
			Column{Title: "Column2", Field: "col2", Width: -40, Align: AlignMiddle},
			Column{Title: "Column3", Field: "col3", Width: -60},
		},
		MaxWidth: 40,
	}
	tv.Print([]map[string]interface{}{
		map[string]interface{}{
			"col1": "123456",
			"col2": "abcde",
			"col3": "+-*/",
		},
	})
	result := buf.String()
	if result != ""+
		"+----------+----------+---------------+\n"+
		"|   Column1| Column2  |Column3        |\n"+
		"+----------+----------+---------------+\n"+
		"|    123456|  abcde   |+-*/           |\n"+
		"+----------+----------+---------------+\n" {
		t.Errorf("Unexpected output\n%v", result)
	}
}

func TestTablePrintEllipsis(t *testing.T) {
	buf := new(bytes.Buffer)
	tv := &Table{
		Output: Output{Writer: buf},
		Columns: []Column{
			Column{Title: "Column1", Field: "col1", MaxWidth: 5},
			Column{Title: "Column2", Field: "col2", Width: 2},
		},
	}
	tv.Print([]map[string]interface{}{
		map[string]interface{}{
			"col1": "123456",
			"col2": "abcde",
		},
	})
	result := buf.String()
	if result != ""+
		"+-----+--+\n"+
		"|Co...|C.|\n"+
		"+-----+--+\n"+
		"|12...|a.|\n"+
		"+-----+--+\n" {
		t.Errorf("Unexpected output\n%v", result)
	}
}

func TestTablePrintWidthAuto(t *testing.T) {
	buf := new(bytes.Buffer)
	tv := &Table{
		Output: Output{Writer: buf},
		Columns: []Column{
			Column{Title: "C1", Field: "col1", Align: AlignRight},
			Column{Title: "C2", Field: "col2", Width: -100},
		},
		MaxWidth: 9,
	}
	tv.Print([]map[string]interface{}{
		map[string]interface{}{
			"col1": "123456",
			"col2": nil,
		},
		map[string]interface{}{
			"col1": "12",
			"col2": nil,
		},
		map[string]interface{}{
			"col1": "123",
			"col2": nil,
		},
		map[string]interface{}{
			"col1": nil,
			"col2": nil,
		},
	})
	result := buf.String()
	if result != ""+
		"+------++\n"+
		"|    C1||\n"+
		"+------++\n"+
		"|123456||\n"+
		"+------++\n"+
		"|    12||\n"+
		"+------++\n"+
		"|   123||\n"+
		"+------++\n"+
		"|      ||\n"+
		"+------++\n" {
		t.Errorf("Unexpected output\n%v", result)
	}
}

func TestTablePrintPadding(t *testing.T) {
	buf := new(bytes.Buffer)
	tv := &Table{
		Output: Output{Writer: buf, Padding: 4},
		Columns: []Column{
			Column{Title: "Column1", Field: "col1", MaxWidth: 5},
			Column{Title: "Column2", Field: "col2", Width: 2},
		},
	}
	tv.Print([]map[string]interface{}{
		map[string]interface{}{
			"col1": "123456",
			"col2": "abcde",
		},
	})
	result := buf.String()
	if result != ""+
		"    +-----+--+\n"+
		"    |Co...|C.|\n"+
		"    +-----+--+\n"+
		"    |12...|a.|\n"+
		"    +-----+--+\n" {
		t.Errorf("Unexpected output\n%v", result)
	}
}

func TestTablePrintStyling(t *testing.T) {
	buf := new(bytes.Buffer)
	tv := &Table{
		Output: Output{
			Writer: buf,
			Styler: func(class, text string, data interface{}) string {
				if data == nil {
					return "<" + class + ">" + text + "</" + class + ">"
				}
				return text
			},
		},
		Columns: []Column{
			Column{Title: "C1", Field: "col1", Align: AlignRight},
			Column{Title: "C2", Field: "col2"},
		},
		MaxWidth: 9,
	}
	tv.Print([]map[string]interface{}{
		map[string]interface{}{
			"col1": "123456",
			"col2": nil,
		},
		map[string]interface{}{
			"col1": "12",
			"col2": nil,
		},
		map[string]interface{}{
			"col1": "123",
			"col2": nil,
		},
		map[string]interface{}{
			"col1": nil,
			"col2": nil,
		},
	})
	result := buf.String()
	if result != ""+
		"+------+--+\n"+
		"|    C1|C2|\n"+
		"+------+--+\n"+
		"|123456|<table:row:col2>  </table:row:col2>|\n"+
		"+------+--+\n"+
		"|    12|<table:row:col2>  </table:row:col2>|\n"+
		"+------+--+\n"+
		"|   123|<table:row:col2>  </table:row:col2>|\n"+
		"+------+--+\n"+
		"|<table:row:col1>      </table:row:col1>|<table:row:col2>  </table:row:col2>|\n"+
		"+------+--+\n" {
		t.Errorf("Unexpected output\n%v", result)
	}
}

func TestTableColumnFormatter(t *testing.T) {
	buf := new(bytes.Buffer)
	tv := &Table{
		Output: Output{Writer: buf, Padding: 4},
		Columns: []Column{
			Column{Title: "Column1", Field: "col1", MaxWidth: 5},
			Column{Title: "Column2", Field: "col2", Width: 2,
				Formatter: func(class string, data interface{}, formatter FormatterFunc) string {
					if strings.HasPrefix(class, "table:row:") {
						return ""
					} else {
						return formatter(class, data, nil)
					}
				},
			},
		},
	}
	tv.Print([]map[string]interface{}{
		map[string]interface{}{
			"col1": "123456",
			"col2": "abcde",
		},
	})
	result := buf.String()
	if result != ""+
		"    +-----+--+\n"+
		"    |Co...|C.|\n"+
		"    +-----+--+\n"+
		"    |12...|  |\n"+
		"    +-----+--+\n" {
		t.Errorf("Unexpected output\n%v", result)
	}
}
