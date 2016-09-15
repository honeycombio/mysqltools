package normalizer

import (
	"strings"

	"github.com/toshok/sqlparser"
)

type Parser struct {
}

func (n *Parser) NormalizeQuery(q string) string {
	sqlAST, err := sqlparser.Parse(q)
	if err != nil {
		return ""
	}

	newAST := transform(sqlAST, n)

	return strings.ToLower(sqlparser.String(newAST))
}

// QuestionMarkExpr is a special SQLNode used to render '?'.  we replace literal values with this in our transformer
type QuestionMarkExpr struct {
}

func (q *QuestionMarkExpr) Format(buf *sqlparser.TrackedBuffer) {
	buf.Myprintf("?")
}

func (*QuestionMarkExpr) IExpr()    {}
func (*QuestionMarkExpr) IValExpr() {}

func (n *Parser) TransformSelect(node *sqlparser.Select) sqlparser.SQLNode {
	if node == nil {
		return nil
	}
	node.SelectExprs, _ = transform(node.SelectExprs, n).(sqlparser.SelectExprs)
	node.Where, _ = transform(node.Where, n).(*sqlparser.Where)
	node.Limit, _ = transform(node.Limit, n).(*sqlparser.Limit)
	return node
}
func (n *Parser) TransformSelectExprs(node sqlparser.SelectExprs) sqlparser.SQLNode {
	var newSlice sqlparser.SelectExprs
	for _, se := range node {
		selectExpr, _ := transform(se, n).(sqlparser.SelectExpr)
		newSlice = append(newSlice, selectExpr)
	}
	return newSlice
}
func (n *Parser) TransformUnion(node *sqlparser.Union) sqlparser.SQLNode {
	if node == nil {
		return nil
	}
	node.Left, _ = transform(node.Left, n).(sqlparser.SelectStatement)
	node.Right, _ = transform(node.Right, n).(sqlparser.SelectStatement)
	return node
}
func (n *Parser) TransformInsert(node *sqlparser.Insert) sqlparser.SQLNode {
	if node == nil {
		return nil
	}
	// remove comments
	node.Comments = make([][]byte, 0)
	node.Rows, _ = transform(node.Rows, n).(sqlparser.InsertRows)
	// XXX(toshok) not yet node.OnDup, _ = transform(node.OnDup, n).(sqlparser.OnDup)
	return node
}
func (n *Parser) TransformUpdate(node *sqlparser.Update) sqlparser.SQLNode {
	if node == nil {
		return nil
	}
	// remove comments
	node.Comments = make([][]byte, 0)

	node.Exprs, _ = transform(node.Exprs, n).(sqlparser.UpdateExprs)
	node.Where, _ = transform(node.Where, n).(*sqlparser.Where)
	node.Limit, _ = transform(node.Limit, n).(*sqlparser.Limit)
	return node
}
func (n *Parser) TransformDelete(node *sqlparser.Delete) sqlparser.SQLNode {
	if node == nil {
		return nil
	}
	// remove comments
	node.Comments = make([][]byte, 0)

	node.Where, _ = transform(node.Where, n).(*sqlparser.Where)
	node.Limit, _ = transform(node.Limit, n).(*sqlparser.Limit)
	return node
}
func (n *Parser) TransformSet(node *sqlparser.Set) sqlparser.SQLNode {
	if node == nil {
		return nil
	}
	// remove comments
	node.Comments = make([][]byte, 0)
	node.Exprs, _ = transform(node.Exprs, n).(sqlparser.UpdateExprs)
	return node
}
func (n *Parser) TransformDDL(node *sqlparser.DDL) sqlparser.SQLNode { return node }
func (n *Parser) TransformColumnDefinition(node *sqlparser.ColumnDefinition) sqlparser.SQLNode {
	return node
}
func (n *Parser) TransformCreateTable(node *sqlparser.CreateTable) sqlparser.SQLNode {
	return node
}
func (n *Parser) TransformStarExpr(node *sqlparser.StarExpr) sqlparser.SQLNode { return node }
func (n *Parser) TransformNonStarExpr(node *sqlparser.NonStarExpr) sqlparser.SQLNode {
	if node == nil {
		return nil
	}
	node.Expr, _ = transform(node.Expr, n).(sqlparser.Expr)
	return node
}
func (n *Parser) TransformAliasedTableExpr(node *sqlparser.AliasedTableExpr) sqlparser.SQLNode {
	return node
}
func (n *Parser) TransformTableName(node *sqlparser.TableName) sqlparser.SQLNode { return node }

//func (n *Parser) TransformParentTableExpr(node *sqlparser.ParentTableExpr) sqlparser.SQLNode {
//	return node
//}
func (n *Parser) TransformJoinTableExpr(node *sqlparser.JoinTableExpr) sqlparser.SQLNode {
	if node == nil {
		return nil
	}
	node.On, _ = transform(node.On, n).(sqlparser.BoolExpr)
	return node
}
func (n *Parser) TransformIndexHints(node *sqlparser.IndexHints) sqlparser.SQLNode /* needed? */ {
	return node
}
func (n *Parser) TransformWhere(node *sqlparser.Where) sqlparser.SQLNode {
	if node == nil {
		return nil
	}

	node.Expr, _ = transform(node.Expr, n).(sqlparser.BoolExpr)
	return node
}
func (n *Parser) TransformAndExpr(node *sqlparser.AndExpr) sqlparser.SQLNode {
	if node == nil {
		return nil
	}
	node.Left, _ = transform(node.Left, n).(sqlparser.BoolExpr)
	node.Right, _ = transform(node.Right, n).(sqlparser.BoolExpr)
	return node
}
func (n *Parser) TransformOrExpr(node *sqlparser.OrExpr) sqlparser.SQLNode {
	if node == nil {
		return nil
	}
	node.Left, _ = transform(node.Left, n).(sqlparser.BoolExpr)
	node.Right, _ = transform(node.Right, n).(sqlparser.BoolExpr)
	return node
}
func (n *Parser) TransformNotExpr(node *sqlparser.NotExpr) sqlparser.SQLNode {
	if node == nil {
		return nil
	}
	node.Expr, _ = transform(node.Expr, n).(sqlparser.BoolExpr)
	return node
}
func (n *Parser) TransformParenBoolExpr(node *sqlparser.ParenBoolExpr) sqlparser.SQLNode {
	if node == nil {
		return nil
	}
	newExpr, _ := transform(node.Expr, n).(sqlparser.BoolExpr)
	node.Expr = newExpr
	return node
}
func (n *Parser) TransformComparisonExpr(node *sqlparser.ComparisonExpr) sqlparser.SQLNode {
	if node == nil {
		return nil
	}
	node.Left, _ = transform(node.Left, n).(sqlparser.ValExpr)
	node.Right, _ = transform(node.Right, n).(sqlparser.ValExpr)
	return node
}
func (n *Parser) TransformRangeCond(node *sqlparser.RangeCond) sqlparser.SQLNode {
	if node == nil {
		return nil
	}
	node.Left, _ = transform(node.Left, n).(sqlparser.ValExpr)
	node.From, _ = transform(node.From, n).(sqlparser.ValExpr)
	node.To, _ = transform(node.To, n).(sqlparser.ValExpr)
	return node
}
func (n *Parser) TransformNullCheck(node *sqlparser.NullCheck) sqlparser.SQLNode   { return node }
func (n *Parser) TransformExistsExpr(node *sqlparser.ExistsExpr) sqlparser.SQLNode { return node }
func (n *Parser) TransformBinaryVal(node sqlparser.BinaryVal) sqlparser.SQLNode {
	return &QuestionMarkExpr{}
}
func (n *Parser) TransformStrVal(node sqlparser.StrVal) sqlparser.SQLNode {
	return &QuestionMarkExpr{}
}
func (n *Parser) TransformNumVal(node sqlparser.NumVal) sqlparser.SQLNode {
	return &QuestionMarkExpr{}
}
func (n *Parser) TransformValArg(node *sqlparser.ValArg) sqlparser.SQLNode {
	return &QuestionMarkExpr{}
}
func (n *Parser) TransformValTuple(node sqlparser.ValTuple) sqlparser.SQLNode {
	var newSlice sqlparser.ValTuple
	for _, val := range node {
		valExpr, _ := transform(val, n).(sqlparser.ValExpr)
		newSlice = append(newSlice, valExpr)
	}
	return newSlice
}
func (n *Parser) TransformNullVal(node *sqlparser.NullVal) sqlparser.SQLNode { return node }
func (n *Parser) TransformColName(node *sqlparser.ColName) sqlparser.SQLNode { return node }
func (n *Parser) TransformSubquery(node *sqlparser.Subquery) sqlparser.SQLNode {
	if node == nil {
		return nil
	}
	node.Select, _ = transform(node.Select, n).(sqlparser.SelectStatement)
	return node
}
func (n *Parser) TransformBinaryExpr(node *sqlparser.BinaryExpr) sqlparser.SQLNode {
	if node == nil {
		return nil
	}
	node.Left, _ = transform(node.Left, n).(sqlparser.Expr)
	node.Right, _ = transform(node.Right, n).(sqlparser.Expr)
	return node
}
func (n *Parser) TransformUnaryExpr(node *sqlparser.UnaryExpr) sqlparser.SQLNode {
	if node == nil {
		return nil
	}
	node.Expr, _ = transform(node.Expr, n).(sqlparser.Expr)
	return node
}
func (n *Parser) TransformFuncExpr(node *sqlparser.FuncExpr) sqlparser.SQLNode {
	if node == nil {
		return nil
	}
	node.Exprs, _ = transform(node.Exprs, n).(sqlparser.SelectExprs)
	return node
}
func (n *Parser) TransformCaseExpr(node *sqlparser.CaseExpr) sqlparser.SQLNode {
	if node == nil {
		return nil
	}
	node.Expr, _ = transform(node.Expr, n).(sqlparser.ValExpr)
	for i, _ := range node.Whens {
		node.Whens[i], _ = transform(node.Whens[i], n).(*sqlparser.When)
	}
	node.Else, _ = transform(node.Else, n).(sqlparser.ValExpr)
	return node
}
func (n *Parser) TransformWhen(node *sqlparser.When) sqlparser.SQLNode {
	if node == nil {
		return nil
	}
	node.Cond, _ = transform(node.Cond, n).(sqlparser.BoolExpr)
	node.Val, _ = transform(node.Val, n).(sqlparser.ValExpr)
	return node
}
func (n *Parser) TransformOrder(node *sqlparser.Order) sqlparser.SQLNode { return node }
func (n *Parser) TransformLimit(node *sqlparser.Limit) sqlparser.SQLNode {
	if node == nil {
		return nil
	}
	node.Offset, _ = transform(node.Offset, n).(sqlparser.ValExpr)
	node.Rowcount, _ = transform(node.Rowcount, n).(sqlparser.ValExpr)
	return node
}
func (n *Parser) TransformUpdateExpr(node *sqlparser.UpdateExpr) sqlparser.SQLNode {
	if node == nil {
		return nil
	}
	node.Expr, _ = transform(node.Expr, n).(sqlparser.ValExpr)
	return node
}

func (n *Parser) TransformValues(node sqlparser.Values) sqlparser.SQLNode {
	var newSlice sqlparser.Values
	for _, rt := range node {
		rowTuple, _ := transform(rt, n).(sqlparser.RowTuple)
		newSlice = append(newSlice, rowTuple)
	}
	return newSlice
}
