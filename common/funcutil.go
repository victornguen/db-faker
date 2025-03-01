package funcutil

import (
	"github.com/samber/mo"
	"strconv"
)

func GetOpt[V any](slice []V, index int) mo.Option[V] {
	if index >= 0 && index < len(slice) {
		return mo.Some(slice[index])
	}
	return mo.None[V]()

}

func ParseInt(s string) (int, bool) {
	var n int
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	return n, true
}
