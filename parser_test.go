package main

import (
	"reflect"
	"testing"
)

func TestParseJSON_BasicTypes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{"null", "null", nil},
		{"true", "true", true},
		{"false", "false", false},
		{"integer", "42", float64(42)},
		{"negative integer", "-17", float64(-17)},
		{"float", "3.14", 3.14},
		{"negative float", "-2.5", -2.5},
		{"exponential", "1.5e2", 150.0},
		{"exponential with plus", "1.5e+2", 150.0},
		{"exponential with minus", "1.5e-2", 0.015},
		{"simple string", `"hello"`, "hello"},
		{"empty string", `""`, ""},
		{"string with spaces", `"hello world"`, "hello world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseJSON(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v (%T), got %v (%T)", tt.expected, tt.expected, result, result)
			}
		})
	}
}

func TestParseJSON_StringEscapes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"escaped quote", `"He said \"Hello\""`, `He said "Hello"`},
		{"escaped backslash", `"C:\\Users\\file.txt"`, `C:\Users\file.txt`},
		{"escaped newline", `"line1\nline2"`, "line1\nline2"},
		{"escaped tab", `"name\tvalue"`, "name\tvalue"},
		{"escaped carriage return", `"line1\rline2"`, "line1\rline2"},
		{"escaped backspace", `"test\bword"`, "test\bword"},
		{"escaped form feed", `"page1\fpage2"`, "page1\fpage2"},
		{"escaped forward slash", `"https:\/\/example.com"`, "https://example.com"},
		{"multiple escapes", `"\"hello\"\n\t\"world\""`, "\"hello\"\n\t\"world\""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseJSON(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestParseJSON_Arrays(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []interface{}
	}{
		{"empty array", "[]", []interface{}{}},
		{"single element", "[1]", []interface{}{float64(1)}},
		{"multiple numbers", "[1, 2, 3]", []interface{}{float64(1), float64(2), float64(3)}},
		{"mixed types", `[1, "hello", true, null]`, []interface{}{float64(1), "hello", true, nil}},
		{"nested arrays", "[[1, 2], [3, 4]]", []interface{}{
			[]interface{}{float64(1), float64(2)},
			[]interface{}{float64(3), float64(4)},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseJSON(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestParseJSON_Objects(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]interface{}
	}{
		{"empty object", "{}", map[string]interface{}{}},
		{"single pair", `{"key": "value"}`, map[string]interface{}{"key": "value"}},
		{"multiple pairs", `{"name": "John", "age": 30}`, map[string]interface{}{
			"name": "John",
			"age":  float64(30),
		}},
		{"mixed types", `{"str": "hello", "num": 42, "bool": true, "null": null}`, map[string]interface{}{
			"str":  "hello",
			"num":  float64(42),
			"bool": true,
			"null": nil,
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseJSON(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestParseJSON_NestedStructures(t *testing.T) {
	input := `{
		"user": {
			"name": "Alice",
			"details": {
				"age": 25,
				"hobbies": ["reading", "coding"]
			}
		},
		"active": true
	}`

	expected := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "Alice",
			"details": map[string]interface{}{
				"age":     float64(25),
				"hobbies": []interface{}{"reading", "coding"},
			},
		},
		"active": true,
	}

	result, err := ParseJSON(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestParseJSON_WithWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{"object with spaces", ` { "key" : "value" } `, map[string]interface{}{"key": "value"}},
		{"array with newlines", "[\n  1,\n  2,\n  3\n]", []interface{}{float64(1), float64(2), float64(3)}},
		{"nested with tabs", "{\n\t\"outer\": {\n\t\t\"inner\": \"value\"\n\t}\n}", map[string]interface{}{
			"outer": map[string]interface{}{"inner": "value"},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseJSON(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestParseJSON_InvalidJSON(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty input", ""},
		{"unclosed object", `{"key": "value"`},
		{"unclosed array", `[1, 2, 3`},
		{"unclosed string", `"unclosed string`},
		{"missing colon", `{"key" "value"}`},
		{"missing comma in object", `{"a": 1 "b": 2}`},
		{"missing comma in array", `[1 2 3]`},
		{"trailing comma in object", `{"key": "value",}`},
		{"trailing comma in array", `[1, 2, 3,]`},
		{"invalid boolean", "truee"},
		{"invalid null", "nulll"},
		{"invalid number", "1.2.3"},
		{"unquoted key", `{key: "value"}`},
		{"single quotes", `{'key': 'value'}`},
		{"invalid escape", `"invalid \x escape"`},
		{"unexpected character", `{"key": @}`},
		{"extra content", `{"key": "value"} extra`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseJSON(tt.input)
			if err == nil {
				t.Errorf("expected error for invalid JSON: %s", tt.input)
			}
		})
	}
}

func TestParseJSON_Numbers(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{"zero", "0", 0},
		{"positive integer", "123", 123},
		{"negative integer", "-456", -456},
		{"positive decimal", "123.45", 123.45},
		{"negative decimal", "-67.89", -67.89},
		{"scientific notation positive", "1e5", 100000},
		{"scientific notation negative", "1e-5", 0.00001},
		{"scientific notation with decimal", "1.23e4", 12300},
		{"scientific notation negative exponent", "1.23e-4", 0.000123},
		{"capital E", "1E5", 100000},
		{"explicit positive exponent", "1e+5", 100000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseJSON(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestParseJSON_InvalidNumbers(t *testing.T) {
	tests := []string{
		"01",    // leading zero
		"1.",    // trailing decimal point
		"1e",    // incomplete exponent
		"1e+",   // incomplete exponent
		"1e-",   // incomplete exponent
		".5",    // missing leading digit
		"1.2.3", // multiple decimal points
		"-",     // just minus
		"+1",    // explicit plus sign
		"1.e5",  // decimal point without fractional part
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, err := ParseJSON(input)
			if err == nil {
				t.Errorf("expected error for invalid number: %s", input)
			}
		})
	}
}

func TestJSONError_ErrorMessage(t *testing.T) {
	_, err := ParseJSON(`{"invalid": }`)
	if err == nil {
		t.Fatal("expected error")
	}

	jsonErr, ok := err.(*JSONError)
	if !ok {
		t.Fatalf("expected JSONError, got %T", err)
	}

	expectedMsg := "JSON parsing error at line 1, column 13: unexpected character '}'"
	if jsonErr.Error() != expectedMsg {
		t.Errorf("expected error message %q, got %q", expectedMsg, jsonErr.Error())
	}
}

// Benchmarks para evaluar el rendimiento
func BenchmarkParseJSON_SimpleObject(b *testing.B) {
	json := `{"name": "John", "age": 30, "active": true}`
	for i := 0; i < b.N; i++ {
		_, err := ParseJSON(json)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseJSON_ComplexNested(b *testing.B) {
	json := `{
		"users": [
			{"name": "Alice", "age": 25, "skills": ["go", "python", "js"]},
			{"name": "Bob", "age": 30, "skills": ["java", "c++"]},
			{"name": "Charlie", "age": 35, "skills": ["rust", "go"]}
		],
		"meta": {
			"total": 3,
			"active": true,
			"created": "2023-01-01"
		}
	}`
	for i := 0; i < b.N; i++ {
		_, err := ParseJSON(json)
		if err != nil {
			b.Fatal(err)
		}
	}
}
