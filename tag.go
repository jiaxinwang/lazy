package lazy

import (
	"errors"
	"strings"
)

func disassembleTag(tag string) (name, id, foreignkeyTable, foreignkey string, err error) {
	part := strings.Split(tag, `;`)
	for _, v := range part {
		if strings.Contains(v, `:`) {
			pair := strings.Split(v, ":")
			if len(pair) != 2 {
				err = errors.New("wrong format")
				return
			}
			switch pair[0] {
			case `foreign`:
				f := strings.Split(pair[1], "->")
				if len(f) != 2 {
					err = errors.New("wrong format")
					return
				}
				id = f[0]
				ff := strings.Split(f[1], ".")
				if len(f) != 2 {
					err = errors.New("wrong format")
					return
				}
				foreignkeyTable = ff[0]
				foreignkey = ff[1]
			}
		} else {
			name = v
		}
	}
	return
}
