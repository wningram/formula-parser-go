package main

import "testing"

func TestGetNext(t *testing.T) {
	formula := "if([field1] < [field2], [field1], [field2])"
	result, nextNdx, err := GetNext(formula)
	// Result should be of type function
	if _, ok := result.(Function); !ok {
		t.Log("result should be a Function type.")
		t.FailNow()
	}

	// Next index should be -1 because there is no object after teh function
	if nextNdx != -1 {
		t.Logf("Next index should be -1. Instead it was %d", nextNdx)
		t.FailNow()
	}

	// There should be no error
	if err != nil {
		t.Log("An error should not have been returned.")
		t.FailNow()
	}
}

func TestGetNextWord(t *testing.T) {
	formula := "if([field1] < [field2], [field1], [field2])"
	// Execution
	curWord, nextNdx := GetNextWord(formula, 0)

	// Assertions
	if curWord != "if" {
		t.Logf("Expected curWord to be 'if'. Instead it was '%s'", curWord)
		t.Fail()
	}

	if nextNdx != 2 {
		t.Logf("Expected nextNdx to be 2. Instead it was %d", nextNdx)
		t.Fail()
	}
}

func TestParse(t *testing.T) {
	formula := "if([field1] < [field2], [field1], [field2])"

	// Execution
	block, err := Parse(formula)

	// Assertions
	if err != nil {
		t.Logf("An error was thrown by the tested function: %s", err)
		t.FailNow()
	}

	obj := block.Objects[0]
	if _, ok := obj.(IBlock); !ok {
		t.Logf("Expected result to be a Block, instead its %T", obj)
		t.Fail()
	}
}
