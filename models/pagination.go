package models

import (
	"net/url"
	"strconv"

	"github.com/mytheresa/go-hiring-challenge/errors"
)

const (
	defaultLimit  = 10
	defaultOffset = 0
	maxLimit      = 100
	minLimit      = 1
)

type PaginationInput struct {
	Offset int
	Limit  int
}

func (p *PaginationInput) Parse(q url.Values) error {
	var err error
	p.Offset = defaultOffset
	p.Limit = defaultLimit

	offset := q.Get("offset")
	if offset != "" {
		p.Offset, err = strconv.Atoi(offset)
		if err != nil {
			return errors.ErrInvalidOffset
		}
	}

	limit := q.Get("limit")
	if limit != "" {
		p.Limit, err = strconv.Atoi(limit)
		if err != nil {
			return errors.ErrInvalidLimit
		}
	}

	return p.Validate()
}

func (p *PaginationInput) Validate() error {
	if p.Limit > maxLimit || p.Limit < minLimit {
		return errors.ErrInvalidLimit
	}
	return nil
}
