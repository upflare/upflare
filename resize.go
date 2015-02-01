package upflare

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var RESIZE_PATTERN = regexp.MustCompile(`^((s[\d]+)|(w[\d]+(-h[\d]+)?)|(h[\d]+))(-c)?$`)

type Resize struct {
	width  int64
	height int64
	crop   bool
}

func (r *Resize) String() string {
	resize := ""
	switch {
	case r.width > 0 && r.height > 0 && r.width == r.height:
		resize = fmt.Sprintf("s%d", r.width)
	case r.width > 0 && r.height > 0:
		resize = fmt.Sprintf("w%d-h%d", r.width, r.height)
	case r.width > 0:
		resize = fmt.Sprintf("w%d", r.width)
	case r.height > 0:
		resize = fmt.Sprintf("h%d", r.height)
	}
	if resize != "" && r.crop {
		resize += "-c"
	}
	return resize
}

func NewResize(width, height int64, crop bool) *Resize {
	if width < 0 {
		width = 0
	}
	if height < 0 {
		height = 0
	}
	if width == 0 && height == 0 && !crop {
		return nil
	}
	return &Resize{
		width:  width,
		height: height,
		crop:   crop,
	}
}

func NewResizeFromString(resize string) *Resize {
	if !RESIZE_PATTERN.MatchString(resize) {
		return nil
	}

	var width, height int64
	var crop bool

	for _, option := range strings.Split(resize, "-") {
		integer, _ := strconv.ParseInt(option[1:], 10, 32)

		switch option[0] {
		case 's':
			width = integer
			height = integer
		case 'w':
			width = integer
		case 'h':
			height = integer
		case 'c':
			crop = true
		default:
			panic("resize: option not recognized")
		}
	}

	return NewResize(width, height, crop)
}
