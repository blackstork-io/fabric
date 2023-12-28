package diagnostics

// More ergonomic wrapper over hcl.Diagnostics.

import "github.com/hashicorp/hcl/v2"

type Diag hcl.Diagnostics //nolint:errname // Diagnostics does implement error interface, but not, itself, an error.

func (d Diag) Error() string {
	return (hcl.Diagnostics)(d).Error()
}

// Appends diag to diagnostics, returns true if the just-appended diagnostic is an error.
func (d *Diag) Append(diag *hcl.Diagnostic) (addedErrors bool) {
	if diag != nil {
		*d = append(*d, diag)
		return diag.Severity == hcl.DiagError
	}
	return false
}

// Add new diagnostic error.
func (d *Diag) Add(summary, detail string) {
	*d = append(*d, &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  summary,
		Detail:   detail,
	})
}

// Appends all diags to diagnostics, returns true if the just-appended diagnostics contain an error.
func (d *Diag) Extend(diags Diag) (haveAddedErrors bool) {
	*d = append(*d, diags...)
	return diags.HasErrors()
}

// Appends all diags to diagnostics, returns true if the just-appended diagnostics contain an error.
func (d *Diag) ExtendHcl(diags hcl.Diagnostics) (haveAddedErrors bool) {
	*d = append(*d, diags...)
	return diags.HasErrors()
}

// HasErrors returns true if the receiver contains any diagnostics of
// severity DiagError.
func (d *Diag) HasErrors() bool {
	return (*hcl.Diagnostics)(d).HasErrors()
}

// Creates diagnostic and appends it if err != nil.
func (d *Diag) AppendErr(err error, summary string) (haveAddedErrors bool) {
	// The body of the function is moved into `appendErr` to convince golang to inline the
	// `AppendErr`, making `err != nil` as cheap as usual.
	// Otherwise each AppendErr would waste a slow golang call just to check that err == nil and
	// return false
	haveAddedErrors = err != nil
	if haveAddedErrors {
		appendErr(d, err, summary)
	}
	return
}

// AppendErr and appendErr together can't be inlined. We're forbiding go from inlining
// appendErr into AppendErr and thus preventing AppendErr inlining.
//
//go:noinline
func appendErr(d *Diag, err error, summary string) {
	*d = append(*d, &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  summary,
		Detail:   err.Error(),
		Extra:    err,
	})
}

func FromHcl(diag *hcl.Diagnostic) Diag {
	if diag != nil {
		return Diag{diag}
	}
	return nil
}

func FromErr(err error, summary string) Diag {
	if err != nil {
		return Diag{{
			Severity: hcl.DiagError,
			Summary:  summary,
			Detail:   err.Error(),
			Extra:    err,
		}}
	}
	return nil
}

func FromErrSubj(err error, summary string, subject *hcl.Range) Diag {
	if err != nil {
		return Diag{{
			Severity: hcl.DiagError,
			Summary:  summary,
			Detail:   err.Error(),
			Extra:    err,
			Subject:  subject,
		}}
	}
	return nil
}