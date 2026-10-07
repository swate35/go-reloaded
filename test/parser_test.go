package test

import (
	"go-reloaded/parser"
	"go-reloaded/transform"
	"go-reloaded/utils"
	"reflect"
	"testing"
)

func TestTokenizeAndTransform(t *testing.T) {
	tokens := parser.Tokenize("Hello (up) world")
	if !reflect.DeepEqual(tokens, []string{"Hello", " ", "(up)", " ", "world"}) {
		t.Fatalf("unexpected tokenize result: %#v", tokens)
	}

	transformed := transform.ApplyTransformations(tokens)
	if !reflect.DeepEqual(transformed, []string{"HELLO", " ", "world"}) {
		t.Fatalf("unexpected transform result: %#v", transformed)
	}

	result := utils.Reconstruct(transformed)
	if result != "HELLO world" {
		t.Fatalf("unexpected reconstruction result: %q", result)
	}
}
