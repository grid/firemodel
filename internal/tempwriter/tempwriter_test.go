package tempwriter_test

import (
	"github.com/mickeyreiss/firemodel/firemodel"
	"github.com/mickeyreiss/firemodel/internal/tempwriter"
)

var _ firemodel.SourceCoder = &tempwriter.TempWriter{}
