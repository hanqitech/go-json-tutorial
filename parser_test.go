package gjson

import (
	"reflect"
	"testing"
)

func TestTokenizer(t *testing.T) {
	data := `
	true`
	expected := []any{true}
	p := newTokenizer([]byte(data))
	if err := p.parseTokens(); err != nil {
		t.Fatalf("tokenizer err %v", err)
	}
	if !reflect.DeepEqual(p.tokens, expected) {
		t.Fatalf("Unmarshal result %v, expected: %v", p.tokens, expected)
	}

	data = `
	false`
	expected = []any{false}
	p = newTokenizer([]byte(data))
	if err := p.parseTokens(); err != nil {
		t.Fatalf("tokenizer err %v", err)
	}
	if !reflect.DeepEqual(p.tokens, expected) {
		t.Fatalf("Unmarshal result %v, expected: %v", p.tokens, expected)
	}

	data = ` "Hello world!"	`
	expected = []any{"Hello world!"}
	p = newTokenizer([]byte(data))
	if err := p.parseTokens(); err != nil {
		t.Fatalf("tokenizer err %v", err)
	}
	if !reflect.DeepEqual(p.tokens, expected) {
		t.Fatalf("Unmarshal result %v, expected: %v", p.tokens, expected)
	}

	// with escape char
	data = ` "Hello \" world!"	`
	expected = []any{`Hello \" world!`}
	p = newTokenizer([]byte(data))
	if err := p.parseTokens(); err != nil {
		t.Fatalf("tokenizer err %v", err)
	}
	if !reflect.DeepEqual(p.tokens, expected) {
		t.Fatalf("Unmarshal result %v, expected: %v", p.tokens, expected)
	}

	data = `	42 `
	expected = []any{float64(42)}
	p = newTokenizer([]byte(data))
	if err := p.parseTokens(); err != nil {
		t.Fatalf("tokenizer err %v", err)
	}
	if !reflect.DeepEqual(p.tokens, expected) {
		t.Fatalf("str Unmarshal result %v, expected: %v", p.tokens, expected)
	}

	data = `	-42 `
	expected = []any{float64(-42)}
	p = newTokenizer([]byte(data))
	if err := p.parseTokens(); err != nil {
		t.Fatalf("tokenizer err %v", err)
	}
	if !reflect.DeepEqual(p.tokens, expected) {
		t.Fatalf("str Unmarshal result %v, expected: %v", p.tokens, expected)
	}

	data = `	42.123 `
	expected = []any{float64(42.123)}
	p = newTokenizer([]byte(data))
	if err := p.parseTokens(); err != nil {
		t.Fatalf("tokenizer err %v", err)
	}
	if !reflect.DeepEqual(p.tokens, expected) {
		t.Fatalf("str Unmarshal result %v, expected: %v", p.tokens, expected)
	}

	data = `null`
	expected = []any{nil}
	p = newTokenizer([]byte(data))
	if err := p.parseTokens(); err != nil {
		t.Fatalf("tokenizer err %v", err)
	}
	if !reflect.DeepEqual(p.tokens, expected) {
		t.Fatalf("str Unmarshal result %v, expected: %v", p.tokens, expected)
	}

	data = ` [116, 943, 234, 38793]`
	expected = []any{float64(116), float64(943), float64(234), float64(38793)}
	p = newTokenizer([]byte(data))
	if err := p.parseTokens(); err != nil {
		t.Fatalf("tokenizer err %v", err)
	}
	if !reflect.DeepEqual(p.tokens, expected) {
		t.Fatalf("str Unmarshal result %v, expected: %v", p.tokens, expected)
	}

	// 	data = `{
	// 		"Image": {
	// 				"Width":  800,
	// 				"Height": 600,
	// 				"Title":  "View from 15th Floor",
	// 				"Thumbnail": {
	// 						"Url":    "http://www.example.com/image/481989943",
	// 						"Height": 125,
	// 						"Width":  100
	// 				},
	// 				"Animated" : false,
	// 				"IDs": [116, 943, 234, 38793],
	// 				"Comment": nil
	// 			}
	// 	}
	// `
	// 	expected = []any{`"Image"`, ":", "{", `"Witdh"`, float64(800), ",", `"Height"`, float64(600),
	// 		",", `"Title"`, ":", `"View from 15th Floor"`, ",",
	// 		`"Thumbnail"`, ":", "{", `"Url"`, ":", `"http://www.example.com/image/481989943"`,
	// 		",", `"Height"`, ":", float64(125), ",", `"Width"`, float64(100), `"Animated"`, ":", "false", ",",
	// 		`"IDs"`, ":", "[", float64(116), float64(943), float64(234), float64(38793), "]", ",",
	// 		`"Comment"`, ":", "nil",
	// 	}
	// 	p = newTokenizer([]byte(data))
	// 	if err := p.parseTokens(); err != nil {
	// 		t.Fatalf("tokenizer err %v", err)
	// 	}

	// 	if reflect.DeepEqual(p.tokens, expected) {
	// 		t.Fatalf("str Unmarshal result %v, expected: %v", p.tokens, expected)
	// 	}
}
func TestAtomValue(t *testing.T) {
	t.Run("null", func(t *testing.T) {
		strData := `"Hello world!"`
		var strResult any
		if err := Unmarshal([]byte(strData), &strResult); err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if strResult != "Hello world!" {
			t.Fatalf("str Unmarshal result %s, expected: %s", strData, strResult)
		}
	})
}

func TestBasicInterfaceUnmarshal(t *testing.T) {
	strData := `"Hello world!"`
	var strResult any
	if err := Unmarshal([]byte(strData), &strResult); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if strResult != "Hello world!" {
		t.Fatalf("str Unmarshal result %s, expected: %s", strData, strResult)
	}
	// fmt.Println("string result", strResult)

	intData := "42"
	var intResult any
	if err := Unmarshal([]byte(intData), &intResult); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if intResult != float64(42) {
		t.Fatalf("int Unmarshal err: [%v] %T", intResult, intResult)
	}

	boolData := "true"
	var boolResult any
	if err := Unmarshal([]byte(boolData), &boolResult); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if !boolResult.(bool) {
		t.Fatalf("bool Unmarshal err")
	}

	nullData := "null"
	var nullResult any = 1
	if err := Unmarshal([]byte(nullData), &nullResult); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if nullResult != nil {
		t.Fatalf("null Unmarshal err")
	}

	objectData := `{
		"Image": {
				"Width":  800,
				"Height": 600,
				"Title":  "View from 15th Floor",
				"Thumbnail": {
						"Url":    "http://www.example.com/image/481989943",
						"Height": 125,
						"Width":  100
				},
				"Animated" : false,
				"IDs": [116, 943, 234, 38793]
			}
	}
`

	var expect map[string]any = map[string]any{
		"Image": map[string]any{
			"Width":    float64(800),
			"Height":   float64(600),
			"Title":    "View from 15th Floor",
			"Animated": false,
			"IDs":      []any{float64(116), float64(943), float64(234), float64(38793)},
			"Thumbnail": map[string]any{
				"Url":    "http://www.example.com/image/481989943",
				"Height": float64(125),
				"Width":  float64(100),
			},
		},
	}

	var objectResult any
	if err := Unmarshal([]byte(objectData), &objectResult); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if !reflect.DeepEqual(objectResult, expect) {
		t.Fatalf("object Unmarshal err, result: %+v\n, expected: %+v", objectResult, expect)
	}

	arrayData := `[
		{
			 "precision": "zip",
			 "Latitude":  37.7668,
			 "Longitude": -122.3959,
			 "Address":   "",
			 "City":      "SAN FRANCISCO",
			 "State":     "CA",
			 "Zip":       "94107",
			 "Country":   "US"
		},
		{
			 "precision": "zip",
			 "Latitude":  37.371991,
			 "Longitude": -122.026020,
			 "Address":   "",
			 "City":      "SUNNYVALE",
			 "State":     "CA",
			 "Zip":       "94085",
			 "Country":   "US"
		}
	]`
	var arrayExpected []any = []any{
		map[string]any{
			"precision": "zip",
			"Latitude":  float64(37.7668),
			"Longitude": float64(-122.3959),
			"Address":   "",
			"City":      "SAN FRANCISCO",
			"State":     "CA",
			"Zip":       "94107",
			"Country":   "US",
		},
		map[string]any{
			"precision": "zip",
			"Latitude":  float64(37.371991),
			"Longitude": float64(-122.026020),
			"Address":   "",
			"City":      "SUNNYVALE",
			"State":     "CA",
			"Zip":       "94085",
			"Country":   "US",
		},
	}
	var arrayResult any
	if err := Unmarshal([]byte(arrayData), &arrayResult); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if !reflect.DeepEqual(arrayResult, arrayExpected) {
		t.Fatalf("array Unmarshal err, result: %+v\n, expected: %+v", arrayResult, arrayExpected)
	}
}
