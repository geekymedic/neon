package id

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var ran = rand.NewSource(time.Now().Unix())

func GenerateId() string{
	const shortForm = "20060102150405.000000"
	t := time.Now()
	temp := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local)
	str := temp.Format(shortForm)
	str = strings.Replace(str,".","",1)
	str += strconv.FormatInt(ran.Int63(),10)[:6]
	return str
}


