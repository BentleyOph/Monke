package parser

import (
	"fmt"
	"strings"
)

var traceLevel int = 0 // Variable to keep track of the current trace level

const traceIdentPlaceholder string = "\t" // Placeholder string used for indentation in tracing

// identLevel returns a string with the appropriate indentation based on the current trace level
func identLevel() string {
	return strings.Repeat(traceIdentPlaceholder, traceLevel-1)
}

// tracePrint prints the traced message with the appropriate indentation
func tracePrint(fs string) {
	fmt.Printf("%s%s\n", identLevel(), fs)
}

// incIdent increments the trace level by 1
func incIdent() {
	traceLevel = traceLevel + 1
}

// decIdent decrements the trace level by 1
func decIdent() {
	traceLevel = traceLevel - 1
}

// trace is a function that starts a new trace and prints the "BEGIN" message
func trace(msg string) string {
	incIdent()
	tracePrint("BEGIN " + msg)
	return msg
}

// untrace is a function that ends the current trace and prints the "END" message
func untrace(msg string) {
	tracePrint("END " + msg)
	decIdent()
}
