package orm

const DefaultPageSize = 20

// xorm支持的操作词
const (
	XormOperatorIsEqual              = "="
	XormOperatorIsGreaterThan        = ">"
	XormOperatorIsGreaterThanOrEqual = ">="
	XormOperatorIsLessThan           = "<"
	XormOperatorIsLessThanOrEqual    = "<="
	XormOperatorIsIn                 = "IN"
	XormOperatorIsNotIn              = "NOT IN"
	XormOperatorIsLike               = "LIKE"
	XormOperatorIsNotLike            = "NOT LIKE"
	XormOperatorIsJsonContains       = "JSON_CONTAINS"
	XormOperatorIsNotEqual           = "<>"
	XormOperatorIsIsNull             = "IS NULL"
	XormOperatorIsIsNotNull          = "IS NOT NULL"
)

// Query通用操作标识
const (
	Operator_is_equal            = 0
	OperatorIsGreaterThan        = 1
	OperatorIsGreaterThanOrEqual = 2
	OperatorIsIn                 = 3
	OperatorIsNotIn              = 4
	OperatorIsLessThan           = 5
	OperatorIsLessThanOrEqual    = 6
	OperatorIsLike               = 7
	OperatorIsNotLike            = 8
	OperatorIsJsonContains       = 9
	OperatorIsNotEqual           = 10
	// OPERATOR_CONTAINS                 = 11;
	// OPERATOR_NOT_CONTAINS             = 12;
	OperatorIsIsNull    = 13
	OperatorIsIsNotNull = 14
	OperatorExists      = 15
)

var QueryOperatorToXorm = map[int]string{
	Operator_is_equal:            XormOperatorIsEqual,
	OperatorIsGreaterThan:        XormOperatorIsGreaterThan,
	OperatorIsGreaterThanOrEqual: XormOperatorIsGreaterThanOrEqual,
	OperatorIsIn:                 XormOperatorIsIn,
	OperatorIsNotIn:              XormOperatorIsNotIn,
	OperatorIsLessThan:           XormOperatorIsLessThan,
	OperatorIsLessThanOrEqual:    XormOperatorIsLessThanOrEqual,
	OperatorIsLike:               XormOperatorIsLike,
	OperatorIsNotLike:            XormOperatorIsNotLike,
	OperatorIsJsonContains:       XormOperatorIsJsonContains,
	OperatorIsNotEqual:           XormOperatorIsNotEqual,
	OperatorIsIsNull:             XormOperatorIsIsNull,
	OperatorIsIsNotNull:          XormOperatorIsIsNotNull,
}
