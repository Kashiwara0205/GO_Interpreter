package ast

import (
	"../token"
	"bytes"
	"strings"
)



/*Node--------------------------------------------------------------------------------------------*/
//特徴：全intercace型のトップ
//     ここに書かれたメソッドは継承しているインターフェースで必ず実現しなければいけない

type Node interface{
	TokenLiteral() string
	String() string
}
/*-------------------------------------------------------------------------------------------------*/



/*Staement_interface-----------------------------------------------------------------------------------*/
//特徴: Statementノードには、どんなStatementでも入る（LetとかReturnとか）

type Statement interface{
	Node
	statementNode()
}
/*-------------------------------------------------------------------------------------------------*/



/*Expresion_interface------------------------------------------------------------------------*/
// 特徴：hoge = 5の5を保持する
// 値の保持に何かと使用する。　

type Expression interface{
	Node
	expressionNode()
}
/*------------------------------------------------------------------------------------------*/



/*ProgramStatementノード----------------------------------------------------------------------*/
//特徴: 全てのStatementの親ノードに値するノード
//      LetやReturnなどのStatementを管理
//      抽象構文木の一番上にいる

type Program struct{
	// Statmentインターフェース型配列
	// この中には構造体のアドレスが入る
	Statements []Statement
}

func (p *Program) TokenLiteral() string{
	if len(p.Statements) > 0{
		return p.Statements[0].TokenLiteral()
	}else{
		return ""
	}
}

func (p *Program) String() string{
	var out bytes.Buffer

	// 1階層下のノードをぽこぽこ吐き出す
	for _, s := range p.Statements{
		out.WriteString(s.String())
	}

	return out.String()
}
/*----------------------------------------------------------------------------------------------*/



/*LetStatementノード------------------------------------------------------------------------------*/
// 特徴: Letという定義文を解析するためのノード
// let x = 5などの解析に使用
// Name →　hoge
// Value →　5

type LetStatement struct{
	Token token.Token
	// Nameのほうが*Identifierで固定なのは、変数名が絶対くるから
	Name *Identifier
	// Valueの方は*Identifierがくることもあれば、別の何かが来ることもある
	// そして、別の何かはString()を持っている
	Value Expression
}

// Statementインターフェース
func (ls *LetStatement) statementNode() {}
// LetStatementのTokenLiteral
func (ls *LetStatement) TokenLiteral() string {return ls.Token.Literal}
// String型
func (ls *LetStatement) String() string{
	var out bytes.Buffer

	// let
	out.WriteString(ls.TokenLiteral() + " ")
	// 変数名　NameはIdenfifiter型なので、別のstring型を呼び出してる
	out.WriteString(ls.Name.String())
	// = 
	out.WriteString(" = ")

	// 代入されるべき数字
	if ls.Value != nil{
		out.WriteString(ls.Value.String())
	}

	// ;
	out.WriteString(";")
	return out.String()
}

/*-----------------------------------------------------------------------------------------------------*/



/*IdentfierStatementノード----------------------------------------------------------------------------------------*/
// 特徴: LetStatementの子ノード
// let x = 5のｘの部分を記憶するために使用

// hoge = 5のhogeを束縛
// 変数
type Identifier struct{
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}
// IdentfierのTokenLiteral
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

// 変数名が帰ってくるだけ
func (i *Identifier) String() string {return i.Value}
/*----------------------------------------------------------------------------------------------------*/



/*ReturｎStatementノード----------------------------------------------------------------------------------------*/
// 特徴: rerutn　1 などのreturn文を解析するために使用するノード

type ReturnStatement struct{
	Token token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode(){}
func (rs *ReturnStatement) TokenLiteral() string{ return rs.Token.Literal }

func(rs *ReturnStatement) String() string{
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil{
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}
/*----------------------------------------------------------------------------------------*/



/*ExpresionStatementノード---------------------------------------------------------------------------*/
//特徴：  myvalue;のような式文解析に使用

type ExpressionStatement struct{
	Token token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

func (es *ExpressionStatement) String() string{
	if es.Expression != nil{
		return es.Expression.String()
	}
	return ""
}
/*-----------------------------------------------------------------------------------------------*/



/*IntegerLiteralノード---------------------------------------------------------------------------*/
//特徴：  5;のような式解析に使用

type IntegerLiteral struct{
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) TokenLiteral() string {return il.Token.Literal}
func (il *IntegerLiteral) String() string {return il.Token.Literal}
/*----------------------------------------------------------------------------------------------*/



/*PrefixExpressionノード--------------------------------------------------------------------------*/
// 特徴: -とか!とかの解析に使う　!true -5
type PrefixExpression struct{
	Token token.Token
	Operator string
	Right Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string {return pe.Token.Literal}
func (pe *PrefixExpression) String() string{
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}
/*-----------------------------------------------------------------------------------------------*/


/*InfixExpressionノード--------------------------------------------------------------------------*/
// 特徴: <Left: 前置型で解析してできたノード> <Operator: 演算子> <Right: 前置型で解析してできたノード>
//       で構成された中置型のノード。　1 + 1など 1 < 1などを管理する。
type InfixExpression struct{
	Token token.Token
	Left Expression
	Operator string
	Right Expression
}

func (oe *InfixExpression) expressionNode() {}
func (oe *InfixExpression) TokenLiteral() string {return oe.Token.Literal}
func (oe *InfixExpression) String() string{
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(oe.Left.String())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.String())
	out.WriteString(")")

	return out.String()
}
/*----------------------------------------------------------------------------------------------*/



/*Booleanノード-------------------------------------------------------------------------------------*/
// 特徴: 真偽値やねん
type Boolean struct{
	Token token.Token
	Value bool
}

func(b *Boolean) expressionNode() {}
func(b *Boolean) TokenLiteral() string {return b.Token.Literal}
func(b *Boolean) String() string {return b.Token.Literal}
/*------------------------------------------------------------------------------------------------*/



/*IfExpressionノード-----------------------------------------------------------------------------------------*/
// 特徴: if文解析ノード
type IfExpression struct{
	Token token.Token
	Condition Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string{
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil{
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}
/*------------------------------------------------------------------------------------------------*/


/*BlockStatementノード------------------------------------------------------------------------------*/
// 特徴: if文の{}←ここのとこ
type BlockStatement struct{
	Token token.Token
	Statements []Statement
}

func (bs *BlockStatement) StatementNode() {}
func (bs *BlockStatement) TokenLiteral() string {return bs.Token.Literal}
func (bs *BlockStatement) String() string{
	var out bytes.Buffer

	for _, s := range bs.Statements{
		out.WriteString(s.String())
	}

	return out.String()
}
/*-------------------------------------------------------------------------------------------------*/



/*FunctionLiteralノード------------------------------------------------------------------------------*/
// 特徴: 関数のやつ
type FunctionLiteral struct{
	Token token.Token
	Parameters []*Identifier
	Body *BlockStatement
}

func(fl *FunctionLiteral) expressionNode() {}
func(fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal}
func(fl *FunctionLiteral) String() string{
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters{
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}
/*--------------------------------------------------------------------------------------------------*/


/*CallExpressionノード-------------------------------------------------------------------------------*/
// 特徴: これは、add(1, 2)のようなメソッドを呼び出すための構文解析パーツである
type CallExpression struct{
	Token token.Token
	Function Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string{
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments{
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}
/*---------------------------------------------------------------------------------------------------*/


/*StringLiteralノード-------------------------------------------------------------------------------*/
type StringLiteral struct{
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal}
func (sl *StringLiteral) String() string {return sl.Token.Literal}
/*--------------------------------------------------------------------------------------------------*/

type ArrayLiteral struct{
	Token token.Token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal}
func (al *ArrayLiteral) String() string{
	var out bytes.Buffer

	elements := []string{}
	for _, el := range al.Elements{
		elements = append(elements, el.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ","))
	out.WriteString("]")

	return out.String()
}

type IndexExpression struct{
	Token token.Token
	Left Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode() {}
func (ie *IndexExpression) TokenLiteral() string {return ie.Token.Literal}
func (ie *IndexExpression) String() string{
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("]")

	return out.String()
}

type HashLiteral struct{
	Token token.Token
	Pairs map[Expression]Expression
}
func (hl *HashLiteral) expressionNode() {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String() string{
	var out bytes.Buffer
	pairs := []string{}
	for key, value := range hl.Pairs{
		pairs = append(pairs, key.String() + ":" + value.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}