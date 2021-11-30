package eomisae

import "context"

type action struct {
	links []string
	dtos  []DTO
}

func (a action) actionFunc(ctx context.Context) error {
	return nil
}

func (a action) getDTOs() []DTO {
	return a.dtos
}

func NewAction(links []string) action {
	return action{
		links: links,
	}
}
