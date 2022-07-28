package gjson

import (
	"reflect"
	"testing"
)

func TestMdParser(t *testing.T) {
	array := []string{"# a", "## b", "## c", "### d", "# e"}
	expected := []mdTitle{
		{level: "1", name: "a"},
		{level: "1.1", name: "b"},
		{level: "1.2", name: "c"},
		{level: "1.2.1", name: "d"},
		{level: "2", name: "e"},
	}
	result := markdownParser(array)
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("expected: %+v, result: %+v", expected, result)
	}

}
