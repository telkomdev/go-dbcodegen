package sqlgen

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Flag struct {
	OutputDirectory string
	OutputTarget    string
	SkipDropTable   bool
}

func NewFlag(dir, target string, skipDrop bool) (*Flag, error) {
	if strings.Contains(target, "/") {
		return nil, errors.New("output target cannot contain \"\\\" character")
	}

	t := time.Now()
	target = fmt.Sprintf("%s_%s", t.Format("20060102150405"), target)
	flag := Flag{
		OutputDirectory: dir,
		OutputTarget:    target,
		SkipDropTable:   skipDrop,
	}
	return &flag, nil
}
