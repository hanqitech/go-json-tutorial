package gjson

import (
	"encoding/json"
	"errors"
	"strconv"
)

func Unmarshal(data []byte, v any) error {
	if json.Valid(data) {
		return errors.New("invalid json")
	}

	return nil
}

type tokenizer struct {
	data []byte
	len  int
	// index 是指向当前字符串流的 char， endIndex 指向当前 token 结尾的 char
	index    int
	endIndex int

	tokens []any
}

func newTokenizer(data []byte) *tokenizer {
	p := new(tokenizer)
	p.data = data
	p.len = len(data)
	return p
}

func (t *tokenizer) curChar() rune {
	return rune(t.data[t.index])
}

func (t *tokenizer) next() {
	t.index += 1
}

func (t *tokenizer) end() bool {
	return t.len <= t.index
}

func (t *tokenizer) isBlank() bool {
	switch t.curChar() {
	case ' ', '\n', '\t', 'r':
		return true
	default:
		return false
	}
}

func (t *tokenizer) tryBool() bool {
	switch t.curChar() {
	case 't':
		t.tokens = append(t.tokens, true)
		t.index += 4
		return true
	case 'f':
		t.tokens = append(t.tokens, false)
		t.index += 5
		return true
	default:
		return false
	}
}

func (t *tokenizer) tryString() bool {
	switch t.curChar() {
	case '"':
		t.tokens = append(t.tokens, "\"")
		next := t.index + 1
		for {
			if t.data[next] == '"' {
				t.endIndex = next
				t.tokens = append(t.tokens, string(t.data[t.index+1:t.endIndex]))
				t.index = next + 1
				t.tokens = append(t.tokens, "\"")
				break
			}
			// 考虑 escape char
			if t.data[next] == '\\' && t.data[next+1] == '"' {
				next += 2
				continue
			}
			next += 1
		}
		return true
	default:
		return false
	}
}

func (t *tokenizer) tryNum() bool {
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
		t.tokens = append(t.tokens, num)
		t.index = next
		return true
	} else {
		return false
	}
}

func (t *tokenizer) tryNull() bool {
	switch t.curChar() {
	case 'n':
		t.tokens = append(t.tokens, nil)
		t.index += 4
		return true
	default:
		return false
	}
}

func (t *tokenizer) tryPrimitive() bool {
	return t.tryBool() || t.tryString() || t.tryNum() || t.tryNull()
}

func (t *tokenizer) tryComma() {
	for {
		if t.isBlank() {
			t.next()
			continue
		}

		if t.curChar() == ']' {
			return
		}

		if t.curChar() == ',' {
			t.tokens = append(t.tokens, ",")
			t.next()
			return
		}
	}
}

func (t *tokenizer) tryArray() bool {
	switch t.curChar() {
	case '[':
		t.tokens = append(t.tokens, "[")
		t.next()
		for {
			if t.curChar() == ']' {
				t.tokens = append(t.tokens, "]")
				t.next()
				break
			}
			if t.isBlank() {
				t.next()
				continue
			}

			if t.tryPrimitive() {
			} else if t.tryArray() {
			} else {
				panic("array object not implemented")
			}
			// 解析数据元素的分隔符
			t.tryComma()
		}
		return true
	default:
		return false
	}
}

func (t *tokenizer) parseTokens() error {
	for {
		if t.end() {
			break
		}
		if t.isBlank() {
			t.next()
			continue
		}

		if t.tryPrimitive() {
			continue
		} else if t.tryArray() {
			continue
		} else {
			panic("not implemented")

		}
		// if t.tryBool() {
		// 	continue
		// } else if t.tryString() {
		// 	continue
		// } else if t.tryNum() {
		// 	continue
		// } else if t.tryNull() {
		// 	continue
		// } else{
		// 	panic("not implemented")
		// }

	}

	return nil
}
