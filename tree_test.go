package cliview

import (
	"bytes"
	"fmt"
	"testing"
)

func TestTreePrintValue(t *testing.T) {
	buf := new(bytes.Buffer)
	tv := &Tree{
		Output: Output{Writer: buf},
		Indent: DefaultIndent,
	}
	tv.Print("Hello")
	tv.Print(true)
	tv.Print(10)
	result := buf.String()
	if result != "Hello\ntrue\n10\n" {
		t.Errorf("Unexpected output\n%v", result)
	}
}

func TestTreePrintMap(t *testing.T) {
	buf := new(bytes.Buffer)
	tv := &Tree{
		Output: Output{Writer: buf},
		Indent: DefaultIndent,
	}
	tv.Print(map[string]interface{}{
		"key1": "Hello",
		"key2": true,
		"key3": 100,
		"key4": nil,
	})
	result := buf.String()
	if result != "key1: Hello\nkey2: true\nkey3: 100\nkey4: \n" {
		t.Errorf("Unexpected output\n%v", result)
	}
}

func TestTreePrintEmptyMap(t *testing.T) {
	buf := new(bytes.Buffer)
	tv := &Tree{
		Output: Output{Writer: buf},
		Indent: DefaultIndent,
	}
	tv.Print(map[string]interface{}{})
	result := buf.String()
	if result != "\n" {
		t.Errorf("Unexpected output\n%v", result)
	}
}

func TestTreePrintEmptyArray(t *testing.T) {
	buf := new(bytes.Buffer)
	tv := &Tree{
		Output: Output{Writer: buf},
		Indent: DefaultIndent,
	}
	tv.Print([]interface{}{})
	result := buf.String()
	if result != "\n" {
		t.Errorf("Unexpected output\n%v", result)
	}
}

func TestTreePrintMixed(t *testing.T) {
	buf := new(bytes.Buffer)
	tv := &Tree{
		Output: Output{Writer: buf},
		Indent: DefaultIndent,
	}
	tv.Print(map[string]interface{}{
		"map1": map[string]interface{}{
			"map11": map[string]interface{}{},
			"arr11": []interface{}{},
			"map12": map[string]interface{}{
				"key121": "Hello",
			},
			"arr12": []interface{}{
				"Hello", true,
			},
		},
		"arr1": []interface{}{
			map[string]interface{}{
				"key121": "Hello",
				"arr122": []interface{}{
					"Hello", 122,
				},
			},
			[]interface{}{
				[]interface{}{
					"NestedArr1", 12,
				},
			},
		},
	})
	result := buf.String()
	if result != ""+
		"arr1: \n"+
		"  - arr122: \n"+
		"      - Hello\n"+
		"      - 122\n"+
		"    key121: Hello\n"+
		"  - \n"+
		"      - \n"+
		"          - NestedArr1\n"+
		"          - 12\n"+
		"map1: \n"+
		"    arr11: \n"+
		"    arr12: \n"+
		"      - Hello\n"+
		"      - true\n"+
		"    map11: \n"+
		"    map12: \n"+
		"        key121: Hello\n" {
		t.Errorf("Unexpected output\n%v", result)
	}
}

func TestTreePrintFormatAndStyling(t *testing.T) {
	buf := new(bytes.Buffer)
	tv := &Tree{
		Output: Output{Writer: buf},
		Indent: DefaultIndent,
	}
	tv.Styler = func(class, text string, data interface{}) string {
		switch data.(type) {
		case bool:
			return "<ok>" + text + "</ok>"
		case int:
			return "<i:" + class + ">" + text + "</i:" + class + ">"
		default:
			return text
		}
	}
	tv.Formatter = func(class string, data interface{}) string {
		switch data.(type) {
		case string:
			return "\"" + data.(string) + "\""
		default:
			return fmt.Sprintf("%v", data)
		}
	}

	tv.Print(map[string]interface{}{
		"map1": map[string]interface{}{
			"map11": map[string]interface{}{},
			"arr11": []interface{}{},
			"map12": map[string]interface{}{
				"key121": "Hello",
			},
			"arr12": []interface{}{
				"Hello", true,
			},
		},
		"arr1": []interface{}{
			map[string]interface{}{
				"key121": "Hello",
				"arr122": []interface{}{
					"Hello", 122,
				},
			},
			[]interface{}{
				[]interface{}{
					"NestedArr1", 12,
				},
			},
		},
	})
	result := buf.String()
	if result != ""+
		"arr1: \n"+
		"  - arr122: \n"+
		"      - \"Hello\"\n"+
		"      - <i:tree:val:arr1/0/arr122/1>122</i:tree:val:arr1/0/arr122/1>\n"+
		"    key121: \"Hello\"\n"+
		"  - \n"+
		"      - \n"+
		"          - \"NestedArr1\"\n"+
		"          - <i:tree:val:arr1/1/0/1>12</i:tree:val:arr1/1/0/1>\n"+
		"map1: \n"+
		"    arr11: \n"+
		"    arr12: \n"+
		"      - \"Hello\"\n"+
		"      - <ok>true</ok>\n"+
		"    map11: \n"+
		"    map12: \n"+
		"        key121: \"Hello\"\n" {
		t.Errorf("Unexpected output\n%v", result)
	}
}
