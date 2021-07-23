package unitsetter_test

import (
	"testing"

	"github.com/nickwells/param.mod/v5/param"
	"github.com/nickwells/param.mod/v5/paramtest"
	"github.com/nickwells/testhelper.mod/testhelper"
	"github.com/nickwells/units.mod/v2/units"
	"github.com/nickwells/unitsetter.mod/v4/unitsetter"
)

const (
	updFlagNameFamilysetter     = "upd-gf-familysetter"
	keepBadFlagNameFamilysetter = "keep-bad-familysetter"
)

var commonFamilysetterGFC = testhelper.GoldenFileCfg{
	DirNames:               []string{"testdata", "FamilySetter"},
	Pfx:                    "gf",
	Sfx:                    "txt",
	UpdFlagName:            updFlagNameFamilysetter,
	KeepBadResultsFlagName: keepBadFlagNameFamilysetter,
}

func init() {
	commonFamilysetterGFC.AddUpdateFlag()
	commonFamilysetterGFC.AddKeepBadResultsFlag()
}

func TestFamilySetter(t *testing.T) {
	const dfltParamName = "set-family"
	var f *units.Family

	testCases := []paramtest.Setter{
		{
			ID:      testhelper.MkID("bad-setter-no-value"),
			PSetter: unitsetter.FamilySetter{},
			ExpPanic: testhelper.MkExpPanic(
				dfltParamName + ": unitsetter.FamilySetter Check failed: " +
					"the Value to be set is nil"),
		},
		{
			ID: testhelper.MkID("good-setter-bad-val"),
			PSetter: unitsetter.FamilySetter{
				Value: &f,
			},
			ParamVal: "nonesuch",
			SetWithValErr: testhelper.MkExpErr(
				`There is no unit family called "nonesuch".`),
		},
		{
			ID: testhelper.MkID("good-setter-bad-val-close-match"),
			PSetter: unitsetter.FamilySetter{
				Value: &f,
			},
			ParamVal: "dostance",
			SetWithValErr: testhelper.MkExpErr(
				`There is no unit family called "dostance".`,
				"Did you mean: distance?"),
		},
		{
			ID: testhelper.MkID("good-setter-good-val"),
			PSetter: unitsetter.FamilySetter{
				Value: &f,
			},
			ParamVal: units.Distance,
		},
		{
			ID: testhelper.MkID("good-setter-good-val-alias"),
			PSetter: unitsetter.FamilySetter{
				Value: &f,
			},
			ParamVal:     "length",
			ValDescriber: true,
		},
	}

	for _, tc := range testCases {
		tc.GFC = commonFamilysetterGFC
		if tc.ParamName == "" {
			tc.ParamName = dfltParamName
		}
		tc.SetVR(param.Mandatory)

		f = units.SampleFamily

		tc.Test(t)
	}
}
