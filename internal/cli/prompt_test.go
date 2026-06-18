package cli

import (
	"strings"
	"testing"
)

func TestAskIntDefault(t *testing.T) {
	p := New(strings.NewReader("\n"), &strings.Builder{})
	n, err := p.AskInt("Port", 8080)
	if err != nil || n != 8080 {
		t.Errorf("boş girdi varsayılanı döndürmeli, bulundu %d err %v", n, err)
	}
}

func TestAskIntValue(t *testing.T) {
	p := New(strings.NewReader("9090\n"), &strings.Builder{})
	n, _ := p.AskInt("Port", 8080)
	if n != 9090 {
		t.Errorf("9090 bekleniyordu, bulundu %d", n)
	}
}

func TestAskYesNo(t *testing.T) {
	cases := map[string]bool{"e\n": true, "h\n": false, "\n": true, "evet\n": true}
	for in, want := range cases {
		p := New(strings.NewReader(in), &strings.Builder{})
		got, _ := p.AskYesNo("Devam?", true)
		if got != want {
			t.Errorf("girdi %q için %v bekleniyordu, bulundu %v", in, want, got)
		}
	}
}
