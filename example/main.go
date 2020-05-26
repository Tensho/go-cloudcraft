package main

import (
	"../cloudcraft"

	log "github.com/sirupsen/logrus"
)

func main() {
	c := cloudcraft.NewClient(nil)
	blueprints, _, err := c.Blueprints.List()
	if err != nil {
		log.Error(err)
		return
	}

	log.Infof("Blueprints\n\n%+v\n", blueprints)

	id := blueprints[0].ID
	blueprint, _, err := c.Blueprints.Get(*id)
	if err != nil {
		log.Error(err)
		return
	}

	log.Infof("Blueprint\n\n%+v\n", blueprint)
}
