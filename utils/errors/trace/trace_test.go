package trace

import (
	"testing"

	"go_base/utils/errors"
	"go_base/utils/testing2"
)

func TestTrace(t *testing.T) {
	tt := testing2.Wrap(t)
	tt.Nil(Trace(nil))

	var e error = errors.Err("Error")

	e2 := Trace(e)
	es := "trace/trace_test.go:16:" + e.Error()
	tt.Eq(es, e2.Error())

	e2 = Trace(e2)
	tt.Eq(es, e2.Error())

	TraceEnabled = false
	e2 = Trace(e)
	tt.Eq(e2, e)
}
