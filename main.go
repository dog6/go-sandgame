package main

import (
	"git.smallzcomputing.com/sand-game/config"
	"git.smallzcomputing.com/sand-game/game"
	"git.smallzcomputing.com/sand-game/util"
)

var Conf config.Configuration

func main() {

	Conf.ReadConfig()
	util.VerboseLogging = Conf.VerboseLogging
	util.Log("Reading config..")

	util.Log("Starting game..")
	game.Start(&Conf)
}
