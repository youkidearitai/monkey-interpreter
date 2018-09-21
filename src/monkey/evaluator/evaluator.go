package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		//fmt.Printf("Program\n")
		return evalProgram(node)
	case *ast.BlockStatement:
		//fmt.Printf("BlockStatement\n")
		return evalBlockStatement(node)
	case *ast.ExpressionStatement:
		//fmt.Printf("ExpressionStatement\n")
		return Eval(node.Expression)
	case *ast.PrefixExpression:
		//fmt.Printf("PrefixExpression\n")
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue)
		//fmt.Printf("ReturnStatement: %d\n", val)
		return &object.ReturnValue{Value: val}
	case *ast.InfixExpression:
		//fmt.Printf("InfixExpression\n")
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		//fmt.Printf("IfExpression\n")
		return evalIfExpression(node)
	case *ast.IntegerLiteral:
		//fmt.Printf("IntegerLiteral: %d\n", node.Value)
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		//fmt.Printf("Boolean: %t\n", node.Value)
		return nativeBoolToBooleanObject(node.Value)
	}
	//fmt.Println("Eval: nil Error")
	return nil
}

func evalProgram(program *ast.Program) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement)

		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue.Value
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object
	/*
		//fmt.Printf("evalBlockStatement: %+v\n", block)
		//fmt.Printf("evalBlockStatement: %+v\n", block.Statements)
	*/

	for _, statement := range block.Statements {
		result = Eval(statement)

		////fmt.Printf("inloop: evalBlockStatement: %+v\n", result)
		if result != nil && result.Type() == object.RETURN_VALUE_OBJ {
			return result
		}
	}

	////fmt.Printf("evalBlockStatement: %+v\n", result)
	return result
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement)

		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue.Value
		}
	}

	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}

	return FALSE
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return NULL
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return NULL
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	//fmt.Println("evalInfixExpression")
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		//fmt.Println("Eval: evalIntegerInfixExpression")
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		//fmt.Println("NULL")
		return NULL
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
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
		//fmt.Println("Eval: evalIntegerInfixExpression NULL")
		return NULL
	}
}

func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)

	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	} else {
		//fmt.Printf("evalIfExpression: %+v\n", condition)
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
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
