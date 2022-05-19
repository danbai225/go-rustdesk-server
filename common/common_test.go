package common

import (
	"encoding/base64"
	logs "github.com/danbai225/go-logs"
	"testing"
)

var pks = []byte{57, 229, 110, 42, 78, 103, 148, 120, 151, 167, 224, 15, 54, 125, 24, 222, 144, 181, 35, 10, 112, 159, 142, 9, 26, 170, 208, 10, 197, 170, 152, 28}

func TestGenKey(t *testing.T) {
	LoadKey()

	decodeString, _ := base64.StdEncoding.DecodeString("OeVuKk5nlHiXp+APNn0Y3pC1Iwpwn44JGqrQCsWqmBw=")
	logs.Info(decodeString)
	logs.Info(pks)
}
