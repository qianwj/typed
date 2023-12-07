package options

import "go.mongodb.org/mongo-driver/mongo/options"

type Collation struct {
	internal options.Collation
}

func NewCollation() *Collation {
	return &Collation{}
}

func (c *Collation) Locale(locale string) *Collation {
	c.internal.Locale = locale
	return c
}

func (c *Collation) CaseLevel(caseLevel bool) *Collation {
	c.internal.CaseLevel = caseLevel
	return c
}

func (c *Collation) CaseFirst(caseFirst string) *Collation {
	c.internal.CaseFirst = caseFirst
	return c
}

func (c *Collation) Strength(strength int) *Collation {
	c.internal.Strength = strength
	return c
}

func (c *Collation) NumericOrdering(numericOrdering bool) *Collation {
	c.internal.NumericOrdering = numericOrdering
	return c
}

func (c *Collation) Alternate(alternate string) *Collation {
	c.internal.Alternate = alternate
	return c
}

func (c *Collation) MaxVariable(maxVariable string) *Collation {
	c.internal.MaxVariable = maxVariable
	return c
}

func (c *Collation) Normalization(normalization bool) *Collation {
	c.internal.Normalization = normalization
	return c
}

func (c *Collation) Backwards(backwards bool) *Collation {
	c.internal.Backwards = backwards
	return c
}

func (c *Collation) Raw() *options.Collation {
	return &(c.internal)
}
