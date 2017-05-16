package swiftcov

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	r := require.New(t)

	f := &Formatter{}
	files, err := f.Search("examples")
	t.Log(files)
	if err != nil {
		t.Fatal(err)
	}
	err = f.Parse()
	r.NoError(err)

	r.Len(f.SourceFiles, 3)

	sf := f.SourceFiles[0]
	r.Equal("examples/Calculator.swift.gcov", sf.Name)
	r.InDelta(70.8, sf.CoveredPercent, 1)
	r.Len(sf.Coverage, 61)
	r.False(sf.Coverage[15].Valid)
	r.False(sf.Coverage[27].Valid)
	r.True(sf.Coverage[26].Valid)
	r.Equal(0, sf.Coverage[53].Int)
	r.Equal(1, sf.Coverage[48].Int)
	r.Equal(2, sf.Coverage[18].Int)

	report, _ := f.Format()
	r.InDelta(68.2, report.CoveredPercent, 1)
}
