package main

import (
	"github.com/traPtitech/traPortfolio/model"
)

func main() {
	if err := model.Setup(); err != nil {
		panic(err)
	}
}
