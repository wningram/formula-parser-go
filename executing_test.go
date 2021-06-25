package main

import (
	"testing"
)

func getTestCsv() CsvFile {
	var (
		fileName string     = "TestData"
		headers  []string   = []string{"Field1", "Field2", "Field3"}
		data     [][]string = [][]string{{"1", "2", "3"}}
	)
	return CsvFile{
		FileName: fileName,
		Headers:  headers,
		DataRows: data,
	}
}

func TestHandleOperator(t *testing.T) {
	file := getTestCsv()
	testExecutor := Executor{file}
	op := Operator{Operator: "<", OperatorType: "cond_op"}
	lOp := Field{FieldName: "Field1"}
	rOp := Field{FieldName: "Field2"}

	// Execution
	result, err := testExecutor.handleOperator(op, lOp, rOp)
	t.Log(result)
	// Assertions
	if err != nil {
		t.Logf("Tested function through an error: %s", err)
		t.FailNow()
	}

	if result[0] != true {
		t.Logf("Result was supposed to be true, instead it was %v", result[0])
	}
}

func TestExecuteOperators(t *testing.T) {
	file := getTestCsv()
	executor := Executor{Data: file}
	block, _ := Parse("if([Field1] < [Field2], [Field1], [Field2])")

	// Execution
	block, err := executor.executeOperators(block)

	// Assertions
	if err != nil {
		t.Logf("Tested function threw an error: %s", err)
		t.FailNow()
	}

	if _, ok := block.Objects[0].(IValueSet); !ok {
		t.Logf("block[0] is not a ValueSet, instead it is %t", block.Objects[0])
		t.Fail()
	}

	if len(block.Objects) != 3 {
		t.Logf("Block does not contain exactly 3 objects, instead it contains %d", len(block.Objects))
	}

}

func TestRemoveObject(t *testing.T) {
	objs := []Object{
		Field{"Fiel1"},
		Field{"Fiel2"},
		Field{"Fiel3"},
	}

	// Execution
	result := removeObject(objs, 1)

	// Assertions
	if len(result) != 2 {
		t.Log("Result length is not 2")
		t.Fail()
	}

	if result[0].(Field).FieldName != "Fiel1" || result[1].(Field).FieldName != "Fiel3" {
		t.Logf("Result has unexpected value: %v", result)
		t.Fail()
	}
}

func TestEvaluateIfFunc(t *testing.T) {
	values := make([]interface{}, 0)
	values = append(values, true, false)
	block := []Object{
		ValueSet{values},
		Field{"Field1"},
		Field{"Field2"},
	}
	executor := Executor{getTestCsv()}

	// Execution
	results, err := executor.evaluateIfFunc(block)

	if err != nil {
		t.Logf("Tested function raised an error: %s", err)
		t.FailNow()
	}

	if len(results) != 2 {
		t.Logf("Expected results length to be 2. Instead it was %d", len(results))
		t.Fail()
	}

	f1, f1Ok := results[0].(Field)
	f2, f2Ok := results[1].(Field)
	if !f1Ok || f1.FieldName != "Field1" {
		t.Logf("First item returned is invalid: %v", results[0])
		t.Fail()
	}

	if !f2Ok || f2.FieldName != "Field2" {
		t.Logf("Second item returned is invalid: %v", results[1])
		t.Fail()
	}
}
