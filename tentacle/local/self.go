package local

import (
	"github.com/sirupsen/logrus"
	"ostopus/shared/tentacle"
)

var self tentacle.Tentacle

func InitSelf(name, address string) {
	self.Name = name
	self.Address = address
	logrus.WithFields(logrus.Fields{"name": name, "address": address}).Info("Tentacle initialised")
}

func GetSelf() tentacle.Tentacle {
	return self
}
