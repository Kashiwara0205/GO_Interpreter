package parser

import(
	"../ast"
	"../lexer"
	"../token"
	"fmt"
	"strconv"
)

//優先度-------------------------------------------------------------------
const (
	_int = iota            //　⇒　0始まり
	LOWSET
	EQUALS                 // ==
	LESSGREATER            // > または <
	SUM                    // +
	PRODUCT                // *
	PREFIX                 //  -X　または　!X
	CALL                   // myFunction(X)
	INDEX                  // array[index]
)
//-----------------------------------------------------------------------------

//エイリアス---------------------------------------------------------------------
type(
	// 前置型
	prefixParseFn func() ast.Expression
	// 中置型
	infixParseFn func(ast.Expression) ast.Expression
)
//------------------------------------------------------------------------------


//precedences---------------------------------------------------------------------
var precedences = map[token.TokenType]int{
	token.EQ: EQUALS,
	token.NOT_EQ: EQUALS,
	token.LT: LESSGREATER,
	token.GT: LESSGREATER,
	token.PLUS: SUM,
	token.MINUS: SUM,
	token.SLASH: PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN: CALL,
	token.LBRACKET: INDEX,
}
//------------------------------------------------------------------------------------



//Parserの定義------------------------------------------------------------------------
type Parser struct{
	// 字句解析インスタンスへのポインタ nextToken()等で使用
	l *lexer.Lexer
	// 現在調べているToken
	curToken token.Token
	//  先読みToken
	peekToken token.Token
	// エラー処理用
	errors []string

	// 構文解析関数
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns map[token.TokenType]infixParseFn
}
//--------------------------------------------------------------------------


/*解析関数以外の関数-------------------------------------------------------*/
func New(l *lexer.Lexer) *Parser{
	// パーサーに字句解析で使う構造体を仕込む
	p := &Parser{l: l,
				 errors: []string{},
				}

	//　前置型構文関数の初期化
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	// これと変わらん ⇒ 変数[token.IDENT] = p.parseIdentifier
	// この時点では、curTokenなどは空が入るので注意
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.LBRACE, p.parseHashLiteral)

	// 中置型構文関数の初期化
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)

	// 最初のpeekTokenには値が保持されていないため
	// ２回呼び出して
	// curToken =  null →　curToken = 値１
	// peekToken = 値１　→　peekToken = 値2 といった具合にする
	p.nextToken()
	p.nextToken()

	return p
}

// ExpressionインターフェースにIdentifierのアドレスを入れたものが戻り値
func (p *Parser) parseIdentifier() ast.Expression{
	// 今現在見ているトークンとリテラルをIdentifierノードとして返却
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// Error処理
func (p *Parser) Errors() []string{
	return p.errors
}

func (p *Parser) peekError(t token.TokenType){
	msg := fmt.Sprintf("expected next token no be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// 先読みTokenと今読みTokenを１つずつすすめる
func (p *Parser) nextToken(){
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool{
	return p.curToken.Type == t
}

// Token先読み
func (p *Parser) peekTokenIs(t token.TokenType) bool{
	return p.peekToken.Type == t
}

// 先読みして引数と同じタイプのTokenなら次のtokenを読み取った後にtrueを返して
// 違ったらfalse
func (p *Parser) expectPeek(t token.TokenType) bool{
	if p.peekTokenIs(t){
		// 同じトークンのタイプだとわかればトークンの位置をすすめる
		p.nextToken()
		return true
	}else{
		p.peekError(t)
		return false
	}
}

// 前置型の登録
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn){
	// メソッドがp.prefixParseFnsに入るんじゃねこれ
	// Tokenのタイプによって使うメソッドを切り分けてる
	p.prefixParseFns[tokenType] = fn
}

// 中置型の登録
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn){
	// メソッドがp.infixParseFnsに入るんじゃねこれ
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) noPrefixParseFnError(t token.TokenType){
	msg := fmt.Sprintf("no prefix parse fuction for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekPrecedence() int{
	if p, ok := precedences[p.peekToken.Type]; ok{
		return p
	}

	// 次の演算子がなかったら自動的に一番低いやつが返却
	return LOWSET
}

func (p *Parser) curPrecedence() int{
	//数値が帰ってくる
	//for example.  EQなら2が帰ってくる
	if p, ok := precedences[p.curToken.Type]; ok{
		return p
	}

	return LOWSET
}
/*---------------------------------------------------------------------------------------*/



/*ParseProgram解析----------------------------------------------------------------------*/
//何を解析するか：全ての解析関数の親元
//              astを生成

// astは抽象構文木
// ここにきた時点で字句解析は終わってる
func (p *Parser) ParseProgram() *ast.Program{
	//　抽象構文木のポインタを代入
	//  抽象構文木はProgram型構造体
	program := &ast.Program{}

	// statementに文字(let x  = 5)が入る
	// Statement型インターフェースを配列でprogram型構造体に仕込む
	// program型構造体はStatement型インターフェースを配列で定義済みだからできる
	program.Statements = []ast.Statement{}

	// 現在調べてるTokenがEOFじゃない限り回し続ける
	for p.curToken.Type != token.EOF{

		// letなのかifなのかfuncなのか　それに対する「木」が帰ってくる
		// letやifなどのキーワードが先頭にあった場合だけエラーを発報するようになっている
		// なので一度エラーがおきても引き続き連鎖で全部死ぬことはない
		stmt := p.parseStatement()
		if stmt != nil{
			// 木をがっちゃんこ
			program.Statements = append(program.Statements, stmt)
		}
		// 次のトークンへ
		p.nextToken()
	}

	// 抽象構文木返却
	return program
}

/*-----------------------------------------------------------------------*/



/*parseStatement解析---------------------------------------------------------*/
//解析内容：　トークンを見て、適切な解析関数に飛ばす

// let x 5みたいにミスったらnilが返却され続けて次のletまでいく
// tokenを解析しツリー：Statementを返却
func (p *Parser) parseStatement() ast.Statement{
	// Tokenが何なのか
	switch p.curToken.Type{
	// let 定義文　let x = 5のように使用
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	// 式
	default:
		return p.parseExpressionStatement()
	}
}
/*-----------------------------------------------------------------------------*/



/*ParseLetStatement解析---------------------------------------------------------*/
//何を解析するか： let a; let o = 0; などのlet文に関して解析する

// letをparseしてツリー：letStatmentを返却
func (p *Parser) parseLetStatement() *ast.LetStatement{
	// LetStatementツリーを作成
	stmt := &ast.LetStatement{Token: p.curToken}

	// letトークンの次がIDENT(変数系)じゃなかったらnil
	if !p.expectPeek(token.IDENT){
		return nil
	}

	// LetStatmentツリーの名前の部分
	// literalなので、letとか５とか+とかそんなんが入る
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// IDENTから１つ先読みしてASSIGNかどうか ASSIGNは=のこと
	if !p.expectPeek(token.ASSIGN){
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWSET)

	if p.peekTokenIs(token.SEMICOLON){
		p.nextToken()
	}

	// LetStatmentツリー返却
	return stmt
}

/*-----------------------------------------------------------------------------------*/



/*parseReturnStatement解析---------------------------------------------------------*/
//何を解析するか： retun hoge;などのretun文を解析する

// returnをパース
func (p *Parser) parseReturnStatement() *ast.ReturnStatement{
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWSET)

	for !p.curTokenIs(token.SEMICOLON){
		p.nextToken()
	}
	return stmt
}
/*---------------------------------------------------------------------------------*/



/*parseExpressionStatement解析--------------------------------------------------*/
// 何を解析するか: 式文の解析　hoge; low; みたいな

// Expression_Statemntノードの作成
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement{

	// ExpressionStatementノード作成
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	// 優先度を設定して新しいExpression読み込み
	// Identのアドレスが格納されたExpressionが飛んできたり色々
	stmt.Expression = p.parseExpression(LOWSET)

	if p.peekTokenIs(token.SEMICOLON){
		p.nextToken()
	}

	return stmt
}
/*-----------------------------------------------------------------------------------*/



/*parseExpression解析--------------------------------------------------------------------*/
//何を解析するか: 「式」を全て解析する
// 全てのparse系の親元に値する
// 　1の場合　⇒　parseIntegerLiteralへ
//  -1の場合　⇒　prefixExpressionへ

func (p *Parser) parseExpression(precedence int) ast.Expression{
	// 一発目に入ってくる場合はLOWSET

	// 今見ている前置型トークンから適するメソッドを引っ張りだす
	// IDENTの場合ならparseIdentifier()とか
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil{
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	// 新たにExpressionノードを作成
	// ここにExpresion型の木がささる　Expresion型の枝には、identfierのアドレス integerliteralのアドレス等
	// 左括弧もここから
	leftExp := prefix()

	// 次点がトークンでない、　かつ　次のトークンの優先度より下

	// 括弧の偉大なトリックの正体は、)のLOWSET比較と変数か数値のLOWESET比較で　100(LOWSET) < )(LOWSET)
	// 下記の状況を作り出すこと
	// 1 + 1とかなら [1] LOWSET <  [ なし ]LOWSETで終了する
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence(){
		// 中置型取得
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil{
			return leftExp
		}

		p.nextToken()

		// 左のやつはここで入る
		// 左のやつは前置型解析関数にかけられ済みのintafeace方であるleftExp
		leftExp = infix(leftExp)
	}

	return leftExp
}
/*---------------------------------------------------------------------------------*/



/*parseIntegerLiteral解析---------------------------------------------------------*/
//何を解析するか: 1, 200などの数値式

func (p *Parser) parseIntegerLiteral() ast.Expression{
	lit := &ast.IntegerLiteral{Token: p.curToken}

	// 文字列　⇒　値　
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil{
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

/*----------------------------------------------------------------------------------------*/



/*prefixExpression解析-------------------------------------------------------------------*/
//何を解析するか: -1, !trueなど前に-や!がつくもの

func (p *Parser) parsePrefixExpression() ast.Expression{
	// PrefixExpressionノードの作成
	expression := &ast.PrefixExpression{
		Token: p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()

	// Identとかintegerliteralとかのアドレスを格納したExpressionが刺さる
	expression.Right = p.parseExpression(PREFIX)

	return expression
}
/*---------------------------------------------------------------------------------------*/


/*parseInfixExpression解析---------------------------------------------------------------*/
//何を解析するか: 中置型の式を解析
//               1 + 1, 4 * 5
// !!!

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression{
	
	expression := &ast.InfixExpression{
		Token: p.curToken,
		Operator: p.curToken.Literal,
		Left: left,
	}

	// 優先度の抽出
	precedence := p.curPrecedence()

	//１回めの比較　⇒ LOWSET 演算子
	//2回めの比較　⇒ 演算子 演算子

	// ここをすすめることにより　式と式の比較がpaeseなんちゃらで可能になる
	p.nextToken()

	// 右も左と同じようにinteface型のものになる
	// ここから多段階構造になっていく

	// a + b * c この場合だと下記のようになる
	// [ a ] + [ b ]
	//		   ([ b ] * [ c ])　←　bが帰ってくるのではなくこいつが帰ってくる
	expression.Right = p.parseExpression(precedence)

	return expression
}
/*---------------------------------------------------------------------------------------*/


/*parseBoolean解析-----------------------------------------------------------------------*/
//何を解析するか: 真偽値
func (p *Parser) parseBoolean() ast.Expression{
	// curTokenISで現在トークンがtrueならtrueが
	//             現在トークンがfalseならfalseが入る

	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}
/*---------------------------------------------------------------------------------------*/


/*parseGroupedExpression解析--------------------------------------------------------*/
//何を解析するか: グループ式 (1 + 1)* 4とか
// 優先度の変更をしている

func (p *Parser)  parseGroupedExpression() ast.Expression{
	// [ ( ] [ a ]の aに移動する
	p.nextToken()

	// 優先度の変更をして再出発
	exp := p.parseExpression(LOWSET)

	// 右括弧がなかったら返却
	if !p.expectPeek(token.RPAREN){
		return nil
	}

	return exp
}
/*-------------------------------------------------------------------------------*/

/*IfExpression解析---------------------------------------------------------------*/
func (p *Parser) parseIfExpression() ast.Expression{
	// IfExpressionの生成
	expression := &ast.IfExpression{Token: p.curToken}

	// LPARENは[ ( ]のとこ
	// 次のトークンが[ ( ]かどうか
	if !p.expectPeek(token.LPAREN){
		return nil
	}

	// [ 式 ] ←こいつになった
	p.nextToken()

	// if (x < y)とかの場合は、InfixExpressionが帰ってくる
	// if(varidate_result(p))などの関数が混入する場合は、わからない
	expression.Condition = p.parseExpression(LOWSET)

	// 2パターン
	// if(式){処理}else{処理}
	// if（式）{処理}

	// [if] ( [x < y] [ ) ]  ←　RAPAENは左括弧
	if !p.expectPeek(token.RPAREN){
		return nil
	}
	// LBRACE ← [ { ]
	if !p.expectPeek(token.LBRACE){
		return nil
	}

	// if (式) ←いまここまでの解析が終わった
	// StatementノードというでっかいノードがConsequesnceに刺さる
	expression.Consequence = p.parseBlockStatement()

	// 次のトークンがELSE式だったら続行
	if p.peekTokenIs(token.ELSE){
		p.nextToken()

		if !p.expectPeek(token.LBRACE){
			return nil
		}
		// Alternativeはelseの処理を管理
		expression.Alternative = p.parseBlockStatement()
	}
	return expression
}
/*-------------------------------------------------------------------------------*/

/*parseBlockStatement解析-------------------------------------------------------*/
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	// ブロックステートメントノードを生成
	block := &ast.BlockStatement{Token: p.curToken}
	// ステートメントノードを持つ
	block.Statements = []ast.Statement{}
	// [ { ] ←　こいつの次だから処理に以降する
	p.nextToken()

	// [ } ] ←　これが来るか　ファイルの終わりまで読み取り
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF){
		// 新たにノードを構築する
		stmt := p.parseStatement()
		if stmt != nil{
			// ノードを追加していく
			block.Statements = append(block.Statements, stmt)
		}
		// 終了したら次のトークンへ
		p.nextToken()
	}

	return block
}
/*------------------------------------------------------------------------------*/

/*parseFunctionLiteral解析-------------------------------------------------------*/
func (p *Parser) parseFunctionLiteral() ast.Expression{
	// FunctiioLiteralノードの作成
	lit := &ast.FunctionLiteral{Token: p.curToken}

	// 左括弧チェック
	if !p.expectPeek(token.LPAREN){
		return nil
	}

	// 引数を入手
	// ここに来た時点でParameterは完成
	lit.Parameters = p.parseFunctionParameters()

	// 右括弧チェックは、parseFunctionParametersで終わっているので[ { ]からのチェックになる
	if !p.expectPeek(token.LBRACE){
		return nil
	}

	// ごついStatementノードが刺さる
	lit.Body = p.parseBlockStatement()
	
	return lit
}
/*----------------------------------------------------------------------------------*/

/*parseFunctionParameters解析-------------------------------------------------------*/
func (p *Parser) parseFunctionParameters() []*ast.Identifier{
	// identifier(変数ノード)を配列で作成
	identifiers := []*ast.Identifier{}

	// 次のトークンが右括弧だったら引数なし関数
	if p.peekTokenIs(token.RPAREN){
		// トークンを勧めておく
		p.nextToken()

		// parameter返却
		// 何も入ってないから空のidentifiersが返却される
		return identifiers
	}

	// ([ x ], y )
	p.nextToken()

	// 変数束縛ノード作成
	// p.curToken ⇒ x
	// Value ⇒　x
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	// 変数束縛ノード配列に追加
	identifiers = append(identifiers, ident)

	// 次がカンマなら実行
	for p.peekTokenIs(token.COMMA){
		// x[,] yカンマ飛ばし
		p.nextToken()
		// x, [y]お目当ての変数へ
		p.nextToken()
		// 上と同じような形式で変数束縛ノードを作成して配列に突っ込んでいく
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}
	// 右括弧チェック
	if !p.expectPeek(token.RPAREN){
		return nil
	}

	// 返却
	return identifiers
}
/*---------------------------------------------------------------------------------*/

/*parseCallExpression解析-----------------------------------------------------------*/
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression{
	// CallExpressioノード作成
	// Expressionの下にIDENTが刺さってる
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	// 引数をもらう
	exp.Arguments = p.parseExpressionList(token.RPAREN)

	// Caller返却
	return exp
}
/*----------------------------------------------------------------------------------*/

/*parseCallArguments解析-----------------------------------------------------------*/
func (p *Parser) parseCallArguments() []ast.Expression{
	// 式ノード配列を生成
	args := []ast.Expression{}

	// 引数なしの関数
	if p.peekTokenIs(token.RPAREN){
		p.nextToken()
		return args
	}

	p.nextToken()
	// Expressionが刺さる
	args = append(args, p.parseExpression(LOWSET))

	//　カンマは飛ばして刺していく
	for p.peekTokenIs(token.COMMA){
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWSET))
	}

	// 右括弧チェック	
	if !p.expectPeek(token.RPAREN){
		return nil
	}

	// 引数返却
	return args
}

/*----------------------------------------------------------------------------------*/

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{ Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseArrayLiteral() ast.Expression{
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Elements = p.parseExpressionList(token.RBRACKET)

	return array
}

func(p *Parser) parseExpressionList(end token.TokenType) []ast.Expression{
	// 値なら何でも刺さるやつを生成
	// ast.Expression型の配列を生成
	list := []ast.Expression{}

	// 閉じ括弧の出現
	if p.peekTokenIs(end){
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWSET))

	for p.peekTokenIs(token.COMMA){
		// カンマ分飛ばし
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWSET))
	}

	// 閉じ括弧がなかったらnil返却
	if !p.expectPeek(end){
		return nil
	}

	return list
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWSET)

	if !p.expectPeek(token.RBRACKET){
		return nil
	}

	return exp
}

func (p *Parser) parseHashLiteral() ast.Expression{
	hash := &ast.HashLiteral{Token: p.curToken}
	// makeはhashを作ることもできる
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(token.RBRACE){
		p.nextToken()
		key := p.parseExpression(LOWSET)

		if !p.expectPeek(token.COLON){
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWSET)

		hash.Pairs[key] = value
		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA){
			return nil
		}
	}

	if !p.expectPeek(token.RBRACE){
		return nil
	}

	return hash
}