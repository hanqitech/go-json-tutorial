package gjson

import (
	"encoding/json"
	"errors"
	"strconv"
)

func Unmarshal(data []byte, v any) error {
	if !json.Valid(data) {
		return errors.New("invalid json")
	}
	p := newParser(data)
	result := p.parse()

	vptr := v.(*any)
	*vptr = result
	return nil
}

type parser struct {
	data []byte
	len  int
	// index 指向当前字符串流的 char
	index int
}

func newParser(data []byte) *parser {
	// fmt.Println(string(data))
	p := new(parser)
	p.data = data
	p.len = len(data)
	return p
}

func (t *parser) curChar() rune {
	return rune(t.data[t.index])
}

func (t *parser) next() {
	t.index += 1
}

func (t *parser) end() bool {
	return t.len <= t.index
}

func (t *parser) isBlank() bool {
	switch t.curChar() {
	case ' ', '\n', '\t', '\r':
		return true
	default:
		return false
	}
}

func (t *parser) tryBool() (any, bool) {
	switch t.curChar() {
	case 't':
		t.index += 4
		return true, true
	case 'f':
		t.index += 5
		return false, true
	default:
		return nil, false
	}
}

func (t *parser) tryString() (any, bool) {
	var result any
	switch t.curChar() {
	case '"':
		next := t.index + 1
		for {
			if t.data[next] == '"' {
				result = string(t.data[t.index+1 : next])
				t.index = next + 1
				break
			}
			// 考虑 escape char
			if t.data[next] == '\\' && t.data[next+1] == '"' {
				next += 2
				continue
			}
			next += 1
		}
		return result, true
	default:
		return result, false
	}
}

func (t *parser) tryNum() (any, bool) {
	if (t.curChar() <= '9' && t.curChar() >= '0') || t.curChar() == '-' {
		// 需要利用终止符寻找数字的结尾
		next := t.index + 1
		for {
			if next == t.len {
				break
			}
			switch t.data[next] {
			case ',', ' ', '\t', '\n', '\r', ']', '}':
				// 可以 break 外层的 loop
				goto FOR
			default:
				next += 1
			}
		}
	FOR:

		num, _ := strconv.ParseFloat(string(t.data[t.index:next]), 64)
		t.index = next
		return num, true
	} else {
		return nil, false
	}
}

func (t *parser) tryNull() (any, bool) {
	switch t.curChar() {
	case 'n':
		t.index += 4
		return nil, true
	default:
		return nil, false
	}
}

func (t *parser) tryPrimitive() (any, bool) {
	var (
		result any
		ok     bool
	)
	if result, ok = t.tryBool(); ok {
	} else if result, ok = t.tryNum(); ok {
	} else if result, ok = t.tryString(); ok {
	} else if result, ok = t.tryNull(); ok {
	} else {
		return nil, false
	}
	return result, true
}

func (t *parser) pass(input ...rune) {
	for {
		if t.curChar() == ']' || t.curChar() == '}' {
			return
		}

		for _, v := range input {
			if t.curChar() == v {
				t.next()
				return
			}
		}
	}
}

func (t *parser) passComma() {
	t.passBlank()
	t.pass(',')
	t.passBlank()
}

func (t *parser) tryArray() (any, bool) {
	var (
		result []any
		ok     bool
		item   any
	)
	switch t.curChar() {
	case '[':
		t.next()
		for {
			if t.curChar() == ']' {
				t.next()
				break
			}
			if t.isBlank() {
				t.next()
				continue
			}
			if item, ok = t.tryPrimitive(); ok {
			} else if item, ok = t.tryArray(); ok {
			} else if item, ok = t.tryObject(); ok {
			} else {
				panic("array invalid")
			}
			result = append(result, item)

			// 跳过数组元素的分隔符
			t.passComma()
		}
		return result, true
	default:
		return nil, false
	}
}

func (t *parser) parseObjectKey() string {
	item, _ := t.tryString()
	return item.(string)
}

func (t *parser) passBlank() {
	for {
		if t.isBlank() {
			t.next()
			continue
		} else {
			return
		}
	}
}

func (t *parser) parseObjectValue() any {
	var (
		item any
		ok   bool
	)
	if item, ok = t.tryPrimitive(); ok {
	} else if item, ok = t.tryArray(); ok {
	} else if item, ok = t.tryObject(); ok {
	} else {
		panic("object value invalid")
	}
	return item
}

func (t *parser) tryObject() (any, bool) {
	var (
		result map[string]any = map[string]any{}
		// item   any
		key   string
		value any
	)
	switch t.curChar() {
	case '{':
		// 进入了 object 的内部
		t.next()
		for {
			// fmt.Printf("try object [%c]\n", t.curChar())
			if t.curChar() == '}' {
				// 离开了 object
				t.next()
				break
			}

			t.passBlank()
			key = t.parseObjectKey()
			t.passBlank()
			t.pass(':')

			t.passBlank()
			value = t.parseObjectValue()
			result[key] = value
			t.passBlank()
			t.passComma()
		}
		return result, true
	default:
		return nil, false
	}

}

func (t *parser) parse() any {
	var (
		ok   bool
		item any
	)
	for {
		// fmt.Printf("parse [%c]\n", t.curChar())
		if t.end() {
			break
		}
		if t.isBlank() {
			t.next()
			continue
		}

		if item, ok = t.tryPrimitive(); ok {
			return item
		} else if item, ok = t.tryArray(); ok {
			return item
		} else if item, ok = t.tryObject(); ok {
			return item
		} else {
			panic("invalid")
		}

	}

	return nil
}
