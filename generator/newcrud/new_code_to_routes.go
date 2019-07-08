package newcrud

import (
	"errors"
)

const (
	expectingROfRoutes = 4
	expectingOOfRoutes = 5
	expectingUOfRoutes = 6
	expectingTOfRoutes = 7
	expectingEOfRoutes = 8
	expectingSOfRoutes = 9

	expectingROfReturn  = 14
	expectingEOfReturn  = 15
	expectingTOfReturn  = 16
	expectingUOfReturn  = 17
	expectingR2OfReturn = 18
	expectingNOfReturn  = 19

	returnRoutesWasFound = 20
)

var returnRouteAt int
var funcRoutesWasFound bool

func insideFuncRoutes(i int) {
	if currentChar == "{" && (len(stack) == 0 || (len(stack) > 0 && stack[len(stack)-1] == "{")) {
		stack = append(stack, currentChar)
	} else if currentChar == "}" && stack[len(stack)-1] == "{" {
		stack = stack[:len(stack)-1]
	} else if currentChar == "`" && stack[len(stack)-1] == "`" && !lastCharWasBackSlash {
		stack = stack[:len(stack)-1]
	} else if currentChar == "`" && stack[len(stack)-1] == "{" {
		stack = append(stack, "`")
	} else if currentChar == "\"" && stack[len(stack)-1] == "\"" && !lastCharWasBackSlash {
		stack = stack[:len(stack)-1]
	} else if currentChar == "\"" && stack[len(stack)-1] == "{" {
		stack = append(stack, "\"")
	} else if currentChar == "\"" && stack[len(stack)-1] == "{" {
		stack = append(stack, "\"")
	} else if currentChar == "/" && lastCharWasSlash && (len(stack) == 0 || stack[len(stack)-1] == "{") {
		currentState = commentSingleLine
	} else if currentChar == "*" && lastCharWasSlash && (len(stack) == 0 || stack[len(stack)-1] == "{") {
		currentState = commentMultiLine
	} else if currentChar == "r" && currentState == expectingROfReturn {
		currentState = expectingEOfReturn
		returnRouteAt = i
	} else if currentChar == "e" && currentState == expectingEOfReturn {
		currentState = expectingTOfReturn
	} else if currentChar == "t" && currentState == expectingTOfReturn {
		currentState = expectingUOfReturn
	} else if currentChar == "u" && currentState == expectingUOfReturn {
		currentState = expectingR2OfReturn
	} else if currentChar == "r" && currentState == expectingR2OfReturn {
		currentState = expectingNOfReturn
	} else if currentChar == "n" && currentState == expectingNOfReturn {
		currentState = returnRoutesWasFound
	} else {
		currentState = expectingROfReturn
	}

	lastCharWasBackSlash = (currentChar == "\\")
	lastCharWasSlash = (currentChar == "/")
}

func AppendNewCodeToRoutes(newCode string, code string) (string, error) {
	var err error
	var i int

	funcRoutesWasFound = false
	currentState = expectingFOfFunc
	stack = nil

	for i = 0; i < len(code); i++ {
		currentChar = code[i : i+1]
		if regexpEmptyChar.MatchString(currentChar) && currentState != commentSingleLine {
			continue
		} else if currentChar == "f" && currentState == expectingFOfFunc {
			currentState = expectingUOfFunc
		} else if currentChar == "u" && currentState == expectingUOfFunc {
			currentState = expectingNOfFunc
		} else if currentChar == "n" && currentState == expectingNOfFunc {
			currentState = expectingCOfFunc
		} else if currentChar == "c" && currentState == expectingCOfFunc {
			currentState = expectingROfRoutes
		} else if currentChar == "R" && currentState == expectingROfRoutes {
			currentState = expectingOOfRoutes
		} else if currentChar == "o" && currentState == expectingOOfRoutes {
			currentState = expectingUOfRoutes
		} else if currentChar == "u" && currentState == expectingUOfRoutes {
			currentState = expectingTOfRoutes
		} else if currentChar == "t" && currentState == expectingTOfRoutes {
			currentState = expectingEOfRoutes
		} else if currentChar == "e" && currentState == expectingEOfRoutes {
			currentState = expectingSOfRoutes
		} else if currentChar == "s" && currentState == expectingSOfRoutes {
			funcRoutesWasFound = true
			currentState = expectingROfReturn
		} else if currentChar == "/" && lastCharWasSlash && currentState != commentMultiLine {
			stateBeforeComment = currentState
			currentState = commentSingleLine
		} else if currentChar == "*" && lastCharWasSlash && currentState != commentSingleLine {
			stateBeforeComment = currentState
			currentState = commentMultiLine
		} else if funcRoutesWasFound && currentState != returnRoutesWasFound {
			insideFuncRoutes(i)
		} else if funcRoutesWasFound && currentState == returnRoutesWasFound {
			break
		} else {
			currentState = expectingFOfFunc
		}
	}

	if currentState == returnRoutesWasFound && len(stack) > 0 {
		outputStr = code[0:returnRouteAt-1] + newCode + "\n\n" + code[returnRouteAt-1:len(code)]
	} else {
		err = errors.New("Could not parse Migrate file. Please, check sintaxe.")
	}

	return outputStr, err
}
