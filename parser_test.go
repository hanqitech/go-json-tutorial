package gjson

import (
	"reflect"
	"testing"
)

func TestParser(t *testing.T) {
	data := `
	true`
	var result any
	var expected any = true
	if err := Unmarshal([]byte(data), &result); err != nil {
		t.Fatalf("bool err %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Unmarshal result %v, expected: %v", result, expected)
	}


	data = `
	false`
	expected = false
	if err := Unmarshal([]byte(data), &result); err != nil {
		t.Fatalf("bool err %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Unmarshal result %v, expected: %v", result, expected)
	}

	data = ` "Hello world! 你好，世界！"	`
	expected = "Hello world! 你好，世界！"
	if err := Unmarshal([]byte(data), &result); err != nil {
		t.Fatalf("str err %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Unmarshal result %v, expected: %v", result, expected)
	}

	// with escape char
	data = ` "Hello \" world!"	`
	expected = `Hello \" world!`
	if err := Unmarshal([]byte(data), &result); err != nil {
		t.Fatalf("str err %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Unmarshal result %v, expected: %v", result, expected)
	}

	data = `	42 `
	expected = float64(42)
	if err := Unmarshal([]byte(data), &result); err != nil {
		t.Fatalf("num err %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Unmarshal result %v, expected: %v", result, expected)
	}

	data = `	-42 `
	expected = float64(-42)
	if err := Unmarshal([]byte(data), &result); err != nil {
		t.Fatalf("num err %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Unmarshal result %v, expected: %v", result, expected)
	}

	data = `	42.123 `
	expected = float64(42.123)
	if err := Unmarshal([]byte(data), &result); err != nil {
		t.Fatalf("num err %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Unmarshal result %v, expected: %v", result, expected)
	}

	data = `null`
	expected = nil
	if err := Unmarshal([]byte(data), &result); err != nil {
		t.Fatalf("nil err %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Unmarshal result %v, expected: %v", result, expected)
	}

	data = ` [116, 943, 234, 38793]`
	expected = []any{float64(116), float64(943), float64(234), float64(38793)}
	if err := Unmarshal([]byte(data), &result); err != nil {
		t.Fatalf("array err %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Unmarshal result %v, expected: %v", result, expected)
	}

	data = ` [116, true, "test", null]`
	expected = []any{float64(116), true, "test", nil}
	if err := Unmarshal([]byte(data), &result); err != nil {
		t.Fatalf("array err %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Unmarshal result %v, expected: %v", result, expected)
	}

	data = ` [116, [117, [118]]]`
	expected = []any{float64(116), []any{float64(117), []any{float64(118)}}}
	if err := Unmarshal([]byte(data), &result); err != nil {
		t.Fatalf("array err %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Unmarshal result %v, expected: %v", result, expected)
	}

	data = `
	[
   {
      "keya":123,
      "keyb":234
   },
   {
      "keya":123,
      "keyb":234
   }
]
`
	expected = []any{
		map[string]any{
			"keya": float64(123),
			"keyb": float64(234),
		},
		map[string]any{
			"keya": float64(123),
			"keyb": float64(234),
		},
	}
	if err := Unmarshal([]byte(data), &result); err != nil {
		t.Fatalf("object err %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Unmarshal result %v, expected: %v", result, expected)
	}

	data = `
	{
		"keyNum" : 123,
		"keyNum2" : -123.234,
		"keyStr": "hello world",
		"keyBool": true,
		"keyNull": null,
		"keyArray": [1, 2, 3]
	}`
	expected = map[string]any{
		"keyNum":   float64(123),
		"keyNum2":  float64(-123.234),
		"keyStr":   "hello world",
		"keyBool":  true,
		"keyNull":  nil,
		"keyArray": []any{float64(1), float64(2), float64(3)},
	}
	if err := Unmarshal([]byte(data), &result); err != nil {
		t.Fatalf("object err %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Unmarshal result %v, expected: %v", result, expected)
	}

	data = `
	{
		"keyObject1":{
			 "keyObject2":{
					"key2":2
			 }
		}
 }
 `
	expected = map[string]any{
		"keyObject1": map[string]any{
			"keyObject2": map[string]any{
				"key2": float64(2),
			},
		},
	}
	if err := Unmarshal([]byte(data), &result); err != nil {
		t.Fatalf("object err %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Unmarshal result %v, expected: %v", result, expected)
	}
}
func TestRFCExample(t *testing.T) {
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
