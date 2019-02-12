package util

import (
	"bufio"
	"io"
)

// state used in reading the front matter.
const (
	stateStart         = 0
	stateInfrontmatter = 1
	stateInBody        = 2
)

// front matter reader.
type frontmatter struct {
	tag string
}

func Newfrontmatter(tag string) *frontmatter {
	return &frontmatter{tag}
}

func (fm *frontmatter) Parse(input io.Reader) (front, body string, err error) {
	s := bufio.NewScanner(input)
	front, body = "", ""
	state := stateStart
	lines := 0
	for s.Scan() {
		t := s.Text()
		if lines == 0 && t != fm.tag {
			state = stateInBody
		}
		if state == stateInfrontmatter && t != fm.tag {
			front += t + "\n"
		} else if state == stateInBody {
			body += t + "\n"
		}
		if t == fm.tag {
			if state < stateInBody {
				state++
			}
		}
		lines++
	}

	err = s.Err()
	if err != nil {
		return "", "", err
	}

	return front, body, nil
}
