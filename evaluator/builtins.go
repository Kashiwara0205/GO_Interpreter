package evaluator

import (
	"../object"
	"fmt"
)

// キーが組み込み関数の名前で、その中に突っ込むのがビルドイン関数のアドレス
var builtins = map[string]*object.Builtin{

	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object{

			// 与えられたオブジェクト引数が複数あるならエラー１つだとOK
			if len(args) != 1{
				return newError("wrong number of arguments. got=%d, want=1",
				len(args))
			}

			// string以外エラー
			switch arg := args[0].(type){
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported, got %s",
					args[0].Type())
			}

		},
	},

	"first": &object.Builtin{
		Fn: func(args ...object.Object) object.Object{
			if len(args) != 1{
				return newError("wrong number of arguments. got = %d, want = 1",
								len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `first` must be ARRAY, got %s",
								args[0].Type())
			}

			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}


			return NULL
		},
	},

	"last" : &object.Builtin{
		Fn: func(args ...object.Object) object.Object{
			if len(args) != 1{
				return newError("wrong number of arguments. got = %d, want = 1",
								len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `first` must be ARRAY, got %s",
								args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if len(arr.Elements) > 0 {
				return arr.Elements[length - 1]
			}

			return NULL
		},
	},

	"reset": &object.Builtin{
		Fn: func(args ...object.Object) object.Object{
			if len(args) != 1{
				return newError("wrong number of arguments. got=%d, want=1",
							len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ{
				return newError("argument to `reset` must be ARRAY, got = %s",
								args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0{
				newElements := make([]object.Object, length-1, length-1)
				// コピー元はいじらない
				// 必要な文だけコピーしてそれを返す
				copy(newElements, arr.Elements[1:length])

				return &object.Array{Elements: newElements}
			}

			return NULL
		},
	},

	"push": &object.Builtin{
		Fn: func(args ...object.Object) object.Object{
			if len(args) != 2{
				return newError("wrong number of arguments. got = %d, want = 2",
								len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ{
				return newError("argument to `push` must be ARRAY, got %s",
								args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)

			// 第一引数：型　第二引数：長さ　第三引数：キャパシティ（容量）
			newElements := make([]object.Object, length + 1, length + 1)
			// 元の配列をコピー
			copy(newElements, arr.Elements)
			// +1した容量に、引数を追加
			newElements[length] = args[1]

			return &object.Array{Elements: newElements}
		},
	},

	"puts": &object.Builtin{
		Fn: func(args ...object.Object) object.Object{
			for _, arg := range args{
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},
}