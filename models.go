package main

type Object interface {
	IsObject() bool
}

type IBlock interface {
	IsBlock() bool
}

type IField interface {
	IsField() bool
}

type IFunction interface {
	IsFunction() bool
}

type IOperator interface {
	IsOperator() bool
}

type IValueSet interface {
	IsValueSet() bool
}

type Block struct {
	Objects           []Object
	InnerFunctionName string
}

func (o Block) IsBlock() bool  { return true }
func (o Block) IsObject() bool { return true }

type Field struct {
	FieldName string
}

func (o Field) IsField() bool  { return true }
func (o Field) IsObject() bool { return true }

type Function struct {
	ArgsString   string
	FunctionName string
}

func (o Function) IsFunction() bool { return true }
func (o Function) IsObject() bool   { return true }

type Operator struct {
	Operator     string
	OperatorType string
}

func (o Operator) IsOperator() bool { return true }
func (o Operator) IsObject() bool   { return true }

type ValueSet struct {
	Values []interface{}
}

func (o ValueSet) IsValueSet() bool { return true }
func (o ValueSet) IsObject() bool   { return true }
