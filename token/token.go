package token

// stringのalias
type TokenType string

// Token構造体
type Token struct{
	Type TokenType
	Literal string
}

// TokenType
const (
	ILLEGAL = "ILLEGAL"   // 未知な文字列・未知なトークン.
	EOF = "EOF"           // ファイル終端.

	// 識別子 + リテラル
	IDENT = "IDENT"
	INT = "INT"

	// 演算子
	ASSIGN = "="
	PLUS = "+"
	MINUS = "_"
	BANG = "!"
	ASTERISK = "*"
	SLASH = "/"
	EQ = "=="
	NOT_EQ = "!="

	LT = "<"
	GT = ">"

	// デリミタ
	COMMA = ","
	SEMICOLON = ";"

	// カッコ
	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"
	LBRACKET = "]"
	RBRACKET = "["

	// キーワード
	FUNCTION = "FUNCTION"
	LET = "LET"
	TRUE = "true"
	FALSE = "false"
	IF = "if"
	ELSE = "else"
	RETURN = "return"

	// 追加対応
	STRING = "STRING"
	COLON = ":"
	LOOP = "loop"
)

// キーワードハッシュ
// let foobarなど特別な意味をもつ文字列を、これを使って処理する
// TokenTypeとつながってる（ただのstringエイリアス）
var keywords = map[string]TokenType {
	"fn": FUNCTION,
	"let": LET,
	"true": TRUE,
	"false": FALSE,
	"if": IF,
	"else": ELSE,
	"return": RETURN,
	"loop": LOOP,
}

//　定義しておいた特別な意味をもつ文字列なのか、どうか検証する
func LookupIdent(ident string)TokenType{
	// keywordsに特別な意味の文字列があるかどうか
	// let とか foobarとか
	if tok, ok := keywords[ident]; ok{
		// あれば、その文字列を返却
		return tok
	}
	// なければIDENTを返却
	return IDENT
}