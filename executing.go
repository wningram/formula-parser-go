package main

import (
	"fmt"
	"strings"

	"github.com/Knetic/govaluate"
)

func removeObject(slice []Object, s int) []Object {
	return append(slice[:s], slice[s+1:]...)
}

type ExecutionError struct {
	Msg string
}

func (e *ExecutionError) Error() string {
	return e.Msg
}

type Executor struct {
	Data CsvFile
}

func (e Executor) handleOperator(op Operator, lOp Object, rOp Object) ([]interface{}, error) {
	var results []interface{}

	_, lOpOk := lOp.(IField)
	_, rOpOk := rOp.(IField)
	if lOpOk && rOpOk {
		for i := 0; i < len(e.Data.DataRows); i++ {
			var evalStr string
			lOpVal := e.Data.GetColumn(lOp.(Field).FieldName)[i]
			rOpVal := e.Data.GetColumn(rOp.(Field).FieldName)[i]
			if op.OperatorType == "text_op" {
				evalStr = fmt.Sprintf("\"%s\"%s\"%s\"", lOpVal, op.Operator, rOpVal)
			} else {
				evalStr = fmt.Sprintf("%s%s%s", lOpVal, op.Operator, rOpVal)
			}
			// Translate operators
			evalStr = strings.ReplaceAll(evalStr, "=", "==")
			evalStr = strings.ReplaceAll(evalStr, "<>", "!=")
			expr, err := govaluate.NewEvaluableExpression(evalStr)
			if err != nil {
				return nil, err
			}
			result, err := expr.Evaluate(nil)
			results = append(results, result)
		}
	}

	return results, nil
}

func (e Executor) executeOperators(block Block) (Block, error) {
	for i := 0; i < len(block.Objects); i++ {
		curObj := block.Objects[i]
		if _, ok := curObj.(IOperator); ok {
			newValue, err := e.handleOperator(
				curObj.(Operator),
				block.Objects[i-1],
				block.Objects[i+1],
			)

			if err != nil {
				return Block{}, err
			}

			// Replace operator with values of evaluation
			block.Objects[i] = ValueSet{newValue}
			// Delete operands from slice
			block.Objects = removeObject(block.Objects, i-1)
			block.Objects = removeObject(block.Objects, i)
			return e.executeOperators(block)
		}

		if _, ok := curObj.(IBlock); ok {
			blockObject := curObj.(Block)
			return e.executeOperators(blockObject)
		}
	}
	return block, nil
}

func (e Executor) evaluateIfFunc(block []Object) ([]interface{}, error) {
	var results []interface{}
	if len(block) != 3 {
		return nil, &ExecutionError{"If func not in correct format."}
	} else {
		cond := block[0]
		t_result := block[1]
		f_result := block[2]
		if _, ok := cond.(IValueSet); !ok {
			return nil, &ExecutionError{"If func condition is not a value."}
		}

		for _, val := range cond.(ValueSet).Values {
			if val == true {
				results = append(results, t_result)
			} else {
				results = append(results, f_result)
			}
		}
	}

	return results, nil
}

func (e Executor) Execute(block Block) ([]interface{}, error) {
	var endResults []interface{}
	curFunc := block.InnerFunctionName
	for i := len(block.Objects); i > -1; i-- {
		switch curFunc {
		case "if":
			// Clear end results
			endResults = make([]interface{}, 0)
			// Evaluate operators
			evaluatedBlock, err := e.executeOperators(block)
			if err != nil {
				return nil, err
			}

			// Evaluate if functions
			results, err := e.evaluateIfFunc(evaluatedBlock.Objects)
			if err != nil {
				return nil, err
			}

			for _, r := range results {
				if field, ok := r.(Field); ok {
					endResults = append(endResults, e.Data.GetColumn(field.FieldName))
				} else if block, ok := r.(Block); ok {
					endResults, err = e.Execute(block)
					if err != nil {
						return nil, err
					}
				}
			}

		case "":
			curObj := block.Objects[i]
			endResults = make([]interface{}, 0)
			if field, ok := curObj.(Field); ok {
				endResults = append(endResults, e.Data.GetColumn(field.FieldName))
			} else if block, ok := curObj.(Block); ok {
				var err error
				endResults, err = e.Execute(block)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return endResults, nil
}
