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

func TestTreePrintArrayMap(t *testing.T) {
	buf := new(bytes.Buffer)
	tv := &Tree{
		Output: Output{Writer: buf},
		Indent: DefaultIndent,
	}
	tv.Print([]map[string]interface{}{
		map[string]interface{}{
			"name": "Jackson",
			"key":  121,
		},
	})
	result := buf.String()
	if result != ""+
		"  - key: 121\n"+
		"    name: Jackson\n" {
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
	tv.Formatter = func(class string, data interface{}, formatter FormatterFunc) string {
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

func TestTreePrintFilterKeys(t *testing.T) {
	buf := new(bytes.Buffer)
	tv := &Tree{
		Output: Output{
			Writer: buf,
			Formatter: func(class string, data interface{}, formatter FormatterFunc) string {
				if class == "tree:key:map1" && data.(string) == "arr11" {
					return ""
				}
				return formatter(class, data, nil)
			},
		},
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
	})
	result := buf.String()
	if result != ""+
		"map1: \n"+
		"    arr12: \n"+
		"      - Hello\n"+
		"      - true\n"+
		"    map11: \n"+
		"    map12: \n"+
		"        key121: Hello\n" {
		t.Errorf("Unexpected output\n%v", result)
	}
}

func TestTreeKeyRanker(t *testing.T) {
	buf := new(bytes.Buffer)
	tv := &Tree{
		Output: Output{Writer: buf},
		Indent: DefaultIndent,
		KeyRanker: func(path, key string) int {
			return 255 - int(key[0])
		},
	}
	tv.Print(map[string]interface{}{
		"a": "a",
		"b": "b",
		"c": "c",
		"0": "0",
	})
	result := buf.String()
	if result != ""+
		"c: c\n"+
		"b: b\n"+
		"a: a\n"+
		"0: 0\n" {
		t.Errorf("Unexpected output\n%v", result)
	}
}

func TestTreeArrayKeyRanker(t *testing.T) {
	buf := new(bytes.Buffer)
	tv := &Tree{
		Output:    Output{Writer: buf},
		Indent:    DefaultIndent,
		KeyRanker: ArrayKeyRanker([]string{"b", "c", "a"}),
	}
	tv.Print(map[string]interface{}{
		"a": "a",
		"b": "b",
		"c": "c",
		"1": "1",
		"0": "0",
	})
	result := buf.String()
	if result != ""+
		"b: b\n"+
		"c: c\n"+
		"a: a\n"+
		"0: 0\n"+
		"1: 1\n" {
		t.Errorf("Unexpected output\n%v", result)
	}
}

func TestTreeSuffixArrayKeyRanker(t *testing.T) {
	buf := new(bytes.Buffer)
	tv := &Tree{
		Output: Output{Writer: buf},
		Indent: DefaultIndent,
		KeyRanker: SuffixArrayKeyRanker(map[string][]string{
			"":   {"b", "c", "a"},
			"l2": {"c", "a", "b"},
		}),
	}
	tv.Print(map[string]interface{}{
		"a": "a",
		"b": "b",
		"c": "c",
		"l2": map[string]interface{}{
			"a": "a",
			"b": "b",
			"c": "c",
		},
	})
	result := buf.String()
	if result != ""+
		"b: b\n"+
		"c: c\n"+
		"a: a\n"+
		"l2: \n"+
		"    c: c\n"+
		"    a: a\n"+
		"    b: b\n" {
		t.Errorf("Unexpected output\n%v", result)
	}
}

func TestTreePatternArrayKeyRanker(t *testing.T) {
	buf := new(bytes.Buffer)
	tv := &Tree{
		Output: Output{Writer: buf},
		Indent: DefaultIndent,
		KeyRanker: PatternArrayKeyRanker(map[string][]string{
			"":            {"b", "c", "a", "l2"},
			"l2$":         {"c", "a", "b"},
			"list/[^/]+$": {"x114", "l72", "l68"},
		}),
	}
	tv.Print(map[string]interface{}{
		"a": "a",
		"b": "b",
		"c": "c",
		"l2": map[string]interface{}{
			"a": "a",
			"b": "b",
			"c": "c",
		},
		"list": []interface{}{
			map[string]interface{}{
				"l4.1": "41",
				"x114": "114",
			},
			map[string]interface{}{
				"l68": "68",
				"l72": "72",
			},
		},
	})
	result := buf.String()
	if result != ""+
		"b: b\n"+
		"c: c\n"+
		"a: a\n"+
		"l2: \n"+
		"    c: c\n"+
		"    a: a\n"+
		"    b: b\n"+
		"list: \n"+
		"  - x114: 114\n"+
		"    l4.1: 41\n"+
		"  - l72: 72\n"+
		"    l68: 68\n" {
		t.Errorf("Unexpected output\n%v", result)
	}
}
