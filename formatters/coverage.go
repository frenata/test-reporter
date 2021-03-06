package formatters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/markbates/pop/nulls"
	"github.com/pkg/errors"
)

type Coverage []nulls.Int

// MarshalJSON marshals the coverage into JSON. Since the Code Climate
// API requires this as a string "[1,2,null]" and not just a straight
// JSON array we have to do a bunch of work to coerce into that format
func (c Coverage) MarshalJSON() ([]byte, error) {
	cc := make([]interface{}, 0, len(c))
	for _, x := range c {
		cc = append(cc, x)
	}
	bb := &bytes.Buffer{}
	err := json.NewEncoder(bb).Encode(cc)
	if err != nil {
		return bb.Bytes(), err
	}
	b, err := json.Marshal(strings.TrimSpace(bb.String()))

	if err != nil {
		fmt.Printf("### err -> %+v\n", err)
		return b, errors.WithStack(err)
	}

	return b, nil
}

func (c *Coverage) UnmarshalJSON(text []byte) error {
	q := []byte("\"")
	text = bytes.TrimPrefix(text, q)
	text = bytes.TrimSuffix(text, q)
	cc := []nulls.Int{}
	err := json.Unmarshal(text, &cc)
	if err != nil {
		return err
	}
	*c = append(*c, cc...)
	return nil
}
