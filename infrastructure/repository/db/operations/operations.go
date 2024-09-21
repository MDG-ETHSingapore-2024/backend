package operations

import "fmt"

const ConnectionString = "host=localhost port=5432 user=postgres password=dpsingh05 dbname=postgres sslmode=disable"

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
	OrderBy:            "",
}

type OperationObject struct {
	Operation  Operation
	ColumnName string
	Value      string
}

func GetOperation(opCode string) Operation {
	return operationEnum[opCode]
}

func GetOperationStatement(op OperationObject) string {
	s := fmt.Sprintf("%s %s %s", op.ColumnName, operationStatement[op.Operation], op.Value)
	return s
}
