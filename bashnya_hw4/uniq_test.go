package main

import (
	"bytes"
	"strings"
	"testing"
)

const sample = `I love music.
I love music.
I love music.

I love music of Kartik.
I love music of Kartik.
Thanks.
`

func run(input string, opt options) string {
	var out bytes.Buffer
	if err := process(strings.NewReader(input), &out, opt); err != nil {
		return "ERR:" + err.Error()
	}
	return out.String()
}

func TestNoFlags(t *testing.T) {
	exp := "I love music.\n\nI love music of Kartik.\nThanks.\n"
	got := run(sample, options{})
	if got != exp {
		t.Fatalf("no flags:\nexp=%q\n got=%q", exp, got)
	}
}

func TestCount(t *testing.T) {
	exp := "3 I love music.\n1 \n2 I love music of Kartik.\n1 Thanks.\n"
	got := run(sample, options{count: true})
	if got != exp {
		t.Fatalf("-c:\nexp=%q\n got=%q", exp, got)
	}
}

func TestOnlyDuplicates(t *testing.T) {
	exp := "I love music.\nI love music of Kartik.\n"
	got := run(sample, options{onlyDup: true})
	if got != exp {
		t.Fatalf("-d:\nexp=%q\n got=%q", exp, got)
	}
}

func TestOnlyUniques(t *testing.T) {
	exp := "\nThanks.\n"
	got := run(sample, options{onlyUniq: true})
	if got != exp {
		t.Fatalf("-u:\nexp=%q\n got=%q", exp, got)
	}
}

func TestIgnoreCase(t *testing.T) {
	input := "I LOVE MUSIC.\nI love music.\nI LoVe MuSiC.\n\nI love MuSIC of Kartik.\nI love music of kartik.\nThanks.\n"
	exp := "I LOVE MUSIC.\n\nI love MuSIC of Kartik.\nThanks.\n"
	got := run(input, options{ignoreCase: true})
	if got != exp {
		t.Fatalf("-i:\nexp=%q\n got=%q", exp, got)
	}
}

func TestSkipFields(t *testing.T) {
	input := "We love music.\nI love music.\nThey love music.\n\nI love music of Kartik.\nWe love music of Kartik.\nThanks.\n"
	exp := "We love music.\n\nI love music of Kartik.\nThanks.\n"
	got := run(input, options{skipFields: 1})
	if got != exp {
		t.Fatalf("-f 1:\nexp=%q\n got=%q", exp, got)
	}
}

func TestSkipChars(t *testing.T) {
	input := "I love music.\nA love music.\nC love music.\n\nI love music of Kartik.\nWe love music of Kartik.\nThanks.\n"
	exp := "I love music.\n\nI love music of Kartik.\nWe love music of Kartik.\nThanks.\n"
	got := run(input, options{skipChars: 1})
	if got != exp {
		t.Fatalf("-s 1:\nexp=%q\n got=%q", exp, got)
	}
}

func TestSkipFieldsAndChars(t *testing.T) {
	input := "X abcdef\nY abcdef\nZ abcxyz\n"
	exp := "X abcdef\nZ abcxyz\n"
	got := run(input, options{skipFields: 1, skipChars: 3})
	if got != exp {
		t.Fatalf("-f 1 -s 3:\nexp=%q\n got=%q", exp, got)
	}
}

func TestMutuallyExclusive(t *testing.T) {
	if err := options{count: true, onlyUniq: true}.validate(); err == nil {
		t.Fatalf("expected error on -c with -u, got nil")
	}
}