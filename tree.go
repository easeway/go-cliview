package cliview

import (
	"fmt"
	"io"
	"sort"
)

const (
	DefaultIndent = 4
)

type Tree struct {
	Output
	Indent int
}

func NewTree() *Tree {
	return &Tree{Indent: DefaultIndent}
}

func (tv *Tree) Print(obj interface{}) {
	tv.render(obj, "", tv.Out(), tv.Padding, false, false)
}

func (tv *Tree) render(obj interface{}, path string, w io.Writer, padding int, skipPadding, forCntr bool) {
	padBuf := PaddingBuffer(padding)
	padStr := padBuf.String()
	empty := false
	switch obj.(type) {
	case map[string]interface{}:
		if mapObj := obj.(map[string]interface{}); len(mapObj) == 0 {
			empty = true
		} else {
			keys := make([]string, 0, len(mapObj))
			for k := range mapObj {
				if tv.Format("tree:key:"+path, k) != "" {
					keys = append(keys, k)
				}
			}
			sort.Strings(keys)
			if skipPadding && !forCntr {
				fmt.Fprintln(w, "")
				skipPadding = false
			}
			for _, k := range keys {
				v := mapObj[k]
				keyStr := tv.Styling("tree:key:"+path, k, v)
				if skipPadding {
					fmt.Fprintf(w, "%s: ", keyStr)
					skipPadding = false
				} else {
					fmt.Fprintf(w, padStr+"%s: ", keyStr)
				}
				subpath := path
				if len(path) > 0 {
					subpath += "/"
				}
				subpath += k
				tv.render(v, subpath, w, padding+tv.Indent, true, false)
			}
		}
	case map[interface{}]interface{}:
		converted := make(map[string]interface{})
		for k, v := range obj.(map[interface{}]interface{}) {
			converted[fmt.Sprintf("%v", k)] = v
		}
		tv.render(converted, path, w, padding, skipPadding, forCntr)
	case []map[string]interface{}:
		data := obj.([]map[string]interface{})
		converted := make([]interface{}, len(data))
		for i, v := range data {
			converted[i] = v
		}
		tv.render(converted, path, w, padding, skipPadding, forCntr)
	case []map[interface{}]interface{}:
		data := obj.([]map[interface{}]interface{})
		converted := make([]interface{}, len(data))
		for i, v := range data {
			converted[i] = v
		}
		tv.render(converted, path, w, padding, skipPadding, forCntr)
	case []interface{}:
		if arrObj := obj.([]interface{}); len(arrObj) == 0 {
			empty = true
		} else {
			if skipPadding {
				fmt.Fprintln(w, "")
				if forCntr {
					padding += tv.Indent
					for i := 0; i < tv.Indent; i++ {
						padBuf.WriteString(" ")
					}
					padStr = padBuf.String()
				}
			}
			if len(padStr) >= tv.Padding+2 {
				padStr = padStr[0:len(padStr)-2] + "- "
			} else {
				padStr = PaddingBuffer(tv.Padding+tv.Indent-2).String() + "- "
				padding = tv.Padding + tv.Indent
			}
			for i, v := range arrObj {
				fmt.Fprintf(w, padStr)
				subpath := fmt.Sprintf("%v", i)
				if len(path) > 0 {
					subpath = path + "/" + subpath
				}
				tv.render(v, subpath, w, padding, true, true)
			}
		}
	default:
		if obj == nil {
			empty = true
		} else {
			if skipPadding {
				padStr = ""
			}
			class := "tree:val:" + path
			fmt.Fprintln(w, padStr+tv.Styling(class, tv.Format(class, obj), obj))
		}
	}

	if empty {
		if skipPadding {
			fmt.Fprintln(w, "")
		} else {
			fmt.Fprintln(w, padStr)
		}
	}
}
