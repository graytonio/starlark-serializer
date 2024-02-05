package starlarkserializer_test

import (
	"testing"

	starlarkserializer "github.com/graytonio/starlark-serializer"
	"github.com/stretchr/testify/assert"
	"go.starlark.net/starlark"
)

func getSimpleTestExpectedOutput() *starlark.Dict {
	results := starlark.NewDict(3)
	results.SetKey(starlark.String("Foo"), starlark.String("test"))
	results.SetKey(starlark.String("Bar"), starlark.MakeInt(12))
	results.SetKey(starlark.String("Baz"), starlark.Bool(false))
	return results
}

func TestMarshal(t *testing.T) {
	var tests = []struct {
		name           string
		input          any
		expectedOutput starlark.Value
		expectError    bool
	}{
		{
			name:           "string",
			input:          "hello world",
			expectedOutput: starlark.String("hello world"),
		},
		{
			name:           "int",
			input:          12,
			expectedOutput: starlark.MakeInt(12),
		},
		{
			name:           "bool",
			input:          false,
			expectedOutput: starlark.Bool(false),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := starlarkserializer.Marshal(tt.input)
			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedOutput, output)
		})
	}
}

func TestMarshalNestedStruct(t *testing.T) {
	type nestedStruct struct {
		Foo string
		Bar int
		Baz bool
	}

	type testStruct struct {
		Foo string
		Bar int
		Baz bool

		Nested nestedStruct
	}

	input := testStruct{
		Foo: "test",
		Bar: 12,
		Baz: true,
		Nested: nestedStruct{
			Foo: "test2",
			Bar: 14,
			Baz: false,
		},
	}

	output, err := starlarkserializer.Marshal(input)
	if assert.NoError(t, err) {
		assert.IsType(t, &starlark.Dict{}, output)

		dict := output.(*starlark.Dict)

		foo, _, err := dict.Get(starlark.String("foo"))
		if assert.NoError(t, err) {
			assert.Equal(t, starlark.String("test"), foo)
		}

		bar, _, err := dict.Get(starlark.String("bar"))
		if assert.NoError(t, err) {
			assert.Equal(t, starlark.MakeInt(12), bar)
		}

		baz, _, err := dict.Get(starlark.String("baz"))
		if assert.NoError(t, err) {
			assert.Equal(t, starlark.Bool(true), baz)
		}

		nested, _, err := dict.Get(starlark.String("nested"))
		if assert.NoError(t, err) {
			assert.IsType(t, &starlark.Dict{}, nested)
		}

		dict = nested.(*starlark.Dict)
		foo, _, err = dict.Get(starlark.String("foo"))
		if assert.NoError(t, err) {
			assert.Equal(t, starlark.String("test2"), foo)
		}

		bar, _, err = dict.Get(starlark.String("bar"))
		if assert.NoError(t, err) {
			assert.Equal(t, starlark.MakeInt(14), bar)
		}

		baz, _, err = dict.Get(starlark.String("baz"))
		if assert.NoError(t, err) {
			assert.Equal(t, starlark.Bool(false), baz)
		}
	}
}

func TestMarshalSimpleStruct(t *testing.T) {
	type testStruct struct {
		Foo string
		Bar int
		Baz bool
	}

	input := testStruct{
		Foo: "test",
		Bar: 12,
		Baz: true,
	}

	output, err := starlarkserializer.Marshal(input)
	if assert.NoError(t, err) {
		assert.IsType(t, &starlark.Dict{}, output)

		dict := output.(*starlark.Dict)

		foo, _, err := dict.Get(starlark.String("foo"))
		if assert.NoError(t, err) {
			assert.Equal(t, starlark.String("test"), foo)
		}

		bar, _, err := dict.Get(starlark.String("bar"))
		if assert.NoError(t, err) {
			assert.Equal(t, starlark.MakeInt(12), bar)
		}

		baz, _, err := dict.Get(starlark.String("baz"))
		if assert.NoError(t, err) {
			assert.Equal(t, starlark.Bool(true), baz)
		}

	}
}


func TestMarshalSimpleStructPointer(t *testing.T) {
	type testStruct struct {
		Foo string
		Bar int
		Baz bool
		unexport string
	}

	input := testStruct{
		Foo: "test",
		Bar: 12,
		Baz: true,
		unexport: "test",
	}

	output, err := starlarkserializer.Marshal(&input)
	if assert.NoError(t, err) {
		assert.IsType(t, &starlark.Dict{}, output)

		dict := output.(*starlark.Dict)

		foo, _, err := dict.Get(starlark.String("foo"))
		if assert.NoError(t, err) {
			assert.Equal(t, starlark.String("test"), foo)
		}

		bar, _, err := dict.Get(starlark.String("bar"))
		if assert.NoError(t, err) {
			assert.Equal(t, starlark.MakeInt(12), bar)
		}

		baz, _, err := dict.Get(starlark.String("baz"))
		if assert.NoError(t, err) {
			assert.Equal(t, starlark.Bool(true), baz)
		}

	}
}
