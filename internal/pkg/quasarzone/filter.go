package quasarzone

import "fmt"

type Filter interface {
	Filter(dto DTO) bool
	Reason() string
}

type filter struct {
	reason string
}

func (f *filter) Reason() string {
	return f.reason
}

type statusFilter struct {
	filter
}

func (f *statusFilter) Filter(dto DTO) bool {
	if dto.Status == "종료" {
		f.reason = fmt.Sprintf("종료 == (%s)", dto.Status)
		return true
	}
	return false
}
