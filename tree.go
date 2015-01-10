package cliview

import (
	"fmt"
	"io"
	"math"
	"regexp"
	"sort"
	"strings"
)

const (
	DefaultIndent = 4
)

type KeyRankFunc func(path, key string) int

func ArrayKeyRanker(keys []string) KeyRankFunc {
	return func(path, key string) int {
		for n, k := range keys {
			if k == key {
				return n
			}
		}
		return -1
	}
}

func PatternArrayKeyRanker(dict map[string][]string) KeyRankFunc {
	return func(path, key string) int {
		for pattern, keys := range dict {
			if path == "" && pattern == "" {
				return ArrayKeyRanker(keys)(path, key)
			} else if pattern != "" {
				if matched, err := regexp.MatchString(pattern, path); matched && err == nil {
					return ArrayKeyRanker(keys)(path, key)
				}
			}
		}
		return -1
	}
}

func SuffixArrayKeyRanker(dict map[string][]string) KeyRankFunc {
	return func(path, key string) int {
		for suffix, keys := range dict {
			if (path == "" && suffix == "") ||
				(suffix != "" && strings.HasSuffix(path, suffix)) {
				return ArrayKeyRanker(keys)(path, key)
			}
		}
		return -1
	}
}

type Tree struct {
	Output
	Indent    int
	KeyRanker KeyRankFunc
}

func NewTree() *Tree {
	return &Tree{Indent: DefaultIndent}
}

func (tv *Tree) Print(obj interface{}) {
	tv.render(obj, "", tv.Out(), tv.Padding, false, false)
}

func (tv *Tree) RankKey(path, key string) uint {
	if tv.KeyRanker != nil {
		if iRank := tv.KeyRanker(path, key); iRank >= 0 {
			return uint(iRank)
		}
	}
	return math.MaxUint32
}

type keyRank struct {
	key  string
	rank uint
}

type keySorter struct {
	keys []*keyRank
}

func (s *keySorter) Len() int {
	return len(s.keys)
}

func (s *keySorter) Swap(i, j int) {
	s.keys[i], s.keys[j] = s.keys[j], s.keys[i]
}

func (s *keySorter) Less(i, j int) bool {
	if s.keys[i].rank == s.keys[j].rank {
		return s.keys[i].key < s.keys[j].key
	}
	return s.keys[i].rank < s.keys[j].rank
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
			keys := &keySorter{keys: make([]*keyRank, 0, len(mapObj))}
			for k := range mapObj {
				if tv.Format("tree:key:"+path, k) != "" {
					rank := tv.RankKey(path, k)
					keys.keys = append(keys.keys, &keyRank{key: k, rank: rank})
				}
			}
			sort.Sort(keys)
			if skipPadding && !forCntr {
				fmt.Fprintln(w, "")
				skipPadding = false
			}
			for _, kr := range keys.keys {
				v := mapObj[kr.key]
				keyStr := tv.Styling("tree:key:"+path, kr.key, v, nil)
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
				subpath += kr.key
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
			fmt.Fprintln(w, padStr+tv.Styling(class, tv.Format(class, obj), obj, nil))
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
