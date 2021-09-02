package main

import SMLogs "github.com/binaryelements/smlogs-go"

var SMLogger SMLogs.Config

func main() {
	SMLogger.New("ABC-PUBLIC-API", "fd3tgg-8081-4d4d5-bfed-431cc3ea9eb82", "https://abclogs.smoothorders.com", "DEBUG", "Y", "Y")
	SMLogger.Error("It's an error!")
}
