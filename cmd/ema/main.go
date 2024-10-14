package main

import (
	"quadrimus.com/ema/tool"
	_ "quadrimus.com/ema/tool/ema2json"
	_ "quadrimus.com/ema/tool/json2ema"
)

func main() {
	tool.DoWithOS()
}
