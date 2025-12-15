package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"io"
)

func NewStruct() *ast.StructType {
	return &ast.StructType{
		Fields: &ast.FieldList{
			List: []*ast.Field{},
		},
	}

}

func AddFieldToStruct(s *ast.StructType, fieldName string, t ast.Expr) {
	s.Fields.List = append(s.Fields.List, &ast.Field{
		Names: []*ast.Ident{
			{
				Name: fieldName,
			},
		},
		Type: t,
	})
}

//func AddAnonFieldToStruct(s *ast.StructType, t ast.Expr) {
//	s.Fields.List = append(s.Fields.List, &ast.Field{
//		Names: []*ast.Ident{},
//		Type:  t,
//	})
//}

func Import(path ...string) ast.Decl {
	specs := make([]ast.Spec, len(path))
	for i, p := range path {
		specs[i] = &ast.ImportSpec{
			Path: StrLit(p),
		}
	}
	return &ast.GenDecl{
		Tok:   token.IMPORT,
		Specs: specs,
	}
}

func Ident(name string) *ast.Ident {
	return ast.NewIdent(name)
}

func VarStmt(name string, t ast.Expr) *ast.DeclStmt {
	return &ast.DeclStmt{Decl: Var(name, t)}
}

func Var(name string, t ast.Expr) ast.Decl {
	return &ast.GenDecl{
		Tok: token.VAR,
		Specs: []ast.Spec{
			&ast.ValueSpec{
				Names: []*ast.Ident{
					{
						Name: name,
					},
				},
				Type: t,
			},
		},
	}
}

func Slice(elementType ast.Expr) ast.Expr {
	return &ast.ArrayType{
		Len: nil, // nil Length means it's a slice, not an array
		Elt: elementType,
	}
}

func NumLit(n int) *ast.BasicLit {
	return NumLitStr(fmt.Sprint(n))
}

func NumLitHex(n int) *ast.BasicLit {
	return NumLitStr(fmt.Sprintf("%X", n))
}

func NumLitStr(n string) *ast.BasicLit {
	return &ast.BasicLit{
		ValuePos: 0,
		Kind:     token.INT,
		Value:    n,
	}
}

func Array(elementType ast.Expr, len int) ast.Expr {
	return &ast.ArrayType{
		Len: NumLit(len),
		Elt: elementType,
	}
}

func TypeDecl(name string, t ast.Expr) ast.Decl {
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: &ast.Ident{
					Name: name,
				},
				Type: t,
			},
		},
	}
}

func Selector(x, sel string) ast.Expr {
	return &ast.SelectorExpr{
		X:   Ident(x),
		Sel: Ident(sel),
	}
}

func SelectorExprAndStr(x ast.Expr, sel string) ast.Expr {
	return &ast.SelectorExpr{
		X:   x,
		Sel: Ident(sel),
	}
}

func Pointer(elementType ast.Expr) ast.Expr {
	return &ast.StarExpr{
		X: elementType,
	}
}

func Call(fn ast.Expr, args ...ast.Expr) *ast.CallExpr {
	return &ast.CallExpr{
		Fun:  fn,
		Args: args,
	}
}

func Assign121(lhs, rhs ast.Expr) *ast.AssignStmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{lhs},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{rhs},
	}
}

func Assign(lhs, rhs []ast.Expr) *ast.AssignStmt {
	return &ast.AssignStmt{
		Lhs: lhs,
		Tok: token.ASSIGN,
		Rhs: rhs,
	}
}

func Define(lhs, rhs []ast.Expr) *ast.AssignStmt {
	return &ast.AssignStmt{
		Lhs: lhs,
		Tok: token.DEFINE,
		Rhs: rhs,
	}
}

func Define121(lhs, rhs ast.Expr) *ast.AssignStmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{lhs},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{rhs},
	}
}

func Exprs(exprs ...ast.Expr) []ast.Expr {
	return exprs
}

func Stmts(stmts ...ast.Stmt) []ast.Stmt {
	return stmts
}

func ExprStmt(expr ast.Expr) *ast.ExprStmt {
	return &ast.ExprStmt{
		X: expr,
	}
}

func If(cond ast.Expr, body *ast.BlockStmt) *ast.IfStmt {
	return &ast.IfStmt{
		Cond: cond,
		Body: body,
	}
}

func NotEquals(lhs, rhs ast.Expr) *ast.BinaryExpr {
	return &ast.BinaryExpr{
		X:  lhs,
		Op: token.NEQ,
		Y:  rhs,
	}
}

func Equals(lhs, rhs ast.Expr) *ast.BinaryExpr {
	return &ast.BinaryExpr{
		X:  lhs,
		Op: token.EQL,
		Y:  rhs,
	}
}

func Nil() ast.Expr {
	return Ident("nil")
}

func NewFile(name string) *ast.File {
	return &ast.File{
		Name: &ast.Ident{
			Name: name,
		},
		Decls: []ast.Decl{},
	}
}

type NameAndType struct {
	Name string
	Type ast.Expr
}

func NewFunc(name string, arguments []NameAndType, returns []NameAndType) *ast.FuncDecl {
	params := &ast.FieldList{
		List: make([]*ast.Field, len(arguments)),
	}
	for i, arg := range arguments {
		params.List[i] = &ast.Field{
			Names: []*ast.Ident{Ident(arg.Name)},
			Type:  arg.Type,
		}
	}

	results := &ast.FieldList{
		List: make([]*ast.Field, len(returns)),
	}
	for i, ret := range returns {
		results.List[i] = &ast.Field{
			Names: []*ast.Ident{{Name: ret.Name}},
			Type:  ret.Type,
		}
	}

	return &ast.FuncDecl{
		Name: Ident(name),
		Type: &ast.FuncType{
			Params:  params,
			Results: results,
		},
		Body: &ast.BlockStmt{},
	}
}

func NewFuncWithReceiver(name, receiverName, receiver string, arguments []NameAndType, returns []NameAndType) *ast.FuncDecl {
	f := NewFunc(name, arguments, returns)
	f.Recv = &ast.FieldList{
		List: []*ast.Field{
			{
				Names: []*ast.Ident{{Name: receiverName}},
				Type:  &ast.Ident{Name: receiver},
			},
		},
	}
	return f
}

func NewBlock(statements []ast.Stmt) *ast.BlockStmt {
	return &ast.BlockStmt{
		List: statements,
	}
}

func NewBlockEllipsis(statements ...ast.Stmt) *ast.BlockStmt {
	return &ast.BlockStmt{
		List: statements,
	}
}

//func MergeBlocks(blocks ...*ast.BlockStmt) *ast.BlockStmt {
//	list := make([]ast.Stmt, 0)
//	for _, block := range blocks {
//		if block != nil {
//			list = append(list, block.List...)
//		}
//	}
//	return &ast.BlockStmt{
//		List: list,
//	}
//}

func Return(exprs ...ast.Expr) *ast.ReturnStmt {
	return &ast.ReturnStmt{
		Results: exprs,
	}
}

func AppendDecl(file *ast.File, decl ast.Decl) {
	file.Decls = append(file.Decls, decl)
}

func PrintToFile(f *ast.File, w io.Writer) error {
	return format.Node(w, token.NewFileSet(), f)
}

func IfErrNil() *ast.IfStmt {
	return If(
		NotEquals(Ident("err"), Nil()), NewBlockEllipsis(Return()),
	)
}

func AddrOf(e ast.Expr) *ast.UnaryExpr {
	return &ast.UnaryExpr{
		Op: token.AND,
		X:  e,
	}
}

func DeRef(e ast.Expr) *ast.StarExpr {
	return &ast.StarExpr{
		X: e,
	}
}

func ForRange(K, X ast.Expr, block *ast.BlockStmt) *ast.RangeStmt {
	return &ast.RangeStmt{
		Key:   K,            // The loop variable "i"
		Value: nil,          // No second variable (common for index-only range loops)
		X:     X,            // The collection being ranged over "l"
		Tok:   token.DEFINE, // ":=" token
		Body:  block,
	}

}

func Index(X, I ast.Expr) *ast.IndexExpr {
	return &ast.IndexExpr{
		X:     X,
		Index: I,
	}
}

func MultiIndex(X ast.Expr, Low, High, Max ast.Expr) *ast.SliceExpr {
	return &ast.SliceExpr{
		X:      X,
		Low:    Low,
		High:   High,
		Max:    Max,
		Slice3: Max != nil,
	}
}

func SwitchStmt(Tag ast.Expr, block *ast.BlockStmt) *ast.SwitchStmt {
	return &ast.SwitchStmt{
		Tag:  Tag,
		Body: block,
	}
}

//func SwitchStmtInit(init ast.Stmt, block *ast.BlockStmt) *ast.SwitchStmt {
//	return &ast.SwitchStmt{
//		Init: init,
//		Body: block,
//	}
//}

func StrLit(s string) *ast.BasicLit {
	return &ast.BasicLit{
		Kind:  token.STRING,
		Value: `"` + s + `"`,
	}
}

func Case(exprs []ast.Expr, stmts []ast.Stmt) *ast.CaseClause {
	return &ast.CaseClause{
		List: exprs,
		Body: stmts,
	}
}

func MapType(key, value ast.Expr) ast.Expr {
	return &ast.MapType{
		Key:   key,
		Value: value,
	}
}

func CompLit(t ast.Expr, elts []ast.Expr) *ast.CompositeLit {
	return &ast.CompositeLit{
		Type: t,
		Elts: elts,
	}
}

//func VarAssign(name string, value ast.Expr) *ast.GenDecl {
//	return &ast.GenDecl{
//		Tok: token.VAR,
//		Specs: []ast.Spec{
//			&ast.ValueSpec{
//				Names: []*ast.Ident{
//					{
//						Name: name,
//					},
//				},
//				Values: []ast.Expr{value},
//			},
//		},
//	}
//}

func KeyValueExpr(key, value ast.Expr) *ast.KeyValueExpr {
	return &ast.KeyValueExpr{
		Key:   key,
		Value: value,
	}
}

func TypeAssert(x, t ast.Expr) *ast.TypeAssertExpr {
	return &ast.TypeAssertExpr{
		X:    x,
		Type: t,
	}
}

func TypeSwitch(x ast.Stmt, block *ast.BlockStmt) *ast.TypeSwitchStmt {
	return &ast.TypeSwitchStmt{
		Assign: x,
		Body:   block,
	}
}

func Not(x ast.Expr) *ast.UnaryExpr {
	return &ast.UnaryExpr{
		Op: token.NOT,
		X:  x,
	}
}

func Sub(lhs, rhs ast.Expr) *ast.BinaryExpr {
	return &ast.BinaryExpr{
		X:  lhs,
		Op: token.SUB,
		Y:  rhs,
	}
}

func BinAnd(lhs, rhs ast.Expr) *ast.BinaryExpr {
	return &ast.BinaryExpr{
		X:  lhs,
		Op: token.AND,
		Y:  rhs,
	}
}
