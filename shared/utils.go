package shared

import (
	jsoniter "github.com/json-iterator/go"
)

// JSON : faster implementation of standard JSON library
var JSON = jsoniter.ConfigCompatibleWithStandardLibrary
