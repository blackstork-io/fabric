package diagnostics

import "github.com/hashicorp/hcl/v2"

type Diagnostics hcl.Diagnostics

func (d Diagnostics) Error() string {
	return (hcl.Diagnostics)(d).Error()
}

// Appends diag to diagnostics, returns true if the just-appended diagnostic is an error
func (d *Diagnostics) Append(diag *hcl.Diagnostic) (addedErrors bool) {
	*d = append(*d, diag)
	return diag.Severity == hcl.DiagError
}

// Appends all diags to diagnostics, returns true if the just-appended diagnostics contain an error
func (d *Diagnostics) Extend(diags Diagnostics) (addedErrors bool) {
	*d = append(*d, diags...)
	return diags.HasErrors()
}

func (d *Diagnostics) ExtendHcl(diags hcl.Diagnostics) (addedErrors bool) {
	*d = append(*d, diags...)
	return diags.HasErrors()
}

// HasErrors returns true if the receiver contains any diagnostics of
// severity DiagError.
func (d *Diagnostics) HasErrors() bool {
	return (*hcl.Diagnostics)(d).HasErrors()
}

// Create diagnostic and append it if err !=nil
func (d *Diagnostics) FromErr(err error, summary string) (addedErrors bool) {
	if err == nil {
		return false
	}
	// for FromErr to be inlined more often
	*d = append(*d, FromErr(err, summary))
	return true
}

func FromErr(err error, summary string) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  summary,
		Detail:   err.Error(),
		Extra:    err,
	}
}
