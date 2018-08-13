package evaluator

import(
	"../ast"
	"../object"
	"fmt"
)

var (
	NULL = &object.Null{}
	TRUE = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// 5 という式の場合こういう順序
// program ⇒　exoression ⇒　integerLiteral ⇒　program　⇒ return
func Eval(node ast.Node, env *object.Enviroment) object.Object{
	// 変換じゃなくて型のswitchをするのが.(type)
	switch node := node.(type){
	case *ast.Program:
		// プログラムノードから、このswitch文に連続で入る
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		// Expression[->] IntegerLiteralなら次はIntegerLiteral型となる
		return Eval(node.Expression, env)

	case *ast.IntegerLiteral:
		// astに格納されているvalueを入れてIntegerオブジェクト返却
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		// !!などの場合再帰呼出ししててわかりにくいから注意
		right := Eval(node.Right, env)
		if isError(right){
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		// 再帰的に回すために、例えば1 + 1 + 1とかなら
		//　最終的にleftには2 rightには1が入る
		left := Eval(node.Left, env)
		if isError(left){
			return left
		}
		right := Eval(node.Right, env)
		if isError(right){
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	case *ast.BlockStatement:
		return evalBlockStament(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val){
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val){
			return val
		}
		env.Set(node.Name.Value, val)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}
	
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function){
			return function
		}
		// argumentsはリストになって帰ってくる
		// callerの引数
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]){
			return args[0]
		}
		
		return applyFunction(function, args)

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)

		// 要素が１でエラー吐いてたらそのまま返却
		if len(elements) == 1 && isError(elements[0]){
			return elements[0]
		}
		return &object.Array{Elements: elements}
	
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left){
			return left
		}
		index := Eval(node.Index, env)
		if isError(index){
			return index
		}
		return evalIndexExpression(left, index)
	}

	return nil
}

// ボコボコ生み出すのではなく一つのtrueとfalseを使い回す
func nativeBoolToBooleanObject(input bool) *object.Boolean{
	if input{
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(operator string, right object.Object) object.Object{
	switch operator{
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object{
	switch right{
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	// !5とかがfalseになるのはこいつのおかげ
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object{
	// 右がintじゃなかったらNULL返却 ⇒　Error返却に変更
	if right.Type() != object.INTEGER_OBJ{
		return newError("unknown operator: -%s", right.Type())
	}
	// すでに作成されているInt型オブジェクトから値を取り出し
	value := right.(*object.Integer).Value
	// -にしてオブジェクトに格納
	return &object.Integer{Value: -value}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object{
	switch{
	// 式の両方が整数だったら
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		// 計算結果をオブジェクトに格納して返却
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s",
				left.Type(), operator, right.Type())
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	default :
		return newError("unknown operator: %s %s %s",
				left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object{
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator{
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s",
				left.Type(), operator, right.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Enviroment) object.Object{
	// 正常な場合ならInfixExpression（式）が評価されて返却される
	// true or false
	condition := Eval(ie.Condition, env)
	if isError(condition){
		return condition
	}

	if isTruthy(condition){
		return Eval(ie.Consequence, env)
	}else if ie.Alternative != nil{
		return Eval(ie.Alternative, env)
	}else{
		return NULL
	}
}

func isTruthy(obj object.Object) bool{
	switch obj{
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

// 大元
func evalProgram(program *ast.Program, env *object.Enviroment) object.Object{
	var result object.Object
	for _, statement := range program.Statements{
		result = Eval(statement, env)

		switch result := result.(type){
		// return返却（ifの中とかfuncの中とかもここに入って終わる）
		case *object.ReturnValue:
			return result.Value
		// エラー返却
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStament(block *ast.BlockStatement, env *object.Enviroment) object.Object {
	var result object.Object

	for _, statement := range block.Statements{
		result = Eval(statement, env)

		if result != nil{
			rt := result.Type()
			// returnかerrorが来たらブロックの読み取り終了
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ{
				return result
			}
		}
	}

	return result
}

// aは呼び出されなかった時にnilとなる
// a := []interface{}{"hello", "world", 42}みたいな感じで格納
func newError(format string, a ...interface{}) *object.Error{
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil{
		// エラーオブジェクトが帰ってきていたらtrueを返却
		return obj.Type() == object.ERROR_OBJ
	}
	// オブジェクトが存在しない場合にどうするか
	return false
}

func evalIdentifier(node *ast.Identifier, env *object.Enviroment) object.Object{
	// x = 100
	// イメージ ⇒　["x"]に100が入ってる
	// １００　 ←　["x"]環境下から取り出し
	//　変数束縛しているアドレスとかも入るので注意

	// letなんて変数定義したら後ろにあるbuitinが動かなくなるので注意
	if val, ok := env.Get(node.Value); ok{
		return val
	}

	// 組み込み関数の呼び出しはここで行われる
	// buitinオブジェクトとして返却
	if builtin, ok := builtins[node.Value]; ok{
		return builtin
	}

	// 取り出してきた値を返却
	return newError("identifier not foutnd:" + node.Value)
}

func evalExpressions(
	exps []ast.Expression,
	env *object.Enviroment,
) []object.Object{
	var result []object.Object

	// expsに引数が入ってる
	for _, e := range exps{
		evaluted := Eval(e, env)
		if isError(evaluted){
			return []object.Object{evaluted}
		}
		result = append(result, evaluted)
	}

	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object{
	switch fn := fn.(type){

	// 普通の関数
	case *object.Function:
		// スコープ変更
		extendEnv := extendFunctionEnv(fn, args)
		// 関数の中身を処理
		evaluted := Eval(fn.Body, extendEnv)
		return unwrapReturnValue(evaluted)

	// 組み込み関数は、スコープの変更はいらない
	// 全然違うところでやるから
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(
	fn *object.Function,
	args []object.Object,
) *object.Enviroment{
	// 現在使用しているスコープをouterに設定して関数内で使うenviromentが帰ってくる
	env := object.NewEncloseEnviroment(fn.Env)

	// fn.Parametersが関数のやつ　fn hoge(x)とかのx
	// argsが引数！ hoge(10)とかの10
	for paramIdx, param := range fn.Parameters{
		env.Set(param.Value, args[paramIdx])
	}
	return env
}

func unwrapReturnValue(obj object.Object) object.Object{
	if returnValue, ok := obj.(*object.ReturnValue); ok{
		return returnValue.Value
	}
	return obj
}

func evalStringInfixExpression(
	operator string,
	left, right object.Object,
)object.Object{
	if operator != "+"{
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	return &object.String{Value: leftVal + rightVal}
}

func evalIndexExpression(left, index object.Object) object.Object{
	switch{
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object{
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return arrayObject.Elements[idx]
}