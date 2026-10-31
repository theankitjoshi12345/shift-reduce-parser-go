package main

import "testing"

func TestParserAcceptsExamples(t *testing.T) {
	for _, input := range [][]string{
		{"id", "+", "id", "*", "id"},
		{"(", "id", ")"},
	} {
		parser := newParser(input)
		for !parser.accepted && !parser.ungrammatical {
			parser.parse1Step()
		}
		if !parser.accepted || parser.ungrammatical {
			t.Fatalf("expected %v to be accepted", input)
		}
		if len(parser.trees.items) != 1 {
			t.Fatalf("expected one parse tree for %v, got %d", input, len(parser.trees.items))
		}
	}
}

func TestParserRejectsInvalidInput(t *testing.T) {
	parser := newParser([]string{"id", "+", "*", "id"})
	for !parser.accepted && !parser.ungrammatical {
		parser.parse1Step()
	}
	if parser.accepted || !parser.ungrammatical {
		t.Fatal("expected invalid input to be rejected")
	}
}

func TestTokenize(t *testing.T) {
	got := tokenize("id + id * id")
	want := []string{"id", "+", "id", "*", "id"}
	if len(got) != len(want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %#v, want %#v", got, want)
		}
	}
}
