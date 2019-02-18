package httpheader

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestHeader_types(t *testing.T) {
	str := "string"
	strPtr := &str
	timeVal := time.Date(2000, 1, 1, 12, 34, 56, 0, time.UTC)

	tests := []struct {
		in   interface{}
		want http.Header
	}{
		{
			// basic primitives
			struct {
				A string
				B int
				C uint
				D float32
				E bool
			}{},
			http.Header{
				"A": []string{""},
				"B": []string{"0"},
				"C": []string{"0"},
				"D": []string{"0"},
				"E": []string{"false"},
			},
		},
		{
			// pointers
			struct {
				A *string
				B *int
				C **string
				D *time.Time
			}{
				A: strPtr,
				C: &strPtr,
				D: &timeVal,
			},
			http.Header{
				"A": []string{str},
				"B": []string{""},
				"C": []string{str},
				"D": []string{"Sat, 01 Jan 2000 12:34:56 GMT"},
			},
		},
		{
			// slices and arrays
			struct {
				A []string
				B []*string
				C [2]string
				D []bool `header:",int"`
			}{
				A: []string{"a", "b"},
				B: []*string{&str, &str},
				C: [2]string{"a", "b"},
				D: []bool{true, false},
			},
			http.Header{
				"A": []string{"a", "b"},
				"B": {"string", "string"},
				"C": []string{"a", "b"},
				"D": {"1", "0"},
			},
		},
		{
			// other types
			struct {
				A time.Time
				B time.Time `header:",unix"`
				C bool      `header:",int"`
				D bool      `header:",int"`
				E http.Header
			}{
				A: time.Date(2000, 1, 1, 12, 34, 56, 0, time.UTC),
				B: time.Date(2000, 1, 1, 12, 34, 56, 0, time.UTC),
				C: true,
				D: false,
				E: http.Header{
					"F": []string{"f1"},
					"G": []string{"gg"},
				},
			},
			http.Header{
				"A": []string{"Sat, 01 Jan 2000 12:34:56 GMT"},
				"B": []string{"946730096"},
				"C": []string{"1"},
				"D": []string{"0"},
				"F": []string{"f1"},
				"G": []string{"gg"},
			},
		},
		{
			nil,
			http.Header{},
		},
		{
			&struct {
				A string
			}{"test"},
			http.Header{
				"A": []string{"test"},
			},
		},
	}

	for i, tt := range tests {
		v, err := Header(tt.in)
		if err != nil {
			t.Errorf("%d. Header(%q) returned error: %v", i, tt.in, err)
		}

		if !reflect.DeepEqual(tt.want, v) {
			t.Errorf("%d. Header(%q) returned %#v, want %#v", i, tt.in, v, tt.want)
		}
	}
}

func TestHeader_omitEmpty(t *testing.T) {
	str := ""
	s := struct {
		a string
		A string
		B string    `header:",omitempty"`
		C string    `header:"-"`
		D string    `header:"omitempty"` // actually named omitempty, not an option
		E *string   `header:",omitempty"`
		F bool      `header:",omitempty"`
		G int       `header:",omitempty"`
		H uint      `header:",omitempty"`
		I float32   `header:",omitempty"`
		J time.Time `header:",omitempty"`
		K struct{}  `header:",omitempty"`
	}{E: &str}

	v, err := Header(s)
	if err != nil {
		t.Errorf("Header(%#v) returned error: %v", s, err)
	}

	want := http.Header{
		"A":         []string{""},
		"Omitempty": []string{""},
		"E":         []string{""}, // E is included because the pointer is not empty, even though the string being pointed to is
	}
	if !reflect.DeepEqual(want, v) {
		t.Errorf("Header(%#v) returned %v, want %v", s, v, want)
	}
}

type A struct {
	B
}

type B struct {
	C string
}

type D struct {
	B
	C string
}

type e struct {
	B
	C string
}

type F struct {
	e
}

func TestHeader_embeddedStructs(t *testing.T) {
	tests := []struct {
		in   interface{}
		want http.Header
	}{
		{
			A{B{C: "foo"}},
			http.Header{"C": []string{"foo"}},
		},
		{
			D{B: B{C: "bar"}, C: "foo"},
			http.Header{"C": []string{"foo", "bar"}},
		},
		{
			F{e{B: B{C: "bar"}, C: "foo"}}, // With unexported embed
			http.Header{"C": []string{"foo", "bar"}},
		},
	}

	for i, tt := range tests {
		v, err := Header(tt.in)
		if err != nil {
			t.Errorf("%d. Header(%q) returned error: %v", i, tt.in, err)
		}

		if !reflect.DeepEqual(tt.want, v) {
			t.Errorf("%d. Header(%q) returned %v, want %v", i, tt.in, v, tt.want)
		}
	}
}

func TestHeader_invalidInput(t *testing.T) {
	_, err := Header("")
	if err == nil {
		t.Errorf("expected Header() to return an error on invalid input")
	}
}

type EncodedArgs []string

func (m EncodedArgs) EncodeHeader(key string, v *http.Header) error {
	for i, arg := range m {
		v.Set(fmt.Sprintf("%s.%d", key, i), arg)
	}
	return nil
}

func TestHeader_Marshaler(t *testing.T) {
	s := struct {
		Args EncodedArgs `header:"arg"`
	}{[]string{"a", "b", "c"}}
	v, err := Header(s)
	if err != nil {
		t.Errorf("Header(%q) returned error: %v", s, err)
	}

	want := http.Header{
		"Arg.0": []string{"a"},
		"Arg.1": []string{"b"},
		"Arg.2": []string{"c"},
	}
	if !reflect.DeepEqual(want, v) {
		t.Errorf("Header(%q) returned %v, want %v", s, v, want)
	}
}

func TestHeader_MarshalerWithNilPointer(t *testing.T) {
	s := struct {
		Args *EncodedArgs `header:"arg"`
	}{}
	v, err := Header(s)
	if err != nil {
		t.Errorf("Header(%q) returned error: %v", s, err)
	}

	want := http.Header{}
	if !reflect.DeepEqual(want, v) {
		t.Errorf("Header(%q) returned %v, want %v", s, v, want)
	}
}

func TestTagParsing(t *testing.T) {
	name, opts := parseTag("field,foobar,foo")
	if name != "field" {
		t.Fatalf("name = %q, want field", name)
	}
	for _, tt := range []struct {
		opt  string
		want bool
	}{
		{"foobar", true},
		{"foo", true},
		{"bar", false},
		{"field", false},
	} {
		if opts.Contains(tt.opt) != tt.want {
			t.Errorf("Contains(%q) = %v", tt.opt, !tt.want)
		}
	}
}
