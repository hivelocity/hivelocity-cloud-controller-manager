package errutils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var errDummy = fmt.Errorf("dummy error")

const (
	dummyInt   = 123
	dummyFloat = 1.234
)

func TestWrap(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		err  error
		msg  string
		args []any
		want string
	}{
		{
			name: "simple case",
			err:  errDummy,
			msg:  "foo",
			args: []any{"myVar", dummyInt},
			want: "[errutils.TestWrap.func1] foo. myVar 123: dummy error",
		},
		{
			name: "more variables",
			err:  errDummy,
			msg:  "foo",
			args: []any{"varInt", dummyInt, "varFloat", dummyFloat, "varStruct", struct{ i int }{dummyInt}},
			want: "[errutils.TestWrap.func1] foo. varInt 123, varFloat 1.234, varStruct {123}: dummy error",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := Wrap(tt.err, tt.msg, tt.args...)
			require.ErrorIs(t, err, errDummy)
			require.Equal(t, tt.want, err.Error())
		})
	}
}
