package tempwriter_test

import (
	"github.com/visor-tax/firemodel"
	"github.com/visor-tax/firemodel/internal/tempwriter"
)

var _ firemodel.SourceCoder = &tempwriter.TempWriter{}
