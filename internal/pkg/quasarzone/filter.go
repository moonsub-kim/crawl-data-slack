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

type pcOnlyFilter struct {
	filter
}

func (f *pcOnlyFilter) Filter(dto DTO) bool {
	if dto.Category != "PC/하드웨어" {
		f.reason = fmt.Sprintf("PC/하드웨어 != (%s)", dto.Category)
		return true
	}
	return false
}

type statusFilter struct {
	filter
}

func (f *statusFilter) Filter(dto DTO) bool {
	if dto.Status != "진행중" {
		f.reason = fmt.Sprintf("진행중 != (%s)", dto.Status)
		return true
	}
	return false
}
