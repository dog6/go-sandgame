package main

import (
	"git.smallzcomputing.com/sand-game/src/config"
	"git.smallzcomputing.com/sand-game/src/game"
	"git.smallzcomputing.com/sand-game/src/util"
)

var Conf config.Configuration

func main() {

	Conf.ReadConfig()
	util.VerboseLogging = Conf.VerboseLogging
	util.Log("Reading config..")

	util.Log("Starting game..")
	game.Start(&Conf)
}
