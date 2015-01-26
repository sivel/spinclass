package forms

import "github.com/mholt/binding"

type CountForm struct {
	Count int
}

func (c *CountForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&c.Count: "count",
	}
}

type PrefixForm struct {
	Prefix string
}

func (p *PrefixForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&p.Prefix: "prefix",
	}
}
