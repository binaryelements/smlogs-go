package main

import "github.com/binaryelements/SMLogs"

var SMLogger SMLogs.Config

func main() {
	SMLogger.New("ABC-PUBLIC-API", "fd3tgg-8081-4d4d5-bfed-431cc3ea9eb82", "https://abclogs.smoothorders.com", "DEBUG", "Y", "Y")
}
