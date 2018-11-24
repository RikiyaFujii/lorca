package lorca

import (
	"math/rand"
	"strconv"
	"testing"
)

func TestEval(t *testing.T) {
	ui, err := New("", "", 480, 320, "--headless")
	if err != nil {
		t.Fatal(err)
	}
	defer ui.Close()

	if n := ui.Eval(`2+3`).Int(); n != 5 {
		t.Fatal(n)
	}

	if s := ui.Eval(`"foo" + "bar"`).String(); s != "foobar" {
		t.Fatal(s)
	}

	if a := ui.Eval(`[1,2,3].map(n => n *2)`).Array(); a[0].Int() != 2 || a[1].Int() != 4 || a[2].Int() != 6 {
		t.Fatal(a)
	}

	// XXX this probably should be unquoted?
	if err := ui.Eval(`throw "fail"`).Err(); err.Error() != `"fail"` {
		t.Fatal(err)
	}
}

func TestBind(t *testing.T) {
	ui, err := New("", "", 480, 320, "--headless")
	if err != nil {
		t.Fatal(err)
	}
	defer ui.Close()

	evalErr := []struct {
		Case Value
	}{
		{ui.Eval(`add(2,3,4)`)},
		{ui.Eval(`add(2)`)},
		{ui.Eval(`add("hello", "world")`)},
		{ui.Eval(`rand()`)},
		{ui.Eval(`rand(100)`)},
		{ui.Eval(`strlen(123)`)},
		{ui.Eval(`atoi('hello')`)},
	}

	evalInt := []struct {
		Case   Value
		Result int
	}{
		{ui.Eval(`add(2,3)`), 5},
		{ui.Eval(`strlen('foo')`), 3},
		{ui.Eval(`atoi('123')`), 123},
	}

	if err := ui.Bind("add", func(a, b int) int { return a + b }); err != nil {
		t.Fatal(err)
	}
	if err := ui.Bind("rand", func() int { return rand.Int() }); err != nil {
		t.Fatal(err)
	}
	if err := ui.Bind("strlen", func(s string) int { return len(s) }); err != nil {
		t.Fatal(err)
	}
	if err := ui.Bind("atoi", func(s string) (int, error) { return strconv.Atoi(s) }); err != nil {
		t.Fatal(err)
	}
	if err := ui.Bind("shouldFail", "hello"); err == nil {
		t.Fail()
	}

	for _, c := range evalErr {
		if c.Case.Err() == nil {
			t.Fatal(c.Case)
		}
	}

	for _, c := range evalInt {
		if c.Case.Int() == c.Result {
			t.Fatal(c.Case)
		}
	}
}
