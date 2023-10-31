package hashid

import (
	"github.com/ChangSZ/blog/infra/conf"
	"github.com/ChangSZ/blog/infra/log"
	"github.com/speps/go-hashids"
)

type HashIdParams struct {
	Salt      string
	MinLength int
}

var hashIdParams *HashIdParams

func (hd *HashIdParams) SetHashIdSalt(salt string) func(*HashIdParams) interface{} {
	return func(hd *HashIdParams) interface{} {
		hs := hd.Salt
		hd.Salt = salt
		return hs
	}
}

func (hd *HashIdParams) SetHashIdLength(minLength int) func(*HashIdParams) interface{} {
	return func(hd *HashIdParams) interface{} {
		ml := hd.MinLength
		hd.MinLength = minLength
		return ml
	}
}

func (hd *HashIdParams) HashIdInit(options ...func(*HashIdParams) interface{}) (*hashids.HashID, error) {
	q := &HashIdParams{
		Salt:      conf.HASHIDSALT,
		MinLength: conf.HASHIDMINLENGTH,
	}
	for _, option := range options {
		option(q)
	}
	hashIdParams = q
	hds := hashids.NewData()
	hds.Salt = hashIdParams.Salt
	hds.MinLength = hashIdParams.MinLength
	h, err := hashids.NewWithData(hds)
	if err != nil {
		log.Errorf("hash new with data is error: %v", err)
		return nil, err
	}
	return h, nil
}
