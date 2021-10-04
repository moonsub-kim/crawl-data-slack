package hackernews

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Filter interface {
	Filter(subText string) bool
	Reason() string
	String() string
}

type filter struct {
	reason string
	parsed string
}

func (f *filter) Reason() string {
	return f.reason
}

func (f *filter) String() string {
	return f.parsed
}

type adFilter struct {
	filter
}

func (f *adFilter) Filter(subText string) bool {
	if !strings.Contains(subText, "comment") && !strings.Contains(subText, "discuss") {
		f.reason = "advertisement filter"
		return true
	}

	return false
}

type ageFilter struct {
	filter
}

func (f *ageFilter) Filter(subText string) bool {
	age := regexp.MustCompile(`(\d+) [A-z]+ ago`).FindStringSubmatch(subText)
	if len(age) != 2 {
		f.reason = "age is not matched"
		return true
	} else if strings.Contains(age[0], "minute") {
		f.reason = "ignore recent 1h post"
		return true
	} else if v, err := strconv.Atoi(age[1]); err != nil {
		f.reason = fmt.Sprintf("Atoi error %v", err)
		return true
	} else if strings.Contains(age[0], "hour") && v < 2 {
		f.reason = "ignore recent 1h post"
		return true
	}
	f.parsed = age[0]

	return false
}

const POINT_THRESHOLD int = 100

type pointFilter struct {
	filter
}

func (f *pointFilter) Filter(subText string) bool {
	point := regexp.MustCompile(`(\d+) poins?t`).FindStringSubmatch(subText)
	if len(point) != 2 {
		f.reason = "point is not matched"
		return true
	} else if v, err := strconv.Atoi(point[1]); err != nil {
		f.reason = fmt.Sprintf("Atoi error %v", err)
		return true
	} else if v < POINT_THRESHOLD {
		f.reason = "ignore less than 40 point"
		return true
	}
	f.parsed = point[0]

	return false
}

type commentFilter struct {
	filter
}

func (f *commentFilter) Filter(subText string) bool {
	comments := regexp.MustCompile(`\d+.comments?`).FindString(subText)
	f.parsed = comments
	if comments == "" {
		f.parsed = "discuss"
	}

	return false
}
