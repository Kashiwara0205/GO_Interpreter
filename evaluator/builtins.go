package evaluator

import "../object"

// キーが組み込み関数の名前で、その中に突っ込むのがビルドイン関数のアドレス
var builtins = map[string]*object.Builtin{

	// 関数名：len
	// 内容: len("hoge")のように使用。　この場合４という整数を返却
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object{
		
			// 与えられたオブジェクト引数が複数あるならエラー１つだとOK
			if len(args) != 1{
				return newError("wrong number of arguments. got=%d, want=1",
				len(args))
			}

			// string以外エラー
			switch arg := args[0].(type){
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported, got %s",
					args[0].Type())
			}

		},

	},
}