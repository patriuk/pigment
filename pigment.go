package main

import (
	"fmt"
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

type StyleFunc = func(string) string

var styles = make(map[string]StyleFunc)

// TODO: create generator to generate all functions for all styles and add them to  struct
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

func createStyle(code string) StyleFunc {
	s := esc + code + "m"
	return func(text string) string {
		return s + text + reset
	}
}

func AddStyle(name string, code string) {
	styles[name] = createStyle(code)
}

func Custom(name string) StyleFunc {
	return styles[name]
}

func main() {
	fmt.Println(Blue("some blue text here"))

	AddStyle("green", "32")
	AddStyle("red", "31")

	fmt.Println(Custom("green")("some green text here"))

	fmt.Println(Blue("Hello") + " World " + Custom("red")("!!!!"))

	fmt.Println(Blue("Hello" + Custom("red")("substring?") + "World"))

	x := fmt.Sprintf("My Name is %s, %s, "+Custom("red")("%s"), "what", Blue("fml"), "hello world")
	fmt.Println(x)

	fmt.Println(Blue(BgBlack(Bold("some blue bold text on black background"))))

	green := Custom("green")
	fmt.Println(green("some green text"))
}
