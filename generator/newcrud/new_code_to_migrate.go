package newcrud

import (
	"errors"
	"regexp"
)

const (
	expectingFOfFunc = 0
	expectingUOfFunc = 1
	expectingNOfFunc = 2
	expectingCOfFunc = 3

	expectingMOfMigrate = 4
	expectingIOfMigrate = 5
	expectingGOfMigrate = 6
	expectingROfMigrate = 7
	expectingAOfMigrate = 8
	expectingTOfMigrate = 9
	expectingEOfMigrate = 10

	funcMigrateWasFound = 11

	commentSingleLine = 12
	commentMultiLine  = 13
)

var lastCloseBracket int
var currentState int
var currentChar string
var regexpEmptyChar = regexp.MustCompile(`[\s\t\n\r\f]`)
var stack []string
var outputStr string
var lastCharWasBackSlash bool
var lastCharWasSlash bool
var lastCharWasStar bool
var stateBeforeComment int

func insideFuncMigrate(i int) {
	if currentChar == "{" && (len(stack) == 0 || (len(stack) > 0 && stack[len(stack)-1] == "{")) {
		stack = append(stack, currentChar)
	} else if currentChar == "}" && stack[len(stack)-1] == "{" {
		stack = stack[:len(stack)-1]
		lastCloseBracket = i + 1
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
	}
}

func insideCommentSingleLine() {
	lastCharWasSlash = false
	if currentChar == "\n" {
		currentState = stateBeforeComment
	}
}

func insideCommentMultiLine() {
	lastCharWasSlash = false
	if currentChar == "/" && lastCharWasStar {
		currentState = stateBeforeComment
	}

	lastCharWasStar = currentChar == "*"
}

func AppendNewCodeToMigrate(newLine string, code string) (string, error) {
	var err error
	var i int

	lastCharWasSlash = false
	lastCharWasSlash = false
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
			currentState = expectingMOfMigrate
		} else if currentChar == "M" && currentState == expectingMOfMigrate {
			currentState = expectingIOfMigrate
		} else if currentChar == "i" && currentState == expectingIOfMigrate {
			currentState = expectingGOfMigrate
		} else if currentChar == "g" && currentState == expectingGOfMigrate {
			currentState = expectingROfMigrate
		} else if currentChar == "r" && currentState == expectingROfMigrate {
			currentState = expectingAOfMigrate
		} else if currentChar == "a" && currentState == expectingAOfMigrate {
			currentState = expectingTOfMigrate
		} else if currentChar == "t" && currentState == expectingTOfMigrate {
			currentState = expectingEOfMigrate
		} else if currentChar == "e" && currentState == expectingEOfMigrate {
			currentState = funcMigrateWasFound
		} else if currentChar == "/" && lastCharWasSlash && currentState != commentMultiLine {
			stateBeforeComment = currentState
			currentState = commentSingleLine
		} else if currentChar == "*" && lastCharWasSlash && currentState != commentSingleLine {
			stateBeforeComment = currentState
			currentState = commentMultiLine
		} else if currentState == funcMigrateWasFound {
			insideFuncMigrate(i)
		} else if currentState == commentSingleLine {
			insideCommentSingleLine()
		} else if currentState == commentMultiLine {
			insideCommentMultiLine()
		} else {
			currentState = expectingFOfFunc
		}

		lastCharWasBackSlash = (currentChar == "\\")
		lastCharWasSlash = (currentChar == "/")
	}

	if currentState == 11 && len(stack) == 0 {
		outputStr = code[0:lastCloseBracket-1] + "\n    " + newLine + "\n" + code[lastCloseBracket-1:len(code)]
	} else {
		err = errors.New("Could not parse Migrate file. Please, check sintaxe.")
	}

	return outputStr, err
}
