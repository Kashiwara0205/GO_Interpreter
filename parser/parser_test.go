package parser

import (
	"testing"
	"../ast"
	"../lexer"
	"fmt"
)

func TestLetStatement(t *testing.T){
	// testにかけるソースコード入力
	input :=`
	let x = 5;
	let y = 10;
	let foobar = 838383;
	`

	// 字句解析導入
	l := lexer.New(input)
	// パーサーのNew
	p := New(l)

	// トークンを元に字句解析
	// parseする
	// parseから、各構造体にトークンやらvalueやらがセットされ抽象構文木が作成される
	program := p.ParseProgram()
	checkParseErrors(t, p)

	// 抽象構文木がnilであるか
	if program == nil{
		// nil返却
		t.Fatalf("ParseProgram() retuned nil")
	}
	// Statmentsはinterfeace型の配列である
	//　３つは、letの定義文が３つあることを示す
	if len(program.Statements) != 3{
		t.Fatalf("program.Statements doen not contain 3 statements got=%d",
		len(program.Statements))
	}

	// Identfierは変数
	tests := []struct{
		expectedIdentifier string
	}{
		// letで定義されている変数
		{"x"},
		{"y"},
		{"foobar"},
	}

	// iは、ただのeach_witdh_index
	for i, tt := range tests{
		// ルートからstatmentツリーをゲット
		// statmentツリーには、各アドレスが入ってる
		stmt := program.Statements[i]

		if !testLetStatement(t, stmt, tt.expectedIdentifier){
			return 
		}
	}
}


func checkParseErrors(t *testing.T, p *Parser){
	errors := p.Errors()
	if len(errors) == 0{
		return 
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors{
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}

func testLetStatement(t *testing.T, s ast.Statement, name string)bool{
	// letが定義されているか
	// TokenLiteralで帰ってくるのは token.goで定義したCONSTのやつ

	// LetStatementnのTokenLiteralメソッド使用
	if s.TokenLiteral() != "let"{
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	// sには、各、letStmtに対するアドレスが入っている
	// sはinterface型なのでLetStatement構造体型に変換をかける
	// interface型の場合、そこで定義されたメソッドしか使えなくて構造体の要素は参照できない
	// 下記のやりかたで変換がかかる
	letStmt, ok := s.(*ast.LetStatement)
	if !ok{
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	// 変数名　hogeとか
	if letStmt.Name.Value != name{
		t.Errorf("letStmt.Name Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	// 上のhogeにもtokenがある
	// そのtokenのリテラルをみる　
	// ぶっつちゃけ変数名と同じである
	if letStmt.Name.TokenLiteral() != name{
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s",
			name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T){
	input :=`
		return 5;
		return 10;
		return 993322;
	`
	// 字句解析
	l := lexer.New(input)
	// 構文解析の初期化
	p := New(l)

	// 抽象構文木の生成
	program := p.ParseProgram()
	// 抽象校分木にエラーがあるか
	checkParseErrors(t, p)

	// return文が３つか
	if len(program.Statements)!=3{
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	for _, stmt := range program.Statements{
		// ReturnStatement構造体を代入
		// Statment構造体をReturnStatment構造体に変換
		returnStmt, ok := stmt.(*ast.ReturnStatement)

		if !ok{
			t.Errorf("stmt not *ast.returnStatement. got=%T", stmt)
		}

		if returnStmt.TokenLiteral() != "return"{
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q",
				returnStmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T){
	input := "foobar;"

	// 字句解析
	l := lexer.New(input)
	// New ⇒　prefixParseFnsにIdentifierのアドレスが登録されたり
	//         Ast作成の前処理実行
	p := New(l)
	// Ast作成
	prgoram := p.ParseProgram()
	//　エラー探知
	checkParseErrors(t, p)

	// 識別子が１つという意味
	if len(prgoram.Statements) != 1{
		t.Fatalf("Program has not enough statements. got=%d",
				len(prgoram.Statements))
	}
	// ExpressionStatementに変換をかける
	stmt, ok := prgoram.Statements[0].(*ast.ExpressionStatement)
	if !ok{
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			prgoram.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)

	if !ok{
		t.Fatalf("exp not *ast.Identifier. got=%s",
			stmt.Expression)
	}
	if ident.Value != "foobar"{
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)

	}
	if ident.TokenLiteral() != "foobar"{
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar",
		ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T){
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1{
		t.Fatalf("program has not enought statements got=%d",
					len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok{
		t.Fatalf("program.Statements[0] is not ast.ExprssionStatement. got=%T",
				program.Statements[0])
	}
	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok{
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got = %d", 5, literal.Value)
	}

	if literal.TokenLiteral() != "5"{
		t.Errorf("literaal.TokenLiteral not %s. got=%s", "5", literal.TokenLiteral())
	}
}


func TestParsingPrefixExpressions(t *testing.T){
	prefixTests := []struct{
		input string
		operator string
		value interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	// prefixTestにある回数だけ字句解析して、Astも、そのつど作成
	for _, tt := range prefixTests{
		// 字句解析
		l := lexer.New(tt.input)
		// astの下準備
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		//!5; ⇒　-15; で１回ずつ
		if len(program.Statements) != 1{
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				l, len(program.Statements))
		}

		// 式ツリーを変換して代入
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok{
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got =%T", stmt.Expression)
		}
		// 式ツリーから、PrefixExpressionに変換
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok{
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		// Operettorの識別子があっているか
		if exp.Operator != tt.operator{
			t.Fatalf("exp.Opertor is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool{
	// PrefixExpressionノードのRightにひっついてるIntgerLiteralをもってくる
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok{
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
	}
	if integ.Value != value{
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}
	if integ.TokenLiteral() != fmt.Sprintf("%d", value){
		t.Errorf("integ.TokenLiteral not %d. got=%s", value,
			integ.TokenLiteral())
		return false
	}
	return true
}

func TestParsingInfixExpressions(t *testing.T){
	infixTests := []struct{
		input string
		leftVale interface{}
		operator string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests{
		// 字句解析
		l := lexer.New(tt.input)
		// Astを構築するための準備
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1{
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				l, len(program.Statements))
		}

		// 式のステートメント変換
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok{
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		//中置型解析数であるInfixExoressionノードに変換
		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok{
			t.Fatalf("exp is not ast.InfixExpression. got=%T", stmt.Expression)
		}


		// 演算子がただしく格納されているか
		if exp.Operator != tt.operator{
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}

		// 右のvalue検証と左のvalue検証が消えてこいつに変わった
		// 右と左に真偽値や数値がきても関係なくテストしてくれる
		if !testInfixExpresison(t, stmt.Expression, tt.leftVale,
				tt.operator, tt.rightValue){
					return
				}
	}
}

func TestOpreatorPrecedenceParsiong(t *testing.T){
	tests := []struct{
		input string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			" a * b / c",
			"((a * b) / c)",
		},


		{
			"a + b / c",
			"(a + (b / c))",
		},

		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{	
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c]* d",
			"((a * ([1,2,3,4][(b * c)]) * d)",

		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2]), (b[1], (2 * ([1,2][1]))",
		},

	}
	for _, tt := range tests{
		// 字句解析開始
		l := lexer.New(tt.input)
		// astノード構築初期設定
		p := New(l)
		// astノード構築
		program := p.ParseProgram()
		// エラー探知
		checkParseErrors(t, p)
		
		// 中置型の構文木が、ちゃんと優先度守って括弧つけられているのかを見る
		actual := program.String()
		if actual != tt.expected{
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool{
	ident, ok := exp.(*ast.Identifier)
	if !ok{
		t.Errorf("exp not *ast.Idenfitifer. got=%T", exp)
	}

	if ident.Value != value{
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}
	if ident.TokenLiteral() != value{
		t.Errorf("ident.TokenLitral not %s. got=%s", value,
				ident.TokenLiteral())
		return false
	}
	return true
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
)bool{
	// 右と左の型によりテストの飛び先を変更
	// 飛び先が変わると,それぞれExpresionで変換をかけてる箇所の枝が変わる
	switch v := expected.(type){
	case int:
		// integerliteralのテストへ
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
        // integerliteralのテストへ
		return testIntegerLiteral(t, exp, v)
	case string:
		//　Identのテストへ
		return testIdentifier(t, exp, v)
	case bool:
		// Boolのテストへ
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got =%T", exp)
	return false
}

func testInfixExpresison(t *testing.T, exp ast.Expression, left interface{},
	operator string, right interface{})bool{
		// 	中置型のinfixノードに置換
		opExp, ok := exp.(*ast.InfixExpression)
		if !ok{
			t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		}
		//　左のvalue検証
		if !testLiteralExpression(t, opExp.Left, left){
			return false
		}
		// 命令
		if opExp.Operator != operator{
			return false
		}
		// 右のvalue検証
		if !testLiteralExpression(t, opExp.Right, right){
			return false
		}

		return true
	}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)
		
		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		boolean, ok := stmt.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("exp not *ast.Boolean. got=%T", stmt.Expression)
		}
		if boolean.Value != tt.expectedBoolean {
			t.Errorf("boolean.Value not %t. got=%t", tt.expectedBoolean,
				boolean.Value)
		}
	}
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool{
	bo, ok := exp.(*ast.Boolean)
	if !ok{
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if bo.Value != value{
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value){
		t.Errorf("bo.TokenLiteral not %t. got=%s",
		value, bo.TokenLiteral())
		return false
	}

	return true
}

func TestIfExpression(t *testing.T){
	input := `if (x < y) { x }`

	// 字句解析
	l := lexer.New(input)
	// 構文解析の準備を整える
	p := New(l)
	// 構文解析完了
	program := p.ParseProgram()
	// エラーチェック
	checkParseErrors(t, p)

	// ノードが１本だけかどうか
	if len(program.Statements) != 1{
		t.Fatalf("program.Statements does not contain %d statements . got=%d\n",
			1, len(program.Statements))
	}

	// 式ノードに変換
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	// 変換がかけれたかどうか
	if !ok{
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	// IfExpressionに変換
	exp, ok := stmt.Expression.(*ast.IfExpression)
	// 変換がかけれたかどうか
	if !ok{
		t.Fatalf("stmt.Expressio is not ast.IfExpression. got=%T",
			stmt.Expression)
	}

	// ifの if(x < y ←　ここはInfixExpression){  return x }
	if !testInfixExpresison(t, exp.Condition, "x", "<", "y"){
		return
	}

	// 内部処理が一個だけってこと xは１つだけだからこれで正解
	if len(exp.Consequence.Statements) != 1{
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}

	// 内部処理に変換　←　ExprssionStatmentに(式文)
	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	// 変換が完了したか
	if !ok{
		t.Fatalf("Statements[0] is not ast.ExprssionStatement. got=%T",
					exp.Consequence.Statements[0])
	}

	// Identifierノードのチェック
	if !testIdentifier(t, consequence.Expression, "x"){
		return
	}

	// elseがない場合はAlternativeはnilになる
	if exp.Alternative != nil{
		t.Errorf("exp.Aliternative.Statements was not nil. got=%+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpresison(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("exp.Alternative.Statements does not contain 1 statements. got=%d\n",
			len(exp.Alternative.Statements))
	}

	// ifの時とかわらない
	// AlternativeにStatementsが査察てる
	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T){
	// x + y; はセミコロンなら中置型式のノード生成は終了するからこれでいい
	// EOFになって勝手に終了するので{}のよこに;をつけなくてもおｋ
	input := `fn(x, y){x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	// funだけで一本のノードになっている
	if len(program.Statements) != 1{
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	// いつもの式変換
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok{
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T",
			stmt.Expression)
	}
	
	// 関数のノードに変換
	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok{
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T",
			stmt.Expression)
	}

	// 引数の数をチェック
	if len(function.Parameters) != 2{
		t.Fatalf("function literal parameters wrong. want 2, got=%d\n",
			len(function.Parameters))
	}

	// Identの検証
	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	
	// x + y;の部分は式になるので１であってる
	if len(function.Body.Statements) != 1{
		t.Fatalf("function.Body.Statments has not 1 statements. got=%d\n",
			len(function.Body.Statements))
	}


	// ExpressionStatementノードに変換
	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok{
		t.Fatalf("function body stmt is not ast.ExprssionStatement. got=%T",
			function.Body.Statements[0])
	}

	fmt.Println(function.Body)
	fmt.Println(function.Body.Statements[0])

	// 	式チェック
	testInfixExpresison(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T){
	tests := []struct{
		input string
		expectedParams []string
	}{
		// {}のよこに;をつける場合、Expression_Statemntで;を飛ばしてくれているからOK
		// ついててもつけなくてもOK 関数は式だということを意識しよう！
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z){};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests{
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expectedParams){
			t.Errorf("length parameters wrong. want %d, got=%d\n",
				len(tt.expectedParams), len(function.Parameters))
		}
		for i, ident := range tt.expectedParams{
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T){
	// callerのテスト
	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1{
		t.Fatalf("Program.Statements does ot contain %d statements. got=%d\n",
					l, len(program.Statements))
	}

	// 式ノードに変換
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok{
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
			stmt.Expression)
	}

	// 関数呼び出しノードに変換
	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok{
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
			stmt.Expression)
	}

	// 関数名チェック
	if !testIdentifier(t, exp.Function, "add"){
		return
	}

	// argumentsチェック
	if len(exp.Arguments) != 3{
		t.Fatalf("Wrong length of arguments. got=%d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpresison(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpresison(t, exp.Arguments[2], 4, "+", 5)
}

func TestLetStatements(t *testing.T){
	tests := []struct{
		input string
		expectedIdentifier string
		expectedValue interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests{
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1{
			t.Fatalf("Program.Statements does not contatin 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier){
			return
		}

		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue){
			return
		}
	}
}

func TestStringLiteralExpression(t *testing.T){
	input := `"hello world";`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok{
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != "hello world"{
		t.Errorf("literal.Value not %q got=%q", "hello_world", literal.Value)
	}
}

func TestParsingArrayLiterals(t *testing.T){
	// [1, 4, 6]
	input := "[1, 2 * 2, 3 + 3]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok{
		t.Fatalf("exp not ast.ArrayLiteral. got = %T", stmt.Expression)
	}
	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpresison(t, array.Elements[1], 2, "*", 2)
	testInfixExpresison(t, array.Elements[2], 3, "+", 3)
}

func TestParsingIndexExpression(t *testing.T){
	input := "myArray[1 + 1]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok{
		t.Fatalf("exp not *ast.IndexExpression. got = %T", stmt.Expression)
	}
	if !testIdentifier(t, indexExp.Left, "myArray"){
		return
	}
	if !testInfixExpresison(t, indexExp.Index, 1, "+", 1){
		return 
	}
}

func TestParsingHashLiteralsStringKeys(t *testing.T){
	input := `{"one": 1, "two": 2, "three": 3}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok{
		t.Fatalf("exp is not ast.HashLiteral. got =%T", stmt.Expression)
	}

	if len(hash.Pairs) != 3{
		t.Errorf("hash.Pairs has wrong length. got =%d", len(hash.Pairs))
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for key, value := range hash.Pairs{
		literal, ok := key.(*ast.StringLiteral)
		if !ok{
			t.Errorf("Key is not ast.StringLiteral. got = %T", key)
		}
		expectedValue := expected[literal.String()]

		testIntegerLiteral(t, value, expectedValue)
	}
	
}

func TestParsingEmptyHashLiteral(t *testing.T){
	input := "{}"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok{
		t.Fatalf("exp is nt ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 0{
		t.Errorf("hash.Pairs has wrong length. got = %d", len(hash.Pairs))
	}
}

func TestParsingHashLiteralsWithExpressions(t *testing.T){
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`

	l := lexer.New(input)
	p := New(l)
	prgoram := p.ParseProgram()
	checkParseErrors(t, p)

	stmt := prgoram.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok{
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}
	if len(hash.Pairs) != 3{
		t.Errorf("hash.Pairs has wrong length. got=%T", len(hash.Pairs))
	}

	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression){
			testInfixExpresison(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression){
			testInfixExpresison(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression){
			testInfixExpresison(t, e, 15, "/", 5)	
		},
	}

	for key, value := range hash.Pairs{
		literaal, ok := key.(*ast.StringLiteral)
		if !ok{
			t.Errorf("key is not ast.StringLiteral. got = %T", key)
			continue
		}

		testFunc, ok := tests[literaal.String()]
		if !ok{
			t.Errorf("No test function for key %q found", literaal.String())
			continue
		}

		testFunc(value)
	}
}