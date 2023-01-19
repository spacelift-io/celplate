package macros

import (
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common"
	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/parser"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

func GetJoinMacro() parser.Macro {
	joinMacro := cel.NewReceiverMacro("join", 1,
		func(meh cel.MacroExprHelper, iterRange *exprpb.Expr, args []*exprpb.Expr) (*exprpb.Expr, *common.Error) {
			delim := args[0]

			// Convert the list elements to string
			// So [1, 2, 3].join(",") becomes ["1", "2", "3"].join(",")
			var newIterRangeElements []*exprpb.Expr
			for _, elem := range iterRange.GetListExpr().GetElements() {
				newIterRangeElements = append(newIterRangeElements, meh.GlobalCall("string", elem))
			}
			stringRange := meh.NewList(newIterRangeElements...)

			iterIdent := meh.Ident("__iter__")
			accuIdent := meh.AccuIdent()
			accuInit := meh.LiteralString("")
			condition := meh.LiteralBool(true)
			step := meh.GlobalCall(
				// __result__.size() > 0 ? __result__  + delim + __iter__ : __iter__
				operators.Conditional,
				meh.GlobalCall(operators.Greater, meh.ReceiverCall("size", accuIdent), meh.LiteralInt(0)),
				meh.GlobalCall(operators.Add, meh.GlobalCall(operators.Add, accuIdent, delim), iterIdent),
				iterIdent)
			return meh.Fold(
				"__iter__",
				stringRange,
				accuIdent.GetIdentExpr().GetName(),
				accuInit,
				condition,
				step,
				accuIdent), nil
		})

	return joinMacro
}
