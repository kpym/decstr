package decstr

import (
	"fmt"
	"testing"
)

func TestDecimalFormatString(t *testing.T) {
	tests := []struct {
		df   DecimalFormat
		want string
	}{
		{DecimalFormat{Point: '.', Group: NoSeparator, Standard: true}, "{`.`, `<none>`, standard}"},
		{DecimalFormat{Point: '.', Group: ' ', Standard: true}, "{`.`, ` `, standard}"},
		{DecimalFormat{Point: ',', Group: '\'', Standard: false}, "{`,`, `'`, non-standard}"},
		{DecimalFormat{Point: '·', Group: NoSeparator, Standard: false}, "{`·`, `<none>`, non-standard}"},
	}

	for _, test := range tests {
		got := test.df.String()
		if got != test.want {
			t.Errorf("(%v).String() = %q, want %q", test.df, got, test.want)
		}
	}
}

func TestGetSign(t *testing.T) {
	testStrings := []struct {
		decimal string
		sign    string
		abs     string
	}{
		{"", "", ""},
		{"  ", "", ""},
		{"0", "", "0"},
		{" 0", "", "0"},
		{"0 ", "", "0"},
		{"+1", "", "1"},
		{"+ 123", "", "123"},
		{"-1", "-", "1"},
		{"  -   123  ", "-", "123"},
	}

	testBytes := []struct {
		decimal []byte
		sign    []byte
		abs     []byte
	}{
		{[]byte(""), []byte(""), []byte("")},
		{[]byte("  "), []byte(""), []byte("")},
		{[]byte("0"), []byte(""), []byte("0")},
		{[]byte(" 0"), []byte(""), []byte("0")},
		{[]byte("0 "), []byte(""), []byte("0")},
		{[]byte("+1"), []byte(""), []byte("1")},
		{[]byte("+ 123"), []byte(""), []byte("123")},
		{[]byte("-1"), []byte("-"), []byte("1")},
		{[]byte("  -   123  "), []byte("-"), []byte("123")},
	}

	for _, test := range testStrings {
		sign, abs := getSign(test.decimal)
		if sign != test.sign || abs != test.abs {
			t.Errorf("GetSign(%q) = (%q, %q), want (%q, %q)", test.decimal, sign, abs, test.sign, test.abs)
		}
	}

	for _, test := range testBytes {
		sign, abs := getSign(test.decimal)
		if string(sign) != string(test.sign) || string(abs) != string(test.abs) {
			t.Errorf("GetSign(%q) = (%q, %q), want (%q, %q)", test.decimal, sign, abs, test.sign, test.abs)
		}
	}
}

func TestDetectFormat(t *testing.T) {
	tests := []struct {
		decimal string
		df      DecimalFormat
		ok      bool
	}{
		{"", DecimalFormat{}, false},
		{"  ", DecimalFormat{}, false},
		{"123", DecimalFormat{Point: NoSeparator, Group: NoSeparator, Standard: true}, true},
		{"1 234", DecimalFormat{Point: NoSeparator, Group: ' ', Standard: true}, true},
		{"1,234", DecimalFormat{}, false}, // ambiguous
		{"1.234", DecimalFormat{}, false}, // ambiguous
		{"1'234", DecimalFormat{}, false}, // ambiguous
		{"1·234", DecimalFormat{Point: '·', Group: NoSeparator, Standard: true}, true},
		{"1 234.56", DecimalFormat{Point: '.', Group: ' ', Standard: true}, true},
		{"1,234.56", DecimalFormat{Point: '.', Group: ',', Standard: true}, true},
		{"1'234.56", DecimalFormat{Point: '.', Group: '\'', Standard: true}, true},
		{"1·234.56", DecimalFormat{}, false},
		{"1 234,56", DecimalFormat{Point: ',', Group: ' ', Standard: true}, true},
		{"1.234,56", DecimalFormat{Point: ',', Group: '.', Standard: true}, true},
		{"1'234,56", DecimalFormat{Point: ',', Group: '\'', Standard: true}, true},
		{"1·234,56", DecimalFormat{}, false},
		{"1.234'56", DecimalFormat{Point: '\'', Group: '.', Standard: true}, true},
		{"1·234'56", DecimalFormat{}, false},
		{"1,234'56", DecimalFormat{}, false},
		{"1 234'56", DecimalFormat{}, false},
		{"1,234·56", DecimalFormat{Point: '·', Group: ',', Standard: true}, true},
		{"1 234·56", DecimalFormat{}, false},
		{"1'234·56", DecimalFormat{}, false},
		{"1.234·56", DecimalFormat{}, false},
		{"1'234'56", DecimalFormat{}, false},
		{"1'234'567", DecimalFormat{Point: NoSeparator, Group: '\'', Standard: true}, true},
		{"1'34'567", DecimalFormat{Point: NoSeparator, Group: '\'', Standard: false}, true},
		{"1 234 56", DecimalFormat{}, false},
		{"1 234 567", DecimalFormat{Point: NoSeparator, Group: ' ', Standard: true}, true},
		{"1 34 567", DecimalFormat{Point: NoSeparator, Group: ' ', Standard: false}, true},
		{"1 234 567.8", DecimalFormat{Point: '.', Group: ' ', Standard: true}, true},
		{"1 34 567.8", DecimalFormat{Point: '.', Group: ' ', Standard: false}, true},
		{".12", DecimalFormat{Point: '.', Group: NoSeparator, Standard: true}, true},
		{"12.", DecimalFormat{Point: '.', Group: NoSeparator, Standard: true}, true},
		{"12.345 678", DecimalFormat{}, false},
		{"12¸345", DecimalFormat{}, false},
		{"1234 567,8", DecimalFormat{}, false},
		{"1'234 567,8", DecimalFormat{}, false},
		{"1'2345'678", DecimalFormat{}, false},
		{"1'23'678'901", DecimalFormat{}, false},
	}

	for _, test := range tests {
		df, ok := DetectFormat(test.decimal)
		if df != test.df || ok != test.ok {
			t.Errorf("DetectFormat(%q) = (%v, %v), want (%v, %v)", test.decimal, df, ok, test.df, test.ok)
		}
	}
}

func ExampleDetectFormat() {
	df, ok := DetectFormat("1 234,56")
	if !ok {
		fmt.Println("not a decimal")
	}
	fmt.Println(df)
	// Output: {`,`, ` `, standard}
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		decimal string
		want    string
	}{
		{"123", "123"},
		{"1 234", "1234"},
		{"1·234", "1.234"},
		{"1 234.56", "1234.56"},
		{"1,234.56", "1234.56"},
		{"1'234.56", "1234.56"},
		{"1 234,56", "1234.56"},
		{"1.234,56", "1234.56"},
		{"1'234,56", "1234.56"},
		{"1.234'56", "1234.56"},
		{"1,234·56", "1234.56"},
		{"1'234'567", "1234567"},
		{"1'34'567", "134567"},
		{"1 234 567", "1234567"},
		{"1 34 567", "134567"},
		{"1 234 567.8", "1234567.8"},
		{"1 34 567.8", "134567.8"},
		{".12", "0.12"},
		{"12.", "12"},
		{"012.", "12"},
		{"012.3", "12.3"},
		{"12.0", "12"},
		{"12.30", "12.3"},
		{"1,234", "1,234"},           // ambiguous
		{"1.234", "1.234"},           // ambiguous
		{"1'234", "1'234"},           // ambiguous
		{"", ""},                     // not a decimal
		{"  ", "  "},                 // not a decimal
		{" test ", " test "},         // not a decimal
		{",", ","},                   // not a decimal
		{"-,", "-,"},                 // not a decimal
		{".", "."},                   // not a decimal
		{"-.", "-."},                 // not a decimal
		{"+.", "+."},                 // not a decimal
		{" - .", " - ."},             // not a decimal
		{"1·234.56", "1·234.56"},     // not a decimal
		{"1·234,56", "1·234,56"},     // not a decimal
		{"1·234'56", "1·234'56"},     // not a decimal
		{"1,234'56", "1,234'56"},     // not a decimal
		{"1 234'56", "1 234'56"},     // not a decimal
		{"1 234·56", "1 234·56"},     // not a decimal
		{"1'234·56", "1'234·56"},     // not a decimal
		{"1.234·56", "1.234·56"},     // not a decimal
		{"1'234'56", "1'234'56"},     // not a decimal
		{"1 234 56", "1 234 56"},     // not a decimal
		{"12.345 678", "12.345 678"}, // not a decimal
	}

	for _, test := range tests {
		got := Normalize(test.decimal)
		if got != test.want {
			t.Errorf("Normalize(%q) = %q, want %q", test.decimal, got, test.want)
		}
		_, ok := DetectFormat(test.decimal)
		// if it was a decimal but the result is not normalized
		if ok && !IsNormalized(got) {
			t.Errorf("Normalize(%q) = %q is not normalized", test.decimal, got)
		}
	}
}

// BenchmarkNormalize compare Normalize and AutoNormalize functions
func BenchmarkNormalizeString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Normalize("1 234,50")
	}
}

func BenchmarkNormalizeSlice(b *testing.B) {
	buf := []byte("1 234,50")
	for i := 0; i < b.N; i++ {
		Normalize(string(buf))
	}
}

func ExampleNormalize() {
	fmt.Println(Normalize(" - 1 234,50 "))
	fmt.Println(Normalize("12 345."))
	// Output:
	// -1234.5
	// 12345
}

func TestNormalizeCheck(t *testing.T) {
	data := []struct {
		decimal string
		want    string
		ok      bool
	}{
		{"123", "123", true},
		{"1 234", "1234", true},
		{"1·234", "1.234", true},
		{"1 234.56", "1234.56", true},
		{"1,234.56", "1234.56", true},
		{"1'234.56", "1234.56", true},
		{"1 234,56", "1234.56", true},
		{"1.234,56", "1234.56", true},
		{"1'234,56", "1234.56", true},
		{"1.234'56", "1234.56", true},
		{"1,234·56", "1234.56", true},
		{"1'234'567", "1234567", true},
		{"1'34'567", "134567", true},
		{"1 234 567", "1234567", true},
		{"1 34 567", "134567", true},
		{"1 234 567.8", "1234567.8", true},
		{"1 34 567.8", "134567.8", true},
		{".12", "0.12", true},
		{"12.", "12", true},
		{"012.", "12", true},
		{"012.3", "12.3", true},
		{"12.0", "12", true},
		{"12.30", "12.3", true},
		{"1,234", "1,234", false},           // ambiguous
		{"1.234", "1.234", false},           // ambiguous
		{"1'234", "1'234", false},           // ambiguous
		{"", "", false},                     // not a decimal
		{"  ", "  ", false},                 // not a decimal
		{" test ", " test ", false},         // not a decimal
		{",", ",", false},                   // not a decimal
		{"-,", "-,", false},                 // not a decimal
		{".", ".", false},                   // not a decimal
		{"-.", "-.", false},                 // not a decimal
		{"+.", "+.", false},                 // not a decimal
		{" - .", " - .", false},             // not a decimal
		{"1·234.56", "1·234.56", false},     // not a decimal
		{"1·234,56", "1·234,56", false},     // not a decimal
		{"1·234'56", "1·234'56", false},     // not a decimal
		{"1,234'56", "1,234'56", false},     // not a decimal
		{"1 234'56", "1 234'56", false},     // not a decimal
		{"1 234·56", "1 234·56", false},     // not a decimal
		{"1'234·56", "1'234·56", false},     // not a decimal
		{"1.234·56", "1.234·56", false},     // not a decimal
		{"1'234'56", "1'234'56", false},     // not a decimal
		{"1 234 56", "1 234 56", false},     // not a decimal
		{"12.345 678", "12.345 678", false}, // not a decimal
	}

	for _, test := range data {
		got, ok := NormalizeCheck(test.decimal)
		if got != test.want || ok != test.ok {
			t.Errorf("NormalizeCheck(%q) = (%q, %v), want (%q, %v)", test.decimal, got, ok, test.want, test.ok)
		}
	}
}

func TestIsNormalized(t *testing.T) {
	data := []struct {
		decimal string
		want    bool
	}{
		{"0", true},
		{"1230", true},
		{"-123", true},
		{"0.1", true},
		{"-0.1", true},
		{"123.45", true},
		{"-123.45", true},
		{"-0", false},       // not standard 0
		{"", false},         // not a decimal
		{"a", false},        // not a decimal
		{"0123", false},     // starts with 0
		{"-0123", false},    // starts with 0
		{".", false},        // starts with '.'
		{".12", false},      // starts with '.'
		{"0.", false},       // trailing '.'
		{"-0.", false},      // trailing '.'
		{"123.", false},     // trailing '.'
		{"-123.", false},    // trailing '.'
		{"0.0", false},      // trailing '0'
		{"0.10", false},     // trailing '0'
		{"1 234", false},    // hase group separator
		{"1·234", false},    // hase '·' character
		{"1 234.56", false}, // hase space
		{" 1234.56", false}, // hase space
		{"1234.56 ", false}, // hase space
	}

	for _, test := range data {
		got := IsNormalized(test.decimal)
		if got != test.want {
			t.Errorf("IsNormalized(%q) = %v, want %v", test.decimal, got, test.want)
		}
	}
}

func ExampleIsNormalized() {
	fmt.Println(IsNormalized("-123.45"))
	fmt.Println(IsNormalized("1 234.5"))
	// Output:
	// true
	// false
}

func TestConvert(t *testing.T) {
	data := []struct {
		df      DecimalFormat
		decimal string
		want    string
		ok      bool
	}{
		{DecimalFormat{Point: '.', Group: NoSeparator, Standard: true}, "123", "123", true},
		{DecimalFormat{Point: '.', Group: ' ', Standard: true}, "+ 1234", "1 234", true},
		{DecimalFormat{Point: '.', Group: ' ', Standard: true}, "123456789", "123 456 789", true},
		{DecimalFormat{Point: '.', Group: ' ', Standard: false}, "123456789", "12 34 56 789", true},
		{DecimalFormat{Point: '.', Group: ' ', Standard: false}, "-23456789", "-2 34 56 789", true},
		{DecimalFormat{Point: '.', Group: ' ', Standard: true}, "123456789.123", "123 456 789.123", true},
		{DecimalFormat{Point: '.', Group: ' ', Standard: false}, "123456789.123", "12 34 56 789.123", true},
		{DecimalFormat{Point: '.', Group: ' ', Standard: false}, "- 23456789.123", "-2 34 56 789.123", true},
		{DecimalFormat{Point: '.', Group: ' ', Standard: true}, "+123.456.789,123", "123 456 789.123", true},
		{DecimalFormat{Point: '.', Group: ' ', Standard: false}, "12 34 56 789,123", "12 34 56 789.123", true},
		{DecimalFormat{Point: '.', Group: ' ', Standard: false}, " - 23 456 789,123", "-2 34 56 789.123", true},
		{DecimalFormat{Point: '·', Group: ',', Standard: false}, " - 23 456 789,123", "-2,34,56,789·123", true},
		{DecimalFormat{Point: '·', Group: ',', Standard: false}, "", "0", false},
		{DecimalFormat{Point: '·', Group: ',', Standard: false}, " ", "0", false},
		{DecimalFormat{Point: '·', Group: ',', Standard: false}, " . ", "0", false},
		{DecimalFormat{Point: '·', Group: ',', Standard: false}, " -. ", "0", false},
		{DecimalFormat{Point: '·', Group: ',', Standard: false}, " - 123 45 6789,123", "0", false},
	}

	for _, test := range data {
		got, ok := test.df.Convert(test.decimal)
		if got != test.want || ok != test.ok {
			t.Errorf("(%v).Convert(%q) = (%q, %v), want (%q, %v)", test.df, test.decimal, got, ok, test.want, test.ok)
		}
	}
}

func ExampleDecimalFormat_Convert() {
	df := DecimalFormat{Point: ',', Group: ' ', Standard: true}
	new, ok := df.Convert("123456789.123")
	if !ok {
		fmt.Println("not a decimal")
	}
	fmt.Println(new)
	// Output: 123 456 789,123
}

// Example demonstrates general usage of the decstr package, including
// normalization, format detection, and conversion of decimal strings.
func Example() {
	decimal := "1'234'567,89"

	// Normalize example
	normalized := Normalize(decimal)
	fmt.Println("Normalized:", normalized)

	// Detect format example
	format, ok := DetectFormat(decimal)
	fmt.Println("Detected format:", format, "ok:", ok)
	// Convert example
	df := DecimalFormat{Point: '.', Group: ' ', Standard: false}
	converted, ok := df.Convert(decimal)
	fmt.Println("Converted:", converted, "ok:", ok)
	// Output:
	// Normalized: 1234567.89
	// Detected format: {`,`, `'`, standard} ok: true
	// Converted: 12 34 567.89 ok: true
}
