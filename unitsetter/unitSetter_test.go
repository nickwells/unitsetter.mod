package unitsetter_test

import (
	"testing"

	"github.com/nickwells/param.mod/v5/param"
	"github.com/nickwells/param.mod/v5/paramtest"
	"github.com/nickwells/testhelper.mod/testhelper"
	"github.com/nickwells/units.mod/units"
	"github.com/nickwells/unitsetter.mod/v3/unitsetter"
)

const (
	updFlagNameUnitsetter     = "upd-gf-unitsetter"
	keepBadFlagNameUnitsetter = "keep-bad-unitsetter"
)

var commonUnitsetterGFC = testhelper.GoldenFileCfg{
	DirNames:               []string{"testdata", "UnitSetter"},
	Pfx:                    "gf",
	Sfx:                    "txt",
	UpdFlagName:            updFlagNameUnitsetter,
	KeepBadResultsFlagName: keepBadFlagNameUnitsetter,
}

func init() {
	commonUnitsetterGFC.AddUpdateFlag()
	commonUnitsetterGFC.AddKeepBadResultsFlag()
}

func TestUnitSetter(t *testing.T) {
	const dfltParamName = "set-unit"
	var u units.Unit
	testCases := []paramtest.Setter{
		{
			ID: testhelper.MkID("bad-value"),
			PSetter: unitsetter.UnitSetter{
				Value: nil,
				UD:    units.GetUnitDetailsOrPanic(units.Distance),
			},
			ExpPanic: testhelper.MkExpPanic(
				dfltParamName + ": unitsetter.UnitSetter Check failed: " +
					"the Value to be set is nil"),
		},
		{
			ID: testhelper.MkID("bad-unit-details"),
			PSetter: unitsetter.UnitSetter{
				Value: &u,
			},
			ExpPanic: testhelper.MkExpPanic(
				dfltParamName + ": unitsetter.UnitSetter Check failed: " +
					"there are no valid alternative units"),
		},
		{
			ID: testhelper.MkID("bad-unit-details"),
			PSetter: unitsetter.UnitSetter{
				Value:  &u,
				UD:     units.GetUnitDetailsOrPanic(units.Distance),
				Checks: []unitsetter.UnitCheckFunc{nil},
			},
			ExpPanic: testhelper.MkExpPanic(
				dfltParamName + ": unitsetter.UnitSetter Check failed: " +
					"one of the check functions is nil"),
		},
		{
			ID: testhelper.MkID("distance-bad-unit"),
			PSetter: unitsetter.UnitSetter{
				Value: &u,
				UD:    units.GetUnitDetailsOrPanic(units.Distance),
			},
			ParamVal: "nonesuch",
			SetWithValErr: testhelper.MkExpErr(
				"'nonesuch' is not a recognised unit of distance."),
		},
		{
			ID: testhelper.MkID("distance-good-unit"),
			PSetter: unitsetter.UnitSetter{
				Value: &u,
				UD:    units.GetUnitDetailsOrPanic(units.Distance),
			},
			ParamVal: "mile",
		},
	}

	for _, tc := range testCases {
		tc.GFC = commonUnitsetterGFC
		if tc.ParamName == "" {
			tc.ParamName = dfltParamName
		}
		tc.SetVR(param.Mandatory)

		u = units.Unit{}

		tc.Test(t)
	}
}
