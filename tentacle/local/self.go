package local

import (
	"fmt"

	"github.com/AHerczeg/ostopus/shared"
	"github.com/sirupsen/logrus"
)

var (
	self        shared.Tentacle
	headAddress string
)

func InitSelf(name, address string) {
	self.Name = name
	self.Address = fmt.Sprintf("https://%v", address)
	logrus.WithFields(logrus.Fields{"name": name, "address": address}).Info("Tentacle initialised")
}

func GetSelf() shared.Tentacle {
	return self
}

func GetHeadAddress() string {
	return headAddress
}

func SetHeadAddress(address string) {
	headAddress = address
}
