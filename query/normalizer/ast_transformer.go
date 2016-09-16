package normalizer

import (
	"fmt"
	"log"
	"reflect"

	"github.com/toshok/sqlparser"
)

type transformer interface {
	TransformSelect(*sqlparser.Select) sqlparser.SQLNode
	TransformSelectExprs(sqlparser.SelectExprs) sqlparser.SQLNode
	TransformUnion(*sqlparser.Union) sqlparser.SQLNode
	TransformInsert(*sqlparser.Insert) sqlparser.SQLNode
	TransformUpdate(*sqlparser.Update) sqlparser.SQLNode
	TransformDelete(*sqlparser.Delete) sqlparser.SQLNode
	TransformSet(*sqlparser.Set) sqlparser.SQLNode
	TransformDDL(*sqlparser.DDL) sqlparser.SQLNode
	TransformColumnDefinition(*sqlparser.ColumnDefinition) sqlparser.SQLNode
	TransformCreateTable(*sqlparser.CreateTable) sqlparser.SQLNode
	TransformStarExpr(*sqlparser.StarExpr) sqlparser.SQLNode
	TransformNonStarExpr(*sqlparser.NonStarExpr) sqlparser.SQLNode
	TransformAliasedTableExpr(*sqlparser.AliasedTableExpr) sqlparser.SQLNode
	TransformTableName(*sqlparser.TableName) sqlparser.SQLNode
	//	TransformParentTableExpr(*sqlparser.ParentTableExpr) sqlparser.SQLNode
	TransformJoinTableExpr(*sqlparser.JoinTableExpr) sqlparser.SQLNode
	TransformWhere(*sqlparser.Where) sqlparser.SQLNode
	TransformIndexHints(*sqlparser.IndexHints) sqlparser.SQLNode // needed?
	TransformAndExpr(*sqlparser.AndExpr) sqlparser.SQLNode
	TransformOrExpr(*sqlparser.OrExpr) sqlparser.SQLNode
	TransformNotExpr(*sqlparser.NotExpr) sqlparser.SQLNode
	TransformParenBoolExpr(*sqlparser.ParenBoolExpr) sqlparser.SQLNode
	TransformComparisonExpr(*sqlparser.ComparisonExpr) sqlparser.SQLNode
	TransformRangeCond(*sqlparser.RangeCond) sqlparser.SQLNode
	TransformNullCheck(*sqlparser.NullCheck) sqlparser.SQLNode
	TransformExistsExpr(*sqlparser.ExistsExpr) sqlparser.SQLNode
	TransformBinaryVal(sqlparser.BinaryVal) sqlparser.SQLNode
	TransformStrVal(sqlparser.StrVal) sqlparser.SQLNode
	TransformNumVal(sqlparser.NumVal) sqlparser.SQLNode
	TransformValArg(*sqlparser.ValArg) sqlparser.SQLNode
	TransformValTuple(sqlparser.ValTuple) sqlparser.SQLNode
	TransformNullVal(*sqlparser.NullVal) sqlparser.SQLNode
	TransformColName(*sqlparser.ColName) sqlparser.SQLNode
	TransformSubquery(*sqlparser.Subquery) sqlparser.SQLNode
	TransformBinaryExpr(*sqlparser.BinaryExpr) sqlparser.SQLNode
	TransformUnaryExpr(*sqlparser.UnaryExpr) sqlparser.SQLNode
	TransformFuncExpr(*sqlparser.FuncExpr) sqlparser.SQLNode
	TransformCaseExpr(*sqlparser.CaseExpr) sqlparser.SQLNode
	TransformWhen(*sqlparser.When) sqlparser.SQLNode
	TransformOrder(*sqlparser.Order) sqlparser.SQLNode
	TransformLimit(*sqlparser.Limit) sqlparser.SQLNode
	TransformUpdateExpr(*sqlparser.UpdateExpr) sqlparser.SQLNode
	TransformValues(sqlparser.Values) sqlparser.SQLNode
	TransformTableExprs(sqlparser.TableExprs) sqlparser.SQLNode
}

var (
	selectType           reflect.Type = reflect.TypeOf((*sqlparser.Select)(nil))
	whereType            reflect.Type = reflect.TypeOf((*sqlparser.Where)(nil))
	comparisonExprType   reflect.Type = reflect.TypeOf((*sqlparser.ComparisonExpr)(nil))
	andExprType          reflect.Type = reflect.TypeOf((*sqlparser.AndExpr)(nil))
	orExprType           reflect.Type = reflect.TypeOf((*sqlparser.OrExpr)(nil))
	notExprType          reflect.Type = reflect.TypeOf((*sqlparser.NotExpr)(nil))
	colNameType          reflect.Type = reflect.TypeOf((*sqlparser.ColName)(nil))
	starExprType         reflect.Type = reflect.TypeOf((*sqlparser.StarExpr)(nil))
	nonStarExprType      reflect.Type = reflect.TypeOf((*sqlparser.NonStarExpr)(nil))
	parenBoolExprType    reflect.Type = reflect.TypeOf((*sqlparser.ParenBoolExpr)(nil))
	limitType            reflect.Type = reflect.TypeOf((*sqlparser.Limit)(nil))
	funcExprType         reflect.Type = reflect.TypeOf((*sqlparser.FuncExpr)(nil))
	rangeCondType        reflect.Type = reflect.TypeOf((*sqlparser.RangeCond)(nil))
	ddlType              reflect.Type = reflect.TypeOf((*sqlparser.DDL)(nil))
	unionType            reflect.Type = reflect.TypeOf((*sqlparser.Union)(nil))
	insertType           reflect.Type = reflect.TypeOf((*sqlparser.Insert)(nil))
	deleteType           reflect.Type = reflect.TypeOf((*sqlparser.Delete)(nil))
	aliasedTableExprType reflect.Type = reflect.TypeOf((*sqlparser.AliasedTableExpr)(nil))
	tableNameType        reflect.Type = reflect.TypeOf((*sqlparser.TableName)(nil))
	joinTableExprType    reflect.Type = reflect.TypeOf((*sqlparser.JoinTableExpr)(nil))

	numValType      reflect.Type = reflect.TypeOf((*sqlparser.NumVal)(nil)).Elem()
	strValType      reflect.Type = reflect.TypeOf((*sqlparser.StrVal)(nil)).Elem()
	binaryValType   reflect.Type = reflect.TypeOf((*sqlparser.BinaryVal)(nil)).Elem()
	selectExprsType reflect.Type = reflect.TypeOf((*sqlparser.SelectExprs)(nil)).Elem()
	valTupleType    reflect.Type = reflect.TypeOf((*sqlparser.ValTuple)(nil)).Elem()
	valuesType      reflect.Type = reflect.TypeOf((*sqlparser.Values)(nil)).Elem()
	tableExprsType  reflect.Type = reflect.TypeOf((*sqlparser.TableExprs)(nil)).Elem()
)

func transform(node sqlparser.SQLNode, t transformer) sqlparser.SQLNode {
	if node == nil {
		return nil
	}
	nodeType := reflect.TypeOf(node)
	switch nodeType {
	case selectType:
		return t.TransformSelect(node.(*sqlparser.Select))
	case whereType:
		return t.TransformWhere(node.(*sqlparser.Where))
	case comparisonExprType:
		return t.TransformComparisonExpr(node.(*sqlparser.ComparisonExpr))
	case andExprType:
		return t.TransformAndExpr(node.(*sqlparser.AndExpr))
	case orExprType:
		return t.TransformOrExpr(node.(*sqlparser.OrExpr))
	case notExprType:
		return t.TransformNotExpr(node.(*sqlparser.NotExpr))
	case colNameType:
		return t.TransformColName(node.(*sqlparser.ColName))
	case starExprType:
		return t.TransformStarExpr(node.(*sqlparser.StarExpr))
	case nonStarExprType:
		return t.TransformNonStarExpr(node.(*sqlparser.NonStarExpr))
	case parenBoolExprType:
		return t.TransformParenBoolExpr(node.(*sqlparser.ParenBoolExpr))
	case limitType:
		return t.TransformLimit(node.(*sqlparser.Limit))
	case funcExprType:
		return t.TransformFuncExpr(node.(*sqlparser.FuncExpr))
	case rangeCondType:
		return t.TransformRangeCond(node.(*sqlparser.RangeCond))
	case ddlType:
		return t.TransformDDL(node.(*sqlparser.DDL))
	case unionType:
		return t.TransformUnion(node.(*sqlparser.Union))
	case insertType:
		return t.TransformInsert(node.(*sqlparser.Insert))
	case deleteType:
		return t.TransformDelete(node.(*sqlparser.Delete))
	case aliasedTableExprType:
		return t.TransformAliasedTableExpr(node.(*sqlparser.AliasedTableExpr))
	case tableNameType:
		return t.TransformTableName(node.(*sqlparser.TableName))
	case joinTableExprType:
		return t.TransformJoinTableExpr(node.(*sqlparser.JoinTableExpr))
	case numValType:
		return t.TransformNumVal(node.(sqlparser.NumVal))
	case strValType:
		return t.TransformStrVal(node.(sqlparser.StrVal))
	case binaryValType:
		return t.TransformBinaryVal(node.(sqlparser.BinaryVal))
	case selectExprsType:
		return t.TransformSelectExprs(node.(sqlparser.SelectExprs))
	case valTupleType:
		return t.TransformValTuple(node.(sqlparser.ValTuple))
	case valuesType:
		return t.TransformValues(node.(sqlparser.Values))
	case tableExprsType:
		return t.TransformTableExprs(node.(sqlparser.TableExprs))
	default:
		log.Fatal(fmt.Sprintf("not handled %+v", reflect.TypeOf(node)))
		return nil
	}
}

func isSelectNode(node sqlparser.SQLNode) bool {
	return isType(node, reflect.TypeOf((*sqlparser.Select)(nil)))
}

func isWhereNode(node sqlparser.SQLNode) bool {
	return isType(node, reflect.TypeOf((*sqlparser.Where)(nil)))
}

func isComparisonExprNode(node sqlparser.SQLNode) bool {
	return isType(node, reflect.TypeOf((*sqlparser.ComparisonExpr)(nil)))
}

func isAndExprNode(node sqlparser.SQLNode) bool {
	return isType(node, reflect.TypeOf((*sqlparser.AndExpr)(nil)))
}

func isOrExprNode(node sqlparser.SQLNode) bool {
	return isType(node, reflect.TypeOf((*sqlparser.OrExpr)(nil)))
}

func isNotExprNode(node sqlparser.SQLNode) bool {
	return isType(node, reflect.TypeOf((*sqlparser.NotExpr)(nil)))
}

func isParenBoolExprNode(node sqlparser.SQLNode) bool {
	return isType(node, reflect.TypeOf((*sqlparser.ParenBoolExpr)(nil)))
}

func isColNameNode(node sqlparser.SQLNode) bool {
	return isType(node, reflect.TypeOf((*sqlparser.ColName)(nil)))
}

func isLimitNode(node sqlparser.SQLNode) bool {
	return isType(node, reflect.TypeOf((*sqlparser.Limit)(nil)))
}

func isFuncExprNode(node sqlparser.SQLNode) bool {
	return isType(node, reflect.TypeOf((*sqlparser.FuncExpr)(nil)))
}

func isNumValNode(node sqlparser.SQLNode) bool {
	return isType(node, reflect.TypeOf((*sqlparser.NumVal)(nil)).Elem())
}

func isStrValNode(node sqlparser.SQLNode) bool {
	return isType(node, reflect.TypeOf((*sqlparser.StrVal)(nil)).Elem())
}

func isBinaryValNode(node sqlparser.SQLNode) bool {
	return isType(node, reflect.TypeOf((*sqlparser.BinaryVal)(nil)).Elem())
}

func isSelectExprsNode(node sqlparser.SQLNode) bool {
	return isType(node, reflect.TypeOf((*sqlparser.SelectExprs)(nil)).Elem())
}

func isValTupleNode(node sqlparser.SQLNode) bool {
	return isType(node, reflect.TypeOf((*sqlparser.ValTuple)(nil)).Elem())
}

func isStarExprNode(node sqlparser.SQLNode) bool {
	return isType(node, reflect.TypeOf((*sqlparser.StarExpr)(nil)))
}

func isNonStarExprNode(node sqlparser.SQLNode) bool {
	return isType(node, reflect.TypeOf((*sqlparser.NonStarExpr)(nil)))
}

func isRangeCondNode(node sqlparser.SQLNode) bool {
	return isType(node, reflect.TypeOf((*sqlparser.RangeCond)(nil)))
}

func isDDLNode(node sqlparser.SQLNode) bool {
	return isType(node, reflect.TypeOf((*sqlparser.DDL)(nil)))
}

func isUnionNode(node sqlparser.SQLNode) bool {
	return isType(node, reflect.TypeOf((*sqlparser.Union)(nil)))
}

func isInsertNode(node sqlparser.SQLNode) bool {
	return isType(node, reflect.TypeOf((*sqlparser.Insert)(nil)))
}

func isValuesNode(node sqlparser.SQLNode) bool {
	return isType(node, reflect.TypeOf((*sqlparser.Values)(nil)).Elem())
}

func isDeleteNode(node sqlparser.SQLNode) bool {
	return isType(node, reflect.TypeOf((*sqlparser.Delete)(nil)))
}

func isTableExprsNode(node sqlparser.SQLNode) bool {
	return isType(node, reflect.TypeOf((*sqlparser.TableExprs)(nil)).Elem())
}

func isAliasedTableExprNode(node sqlparser.SQLNode) bool {
	return isType(node, reflect.TypeOf((*sqlparser.AliasedTableExpr)(nil)))
}

func isJoinTableExprNode(node sqlparser.SQLNode) bool {
	return isType(node, reflect.TypeOf((*sqlparser.JoinTableExpr)(nil)))
}

func isTableNameNode(node sqlparser.SQLNode) bool {
	return isType(node, reflect.TypeOf((*sqlparser.TableName)(nil)))
}

func isType(node sqlparser.SQLNode, ty reflect.Type) bool {
	return reflect.TypeOf(node) == ty
}
