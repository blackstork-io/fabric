package circularRefDetector_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/assert"

	"github.com/blackstork-io/fabric/pkg/circularRefDetector"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

func TestBasic(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	obj := new(struct{})

	assert.False(circularRefDetector.Check(obj))
	circularRefDetector.Add(obj, nil)
	assert.True(circularRefDetector.Check(obj))
	// duplicate add does nothing
	circularRefDetector.Add(obj, nil)
	assert.True(circularRefDetector.Check(obj))
	circularRefDetector.Remove(obj, nil)
	assert.False(circularRefDetector.Check(obj))
}

func TestDiagnostic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		addRange  bool
		addMarker bool

		wantUnchanged        bool
		filenameForTraceback string
	}{
		{
			addRange:      false,
			addMarker:     false,
			wantUnchanged: true,
		},
		{
			addRange:      true,
			addMarker:     false,
			wantUnchanged: true,
		},
		{
			addRange:      false,
			addMarker:     true,
			wantUnchanged: false,
		},
		{
			addRange:             true,
			addMarker:            true,
			wantUnchanged:        false,
			filenameForTraceback: "some_name",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(
			fmt.Sprintf("CircularRef_range:%v_marker:%v", tt.addRange, tt.addMarker),
			func(t *testing.T) {
				t.Parallel()
				assert := assert.New(t)

				obj := new(int)
				var rng *hcl.Range
				if tt.addRange {
					rng = &hcl.Range{
						Filename: tt.filenameForTraceback,
					}
				}

				circularRefDetector.Add(obj, rng)
				assert.True(circularRefDetector.Check(obj))
				orig_diag := &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Some diag",
				}

				if tt.addMarker {
					orig_diag.Extra = diagnostics.NewTracebackExtra()
				}
				orig_extra_len := 0

				diags := diagnostics.Diag{
					orig_diag,
				}
				circularRefDetector.Remove(obj, &diags)

				assert.Len(diags, 1)
				var new_tb_len int
				var tbe diagnostics.TracebackExtra
				var ok bool
				if tbe, ok = diags[0].Extra.(diagnostics.TracebackExtra); ok {
					new_tb_len = len(tbe.Traceback)
				}
				if tt.wantUnchanged {
					assert.Equal(orig_extra_len, new_tb_len)
				} else {
					assert.NotEqual(orig_extra_len, new_tb_len)
					if tt.filenameForTraceback != "" {
						assert.Contains(tbe.Traceback[0].Filename, tt.filenameForTraceback)
					}
				}
			})
	}
}
