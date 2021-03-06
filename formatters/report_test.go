package formatters

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/markbates/pop/nulls"
	"github.com/stretchr/testify/require"
)

func Test_Report_Merge_Bad_GitHead(t *testing.T) {
	r := require.New(t)
	a := &Report{
		Git: ccGit{
			Head: "a",
		},
	}
	b := &Report{
		Git: ccGit{
			Head: "b",
		},
	}
	err := a.Merge(b)
	r.Error(err)
	r.Equal("git heads do not match", err.Error())
}

func Test_Report_Merge_MismatchedCoverageLength(t *testing.T) {
	r := require.New(t)
	a := &Report{
		Git: ccGit{
			Head: "a",
		},
		SourceFiles: SourceFiles{
			"a.go": {
				Name:     "a.go",
				Coverage: Coverage{nulls.NewInt(1)},
			},
		},
	}
	b := &Report{
		Git: ccGit{
			Head: "a",
		},
		SourceFiles: SourceFiles{
			"a.go": {
				Name:     "a.go",
				Coverage: Coverage{nulls.NewInt(1), nulls.NewInt(2)},
			},
		},
	}
	err := a.Merge(b)
	r.Error(err)
	r.Equal("coverage length mismatch for a.go", err.Error())
}

func Test_Report_Merge(t *testing.T) {
	r := require.New(t)
	reps := []*Report{}
	for i := 0; i < 4; i++ {
		rep, err := NewReport()
		r.NoError(err)

		f, err := os.Open(fmt.Sprintf("../examples/codeclimate.%d.json", i))
		r.NoError(err)
		err = json.NewDecoder(f).Decode(&rep)
		r.NoError(err)

		sf := rep.SourceFiles["config/initializers/resque.rb"]
		r.NotNil(sf)
		r.Equal(14, sf.LineCounts.Total)

		reps = append(reps, &rep)
	}
	main := reps[0]
	main.Merge(reps[1:]...)
	r.Equal(19379, main.LineCounts.Total)
	r.Equal(2564, main.LineCounts.Missed)
	r.Equal(16815, main.LineCounts.Covered)
	r.InDelta(86.76, main.LineCounts.CoveredPercent(), 1)

	sf := main.SourceFiles["lib/code_climate/polymorphic_routes.rb"]
	r.NotNil(sf)
	r.Equal(196, sf.LineCounts.Total)
	r.Equal(59, sf.LineCounts.Missed)
	r.Equal(137, sf.LineCounts.Covered)
	r.InDelta(69.8, sf.CoveredPercent, 1)
}

func Test_Report_JSON_Unmarshal(t *testing.T) {
	r := require.New(t)
	f, err := os.Open("../examples/codeclimate.json")
	r.NoError(err)

	rep, err := NewReport()
	r.NoError(err)
	err = json.NewDecoder(f).Decode(&rep)
	r.NoError(err)

	r.Equal(20, len(rep.SourceFiles))
	r.Equal("/go/src/github.com/codeclimate/test-reporter/simplecov-test-reporter", rep.Environment.PWD)

	sf := rep.SourceFiles["lib/code_climate/test_reporter/client.rb"]
	r.NotNil(sf)
	r.InDelta(87.87, sf.CoveredPercent, 1)

	lc := sf.LineCounts
	r.Equal(8, lc.Missed)
	r.Equal(58, lc.Covered)
	r.Equal(66, lc.Total)
}

func Test_Merge_Issue_103(t *testing.T) {
	r := require.New(t)

	a, err := NewReport()
	r.NoError(err)

	a.SourceFiles = SourceFiles{
		"app/jobs/initialize_account_seats.rb": SourceFile{
			Name:           "app/jobs/initialize_account_seats.rb",
			CoveredPercent: 100,
			Coverage:       Coverage{nulls.NewInt(1), nulls.NewInt(1), nulls.NewInt(15), nulls.Int{}, nulls.Int{}, nulls.NewInt(1), nulls.NewInt(3), nulls.Int{}, nulls.Int{}, nulls.NewInt(1), nulls.Int{}, nulls.NewInt(1), nulls.NewInt(3), nulls.Int{}, nulls.Int{}},
			LineCounts: LineCounts{
				Missed:  0,
				Covered: 8,
				Total:   8,
			},
		},
	}

	b, err := NewReport()
	r.NoError(err)

	b.SourceFiles = SourceFiles{
		"app/jobs/initialize_account_seats.rb": SourceFile{
			Name:           "app/jobs/initialize_account_seats.rb",
			CoveredPercent: 62.5,
			Coverage:       Coverage{nulls.NewInt(1), nulls.NewInt(1), nulls.NewInt(0), nulls.Int{}, nulls.Int{}, nulls.NewInt(1), nulls.NewInt(0), nulls.Int{}, nulls.Int{}, nulls.NewInt(1), nulls.Int{}, nulls.NewInt(1), nulls.NewInt(0), nulls.Int{}, nulls.Int{}},
			LineCounts: LineCounts{
				Missed:  3,
				Covered: 5,
				Total:   8,
			},
		},
	}

	err = a.Merge(&b)
	r.NoError(err)

	r.InDelta(100, a.CoveredPercent, 1)
	r.Equal(0, a.LineCounts.Missed)
	r.Equal(8, a.LineCounts.Covered)
	r.Equal(8, a.LineCounts.Total)
}
