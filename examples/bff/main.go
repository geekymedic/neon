package main

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/geekymedic/neon/bff"
)

func main() {
	viper.Set("address", ":9810")
	if err := bff.Main(); err != nil {
		fmt.Println(err)
	}
}