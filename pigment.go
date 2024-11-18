package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

var esc = "\x1B["
var reset = esc + "0m"

var ansiCodes = map[string]string{
	"Bold":    "1",
	"Italic":  "3",
	"Black":   "30",
	"Green":   "32",
	"Blue":    "34",
	"White":   "37",
	"BgBlack": "40",
	"BgGreen": "42",
	"BgBlue":  "44",
	"BgWhite": "47",
}

type Style = func(string) string

var defaultStyle Style = func(text string) string { return text }

var styleMap = make(map[string]Style)

// TODO: create generator to generate all functions for all styles and add them to struct
func Blue(text string) string {
	s := esc + ansiCodes["Blue"] + "m"
	return s + text + reset
}

func Bold(text string) string {
	s := esc + ansiCodes["Bold"] + "m"
	return s + text + reset
}

func BgBlack(text string) string {
	s := esc + ansiCodes["BgBlack"] + "m"
	return s + text + reset
}

func createStyle(codes ...string) Style {
	s := esc + strings.Join(codes, ";") + "m"

	return func(text string) string {
		return s + text + reset
	}
}

func Add(name string, codes ...string) Style {
	styleMap[name] = createStyle(codes...)
	return styleMap[name]
}

func Apply(name string) Style {
	s, ok := styleMap[name]

	if !ok {
		log.Printf("Style '%s' is not defined. Please define it using 'pigment.Add' before applying it. Applying default styling.", name)

		// TODO: conditional stack tracing based on PIGMENT_STACK_TRACE=true
		// make it false by default
		buf := make([]byte, 1024)
		n := runtime.Stack(buf, false)
		log.Printf("Stack trace:\n%s", buf[:n])

		return defaultStyle
	}
	return s
}

// TODO: accept interface{} and check for type string to add name, type Style for styling
func Mix(name string, styles ...Style) Style {
	mixed := func(text string) string {
		for _, style := range styles {
			text = style(text)
		}
		return text
	}

	if name != "" {
		styleMap[name] = mixed
	}

	return mixed
}

func main() {
	fmt.Println(Blue("some blue text here"))

	Add("green", "32")

	fmt.Println(Apply("green")("some green text here"))

	red := Add("red", "31")
	fmt.Println(red("some red text here"))

	fmt.Println(Blue("Hello") + " World " + Apply("red")("!!!!"))

	fmt.Println(Blue("Hello" + Apply("red")("substring?") + "World"))

	x := fmt.Sprintf("My Name is %s, %s, "+Apply("red")("%s"), "what", Blue("fml"), "hello world")
	fmt.Println(x)

	fmt.Println(Blue(BgBlack(Bold("some blue bold text on black background"))))

	green := Apply("green")
	fmt.Println(green("some green text"))
	fmt.Println("\x1B[30;31;32;43;44msomethingggggg" + reset)

	BlueBgBlackBold := Mix("", Blue, BgBlack, Bold, func(text string) string { return text + "fml" })
	fmt.Println(BlueBgBlackBold("the mix text here"))

	fmt.Println("\nmissing style case:")
	fmt.Println(Apply("BlueBgBlackBold")("the mix text there"))

	BlueBgBlackBoldX := Add("BlueBgBlackBold", "34", "40", "1")
	fmt.Println(BlueBgBlackBoldX("multi code style"))
}
