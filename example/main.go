package main

import (
	"fmt"
	SMLogs "github.com/binaryelements/smlogs-go"
)

var SMLogger SMLogs.Config

func main() {
	var er error
	fmt.Println("Running test")

	test := []string{"abc", "def"}
	SMLogger.New("MAIN", "na", "https://logs.smoothorders.com", "DEBUG", "Y", "Y")
	SMLogger.Error("It's an error!", test, er)
}
