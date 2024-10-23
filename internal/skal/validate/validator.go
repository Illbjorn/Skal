package validate

import (
	"slices"

	"github.com/illbjorn/skal/internal/skal/typeset"
)

var v = NewValidator()

func Reset() {
	clear(v.known)
}

/*------------------------------------------------------------------------------
 * Validator
 *----------------------------------------------------------------------------*/

func NewValidator() *Validator {
	return &Validator{known: make([]*knownType, 0)}
}

type Validator struct {
	query *query
	known []*knownType
}

func (v *Validator) Register(id, dataT, skalT, parent string, object typeset.SkalType) *Validator {
	v.known = append(
		v.known,
		&knownType{
			id:       id,
			dataType: dataT,
			skalType: skalT,
			parent:   parent,
			object:   object,
		},
	)

	return v
}

func (v *Validator) Select() *Validator {
	v.query = newQuery()
	return v
}

// Append a `queryFn` filtering on the `id` field.
func (v *Validator) ID(id string) *Validator {
	// Don't add blank values.
	if id == "" {
		return v
	}

	v.query.queries = append(
		v.query.queries,
		func(kt *knownType) bool {
			return kt.id == id
		})

	return v
}

// Append a `queryFn` filtering on the DataType field.
func (v *Validator) DataType(dataT string) *Validator {
	// Don't add blank values.
	if dataT == "" {
		return v
	}

	v.query.queries = append(
		v.query.queries,
		func(kt *knownType) bool {
			return kt.dataType == dataT
		})

	return v
}

// Append a `queryFn` filtering on the SkalType field.
func (v *Validator) SkalType(skalTs ...string) *Validator {
	// Don't add blank values.
	if len(skalTs) == 0 {
		return v
	}

	v.query.queries = append(
		v.query.queries,
		func(kt *knownType) bool {
			return slices.Contains(skalTs, kt.skalType)
		})

	return v
}

// Append a `queryFn` filtering on the Parent field.
func (v *Validator) Parent(parent string) *Validator {
	// Don't add blank values.
	if parent == "" {
		return v
	}

	v.query.queries = append(
		v.query.queries,
		func(kt *knownType) bool {
			return kt.parent == parent
		})

	return v
}

// Simply wraps `Execute()`, producing a boolean `true`/`false` whether the
// query was successful. Successful meaning returned > 0 results.
func (v *Validator) Exists() bool {
	return len(v.Execute()) > 0
}

// Execute the built query, applying all assembled `queryFn` to the provided
// `validator` instance's `known` slice. Finally, returning any results for
// which all `queryFn`s produced a `true` result.
func (v *Validator) Execute() []*knownType {
	if len(v.query.queries) == 0 || len(v.known) == 0 {
		return nil
	}

	var res []*knownType
	for _, kt := range v.known {
		if kt == nil {
			continue
		}

		if applyQueryFns(kt, v.query.queries) {
			res = append(res, kt)
		}
	}

	clear(v.query.queries)

	return res
}

// Wraps `Execute()` returning the first result.
func (v *Validator) First() *knownType {
	res := v.Execute()
	if len(res) == 0 {
		return nil
	}

	return res[0]
}

/*------------------------------------------------------------------------------
 * Known Type
 *----------------------------------------------------------------------------*/

type knownType struct {
	object   typeset.SkalType
	id       string
	dataType string
	skalType string
	parent   string
}

/*------------------------------------------------------------------------------
 * Validator Query
 *----------------------------------------------------------------------------*/

func newQuery() *query {
	return &query{queries: make([]queryFn, 0)}
}

type queryFn func(*knownType) bool

type query struct {
	queries []queryFn
}

func applyQueryFns(kt *knownType, fns []queryFn) bool {
	for _, fn := range fns {
		if !fn(kt) {
			return false
		}
	}

	return true
}

/*------------------------------------------------------------------------------
 * Generic Supporting Functions
 *----------------------------------------------------------------------------*/

func ascendTo[T typeset.SkalType](v typeset.SkalType) T {
	if v == nil {
		var x T
		return x
	}

	switch x := v.(type) {
	case T:
		return x
	default:
		return ascendTo[T](x.Parent())
	}
}
