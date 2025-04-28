package unitsetter_test

import (
	"testing"

	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/paramtest"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
	"github.com/nickwells/units.mod/v2/units"
	"github.com/nickwells/unitsetter.mod/v4/unitsetter"
)

const (
	updFlagNameTaglistappender     = "upd-gf-taglistappender"
	keepBadFlagNameTaglistappender = "keep-bad-taglistappender"
)

var commonTaglistappenderGFC = testhelper.GoldenFileCfg{
	DirNames:               []string{"testdata", "TagListAppender"},
	Pfx:                    "gf",
	Sfx:                    "txt",
	UpdFlagName:            updFlagNameTaglistappender,
	KeepBadResultsFlagName: keepBadFlagNameTaglistappender,
}

func init() {
	commonTaglistappenderGFC.AddUpdateFlag()
	commonTaglistappenderGFC.AddKeepBadResultsFlag()
}

func TestTagListAppender(t *testing.T) {
	const dfltParamName = "set-tag"

	var tags []units.Tag

	testCases := []paramtest.Setter{
		{
			ID:      testhelper.MkID("bad-setter-no-value"),
			PSetter: unitsetter.TagListAppender{},
			ExpPanic: testhelper.MkExpPanic(
				dfltParamName + ": unitsetter.TagListAppender Check failed: " +
					"the Value to be set is nil"),
		},
		{
			ID: testhelper.MkID("good-setter-bad-val"),
			PSetter: unitsetter.TagListAppender{
				Value: &tags,
			},
			ParamVal: "nonesuch",
			SetWithValErr: testhelper.MkExpErr(
				`there is no unit tag called "nonesuch"`),
		},
		{
			ID: testhelper.MkID("good-setter-bad-val-close-match"),
			PSetter: unitsetter.TagListAppender{
				Value: &tags,
			},
			ParamVal: "histeric",
			SetWithValErr: testhelper.MkExpErr(
				`there is no unit tag called "histeric",`,
				`did you mean "historic"?`),
		},
		{
			ID: testhelper.MkID("good-setter-good-val"),
			PSetter: unitsetter.TagListAppender{
				Value: &tags,
			},
			ParamVal: string(units.TagHist),
		},
	}

	for _, tc := range testCases {
		tc.GFC = commonTagsetterGFC
		if tc.ParamName == "" {
			tc.ParamName = dfltParamName
		}

		tc.SetVR(param.Mandatory)

		tags = []units.Tag{}

		tc.Test(t)
	}
}
