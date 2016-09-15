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

func transform(node sqlparser.SQLNode, t transformer) sqlparser.SQLNode {
	if node == nil {
		return nil
	}
	switch {
	case isSelectNode(node):
		return t.TransformSelect(node.(*sqlparser.Select))
	case isWhereNode(node):
		return t.TransformWhere(node.(*sqlparser.Where))
	case isComparisonExprNode(node):
		return t.TransformComparisonExpr(node.(*sqlparser.ComparisonExpr))
	case isAndExprNode(node):
		return t.TransformAndExpr(node.(*sqlparser.AndExpr))
	case isOrExprNode(node):
		return t.TransformOrExpr(node.(*sqlparser.OrExpr))
	case isNotExprNode(node):
		return t.TransformNotExpr(node.(*sqlparser.NotExpr))
	case isColNameNode(node):
		return t.TransformColName(node.(*sqlparser.ColName))
	case isNumValNode(node):
		return t.TransformNumVal(node.(sqlparser.NumVal))
	case isStrValNode(node):
		return t.TransformStrVal(node.(sqlparser.StrVal))
	case isBinaryValNode(node):
		return t.TransformBinaryVal(node.(sqlparser.BinaryVal))
	case isSelectExprsNode(node):
		return t.TransformSelectExprs(node.(sqlparser.SelectExprs))
	case isStarExprNode(node):
		return t.TransformStarExpr(node.(*sqlparser.StarExpr))
	case isNonStarExprNode(node):
		return t.TransformNonStarExpr(node.(*sqlparser.NonStarExpr))
	case isValTupleNode(node):
		return t.TransformValTuple(node.(sqlparser.ValTuple))
	case isParenBoolExprNode(node):
		return t.TransformParenBoolExpr(node.(*sqlparser.ParenBoolExpr))
	case isLimitNode(node):
		return t.TransformLimit(node.(*sqlparser.Limit))
	case isFuncExprNode(node):
		return t.TransformFuncExpr(node.(*sqlparser.FuncExpr))
	case isRangeCondNode(node):
		return t.TransformRangeCond(node.(*sqlparser.RangeCond))
	case isDDLNode(node):
		return t.TransformDDL(node.(*sqlparser.DDL))
	case isUnionNode(node):
		return t.TransformUnion(node.(*sqlparser.Union))
	case isInsertNode(node):
		return t.TransformInsert(node.(*sqlparser.Insert))
	case isValuesNode(node):
		return t.TransformValues(node.(sqlparser.Values))
	case isDeleteNode(node):
		return t.TransformDelete(node.(*sqlparser.Delete))
	case isTableExprsNode(node):
		return t.TransformTableExprs(node.(sqlparser.TableExprs))
	case isAliasedTableExprNode(node):
		return t.TransformAliasedTableExpr(node.(*sqlparser.AliasedTableExpr))
	case isTableNameNode(node):
		return t.TransformTableName(node.(*sqlparser.TableName))
	case isJoinTableExprNode(node):
		return t.TransformJoinTableExpr(node.(*sqlparser.JoinTableExpr))
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
	if node == nil {
		return false
	}
	return reflect.TypeOf(node) == ty
}
