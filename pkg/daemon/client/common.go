package client

import (
	"time"

	jsoniter "github.com/json-iterator/go"
)

const (
	FioTool = "fio"
)

var (
	FioCommandsTimeout = 15 * time.Second

	json = jsoniter.ConfigCompatibleWithStandardLibrary
)
