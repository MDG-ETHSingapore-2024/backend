package operations

import "fmt"

type Operation int

const (
	EqualTo Operation = iota
	GreaterThan
	LesserThan
	GreaterThanEqualTo
	LesserThanEqualTo
	NotEqualTo
	OrderBy
)

var operationEnum = map[string]Operation{
	"eq":    EqualTo,
	"gt":    GreaterThan,
	"lt":    LesserThan,
	"gte":   GreaterThanEqualTo,
	"lte":   LesserThanEqualTo,
	"neq":   NotEqualTo,
	"order": OrderBy,
}

var operationStatement = map[Operation]string{
	EqualTo:            "=",
	GreaterThan:        ">",
	LesserThan:         "<",
	GreaterThanEqualTo: ">=",
	LesserThanEqualTo:  "<=",
	NotEqualTo:         "!=",
	// TODO: Remove orderby
	OrderBy: "",
}

type OperationObject struct {
	operation  Operation
	columnName string
	value      string
}

func GetOperation(opCode string) Operation {
	return operationEnum[opCode]
}

func GetOperationStatement(op OperationObject) string {
	s := fmt.Sprintf("%s %s %s", op.columnName, operationStatement[op.operation], op.value)
	return s
}