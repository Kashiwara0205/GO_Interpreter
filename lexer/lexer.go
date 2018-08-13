package lexer

import(
	"../token"
)

// 構造体は、何もいれないと、それぞれ初期値が入る。
type Lexer struct{
	// バイト型スライスと考える　Monkye言語のソースがすべて入る
	input string
	// 現在の読み取り位置
	position int
	// 次に読み取る位置
	readPosition int
	// 現在の調査文字
	ch byte
}

// ポインタ使ってるから値は上書き
func (l *Lexer) readChar() {
	// 文字数の数 ＝　バイト数としているlenなのでASCLLのみに対応
	if l.readPosition >= len(l.input) {
		// 終わりに達したら
		l.ch = 0
	} else {
		// 解析内容をbyteで読み込み（１文字読み込み）
		l.ch = l.input[l.readPosition]
	}
	// 現在の読み取り位置を格納
	l.position = l.readPosition
	// 次の読み取り内容へ
	l.readPosition += 1
}

// Lexter構造体を新しく作成
func New(input string) *Lexer{
	// Lexter構造体に分析する文字列を格納し、そのアドレスをlに入れる。
	l := &Lexer{input: input}
	// 最初の文字を読み込む
	l.readChar()
	return l
}

// Lexer構造体の関数
func (l *Lexer) NextToken() token.Token {
	// Toke構造体を定義
	var tok token.Token

	// let fiveのようなinputであった場合、[let][][five]となってしまい[]がイリーガルエラーとなるので
	// スペースの場合は、読み飛ばす
	l.skipWhitespace()

	// 読み取った文字をswitchにかける
	switch l.ch{
	case '=':
		// 2文字タイプの演算子、例えば　== なら[=][=]となり=が２回続いていると判断されるので
		// それをくいとめる
		// ==
		if l.peekChar() == '='{
			ch := l.ch
			// １イコール分をよみとっておく
			l.readChar()
			// リテラルとして==を生成
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		}else{
		// 新しいtokenを生成
		tok = newToken(token.ASSIGN, l.ch)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		// != ２文字タイプ判定
		if l.peekChar() == '='{
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		}else{
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case 0:
		// ここだけ書き方違うけど多分 l.chでかけないから崩してるだけで
		// やってることは同じ 空白　＝　１トークン
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch){
			// max_sum などの変数？がliteralに代入される
			tok.Literal = l.readIdentifiter()
			// let や　funcなどの特別な文字列かどうか
			// ここのTypeには特別な文字列の場合　→ LET FUNCTIONとか入る
			// 普通だったらIDENTが入る
			tok.Type = token.LookupIdent(tok.Literal)
			// リテラルを返却
			return tok
		}else if isDigit(l.ch){
			// 数値であらわすことを定義するINTをTypeに代入
			tok.Type = token.INT
			// ここでは　'542'とか '0'が帰ってくる
			// だいじなのは数値でなく数字が帰ってくることである
			tok.Literal = l.readNumber()

			return tok
		}else{
			// **や==~など認識されない文字列が出現した場合、例外発生
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	// 次の文字へ
	l.readChar()
	return tok
}

// 空白や改行をスキップさせる
func (l *Lexer) skipWhitespace(){
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r'{
		l.readChar()
	}
}

func (l *Lexer) readIdentifiter() string{
	position := l.position
	// キーワード、もしくは認識されていない文字が出現するまで回る
	for isLetter(l.ch){
		// 次の文字へ
		l.readChar()
	}
	

	// 読み取った部分をスライスで切り取ってreturn
	//　例　max_num, hogeなど
	return l.input[position:l.position]
}

// IDENTや特別な文字列（キーワード）かを判別するときに使う
func isLetter(ch byte) bool{
	// a-z A-Z _ ならtrue
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// 数字の読み取り
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch){
		l.readChar()
	}

	// 数値を返却
	// ここで返却されるのは数字であって数値ではない
	return l.input[position:l.position]
}

// 数字かどうか
func isDigit(ch byte) bool{
	return '0' <= ch && ch <= '9'
}

// 新しいトークンを作成する 1文字　＝　1トークン
func newToken(tokenType token.TokenType, ch byte) token.Token{
	// token構造体を作成
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// 一つ先を読み取る
// peek(のぞき見)
func (l *Lexer) peekChar() byte{
	// 読み取り内容が次で終わってたのであれば、0を返却
	if l.readPosition >= len(l.input){
		return 0
	}else{
		// 次の読み取り位置を返却
		return l.input[l.readPosition]
	}
}

func (l *Lexer) readString() string{
	// " "で区切られた箇所をトークン化
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0{
			break
		}
	}
	return l.input[position:l.position]
}