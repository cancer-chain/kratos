package chaindb

import (
	"encoding/base64"
	"encoding/hex"
	"github.com/go-pg/pg/v10"
	"time"
)
import "strings"

var ErrDatabase *pg.DB

func ToJson(str string) string {
	str = strings.Replace(str, "\n", ",", -1)
	str = strings.Replace(str, "\"", "", -1)
	a := strings.Split(str, ",")

	k := "{"
	for i := 0; i < len(a); i++ {
		k += getV(a[i]) + ","
	}
	k = strings.TrimRight(k, ",")
	k += "}"
	return k
}

func getV(s string) string {
	v := ""
	a := strings.Split(s, ":")

	key := strings.Trim(a[0], " ")
	if len(key) > 0 {
		v = `"` + key + `"` + ":" + `"`
		for i := 1; i < len(a); i++ {
			value := strings.Trim(a[i], " ")
			v += value
		}
		v += `"`
	}
	return v
}

func Hash2base64(hash []byte) string {
	return base64.StdEncoding.EncodeToString(hash)
}
func Hash2Hex(hash []byte) string {
	return hex.EncodeToString(hash)
}

func TimeFormat(t time.Time) string {
	return t.Format("2006-01-02T15:04:05.999999999Z")
}
