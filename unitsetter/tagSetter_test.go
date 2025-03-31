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
	updFlagNameTagsetter     = "upd-gf-tagsetter"
	keepBadFlagNameTagsetter = "keep-bad-tagsetter"
)

var commonTagsetterGFC = testhelper.GoldenFileCfg{
	DirNames:               []string{"testdata", "TagSetter"},
	Pfx:                    "gf",
	Sfx:                    "txt",
	UpdFlagName:            updFlagNameTagsetter,
	KeepBadResultsFlagName: keepBadFlagNameTagsetter,
}

func init() {
	commonTagsetterGFC.AddUpdateFlag()
	commonTagsetterGFC.AddKeepBadResultsFlag()
}

func TestTagSetter(t *testing.T) {
	const dfltParamName = "set-tag"

	var tag units.Tag

	testCases := []paramtest.Setter{
		{
			ID:      testhelper.MkID("bad-setter-no-value"),
			PSetter: unitsetter.TagSetter{},
			ExpPanic: testhelper.MkExpPanic(
				dfltParamName + ": unitsetter.TagSetter Check failed: " +
					"the Value to be set is nil"),
		},
		{
			ID: testhelper.MkID("good-setter-bad-val"),
			PSetter: unitsetter.TagSetter{
				Value: &tag,
			},
			ParamVal: "nonesuch",
			SetWithValErr: testhelper.MkExpErr(
				`there is no unit tag called "nonesuch".`),
		},
		{
			ID: testhelper.MkID("good-setter-bad-val-close-match"),
			PSetter: unitsetter.TagSetter{
				Value: &tag,
			},
			ParamVal: "histeric",
			SetWithValErr: testhelper.MkExpErr(
				`there is no unit tag called "histeric".`,
				`Did you mean: "historic"?`),
		},
		{
			ID: testhelper.MkID("good-setter-good-val"),
			PSetter: unitsetter.TagSetter{
				Value: &tag,
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

		tag = units.Tag("")

		tc.Test(t)
	}
}
