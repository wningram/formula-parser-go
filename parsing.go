package main

import (
	"fmt"
	"strings"
)

func getConditionalOperators() []string {
	return []string{"<", ">", "<>", "=", "<=", ">="}
}

func getTextOperators() []string {
	return []string{"&"}
}

func getArithmeticOperators() []string {
	return []string{"+", "-", "/", "^"}
}

func isTextOperator(str string) bool {
	for _, e := range getTextOperators() {
		if str == e {
			return true
		}
	}
	return false
}

func isArithmeticOperator(str string) bool {
	for _, e := range getArithmeticOperators() {
		if str == e {
			return true
		}
	}
	return false
}

func isConditionalOperator(str string) bool {
	for _, e := range getConditionalOperators() {
		if str == e {
			return true
		}
	}
	return false
}

type ParsingError struct {
	Msg     string
	ExStart int
	ExEnd   int
}

func (e *ParsingError) Error() string {
	return e.Msg
}

//type Parser interface {
//	getNextWord(block string, startNdx int) (string, int)
//	getNext(formulaSeg string) Object
//}

func GetNextWord(block string, startNdx int) (string, int) {
	if startNdx >= len(block) {
		return "", -1
	}

	if startNdx <= len(block)-2 && block[startNdx:startNdx+2] == "if" {
		return "if", startNdx + 2
	}

	if block[startNdx] == '[' {
		return "field_begin", startNdx + 1
	}

	if block[startNdx] == ']' {
		return "field_end", startNdx + 1
	}

	if block[startNdx] == '(' {
		return "func_begin", startNdx + 1
	}

	if block[startNdx] == ')' {
		return "func_end", startNdx + 1
	}

	if block[startNdx] == ' ' {
		return "space", startNdx + 1
	}

	if startNdx < len(block)-2 && isConditionalOperator(block[startNdx:startNdx+2]) {
		return "cond_op", startNdx + 2
	}

	if isConditionalOperator(string(block[startNdx])) {
		return "cond_op", startNdx + 1
	}

	if isTextOperator(string(block[startNdx])) {
		return "text_op", startNdx + 1
	}

	if isArithmeticOperator((string(block[startNdx]))) {
		return "arithmetic_op", startNdx + 1
	}

	if block[startNdx] == ',' {
		return "arg_sep", startNdx + 1
	}

	return string(block[startNdx]), startNdx + 1
}

func GetNext(formulaSeg string) (Object, int, error) {
	var (
		curNdx           = 0
		curWord, nextNdx = GetNextWord(formulaSeg, 0)
		isField          = false
		isFunc           = false
		isOp             = false
		endWord          = ""
		funcCloseCount   = 0
		funcArgStr       = ""
		fieldName        = ""
		operator         = ""
		curFunc          = ""
		loop             = true
	)

	for loop {
		switch curWord {
		case "func_begin":
			if isFunc {
				funcArgStr += formulaSeg[curNdx:nextNdx]
			}
			if isField {
				return nil, -1, &ParsingError{"Cannot begin a function within a field", curNdx, nextNdx}
			}
			endWord = "func_end"
			isFunc = true
			funcCloseCount++
		case "func_end":
			funcCloseCount--
			if funcCloseCount != 0 {
				funcArgStr += formulaSeg[curNdx:nextNdx]
			}
		case "field_begin":
			if !isFunc {
				endWord = "field_end"
				isField = true
			} else {
				funcArgStr += formulaSeg[curNdx:nextNdx]
			}
		case "if":
			if isFunc {
				funcArgStr += formulaSeg[curNdx:nextNdx]
			}
			curFunc = "if"
		default:
			if strings.HasSuffix(curWord, "_op") {
				if !isFunc {
					if isField {
						return nil, -1, &ParsingError{"Operator cannot exist within field name.", curNdx, nextNdx}
					}
					isOp = true
					operator = formulaSeg[curNdx:nextNdx]
					break
				} else {
					funcArgStr += formulaSeg[curNdx:nextNdx]
				}
			} else {
				if isFunc {
					funcArgStr += formulaSeg[curNdx:nextNdx]
				} else if isField {
					fieldName += formulaSeg[curNdx:nextNdx]
				}
			}
		}
		curNdx = nextNdx
		curWord, nextNdx = GetNextWord(formulaSeg, curNdx)
		// Determine if we should continue looping
		if nextNdx == -1 {
			loop = false
		}

		if isOp {
			loop = false
		}

		if curWord == endWord {
			if !isFunc {
				loop = false
			} else {
				if funcCloseCount == 0 && curWord == endWord {
					loop = false
				}
			}
		}
	}
	if isField {
		return Field{
			FieldName: fieldName,
		}, nextNdx, nil
	} else if isFunc {
		return Function{
			ArgsString:   funcArgStr,
			FunctionName: curFunc,
		}, nextNdx, nil
	} else if isOp {
		return Operator{
			Operator:     operator,
			OperatorType: curWord,
		}, nextNdx, nil
	} else {
		return nil, -1, fmt.Errorf("'Next' is not an operator, field or function.")
	}
}

func Parse(formulaSeg string) (Block, error) {
	var (
		objs              []Object
		newFormulaSeg     string = formulaSeg
		funcName          string
		obj, nextNdx, err = GetNext(formulaSeg)
	)

	for true {
		if funcObj, ok := obj.(Function); ok {
			funcName = funcObj.FunctionName
			block, err := Parse(funcObj.ArgsString)
			if err != nil {
				return Block{}, err
			}
			block.InnerFunctionName = funcName
			objs = append(objs, block)
		} else {
			objs = append(objs, obj)
		}

		if nextNdx == -1 {
			break
		}
		newFormulaSeg = newFormulaSeg[nextNdx:]
		if newFormulaSeg == "" {
			break
		}
		obj, nextNdx, err = GetNext(newFormulaSeg)
		if err != nil {
			// TODO: Need to see if putting error in block is correct way to handle errors
			return Block{}, err
		}
	}
	return Block{
		Objects:           objs,
		InnerFunctionName: "",
	}, nil
}
