package internal

import (
	"fmt"
	"strconv"
	"strings"
)

type Semantics struct {
	major, minor, patch uint8
}

func NewSemantics(version string) (sem Semantics, err error) {
	version = strings.TrimLeft(version, "v")
	s := strings.Split(version, ".")

	if len(s) < 2 || len(s) > 3 {
		err = ErrInvalidVersionFormat
		return
	}

	for idx, semverItem := range s {
		var num int
		num, err = strconv.Atoi(semverItem)
		if err != nil {
			err = ErrInvalidVersionFormat
			return
		}
		switch idx {
		case 0:
			sem.major = uint8(num)
		case 1:
			sem.minor = uint8(num)
		case 2:
			sem.patch = uint8(num)
		}
	}

	return
}

func (s Semantics) String() string {
	v := fmt.Sprintf("%d.%d", s.major, s.minor)
	if s.patch != 0 {
		v = fmt.Sprintf("%s.%d", v, s.patch)
	}
	return v
}
