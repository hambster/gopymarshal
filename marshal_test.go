package gopymarshal

import (
	"testing"
)

func TestMarshal(t *testing.T) {
	var ret []byte
	var err error

	var n int32
	n = 10
	if ret, err = Marshal(n); nil != err {
		t.Fatalf("err - %s", err.Error())
	}
	t.Logf("%d = %#v", n, ret)

	var f float64
	f = 100.000
	if ret, err = Marshal(f); nil != err {
		t.Fatalf("err - %s", err)
	}
	t.Logf("%f = %#v", f, ret)

	nList := []interface{}{int32(100), int32(-100), int32(1), int32(2), int32(100)}
	if ret, err = Marshal(nList); nil != err {
		t.Fatalf("err - %s", err)
	}
	t.Logf("%#v", ret)
}
