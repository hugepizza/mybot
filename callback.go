package main

import (
	"fmt"
	"strconv"
	"strings"
)

type CallbackType string
type CallbackData string

const (
	CallbackTypeRouteKey CallbackType = "r"
	CallbackTypePointKey              = "p"
	CallbackTypeVerify                = "f"
)

func (d CallbackData) String() string {
	return string(d)
}

func NewCallbackData(ctype CallbackType, routeIndexes ...int) CallbackData {
	var strs []string
	for _, index := range routeIndexes {
		strs = append(strs, strconv.Itoa(index))
	}
	return CallbackData(fmt.Sprintf("%s:%s", ctype, strings.Join(strs, ",")))
}

func (d CallbackData) GetCallbackType() CallbackType {
	return CallbackType(strings.Split(string(d), ":")[0])
}

func (d CallbackData) GetCallbackIndex() []int {
	var indexes []int
	indexStr := strings.Split(string(d), ":")[1]
	indexStrs := strings.Split(indexStr, ",")
	for _, indexStr := range indexStrs {
		index, _ := strconv.Atoi(indexStr)
		indexes = append(indexes, index)
	}
	return indexes
}
