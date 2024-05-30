package diagnostics

// More ergonomic wrapper over hcl.Diagnostics.

import (
	"github.com/hashicorp/hcl/v2"
)

// Generic type matching both hcl.Diagnostics and these Diags
type Generic interface {
	~[]*hcl.Diagnostic
}

func From[D Generic](diags D) Diag {
	return Diag(diags)
}

type Diag hcl.Diagnostics // Diagnostics does implement error interface, but not, itself, an error.

type repeatedError struct{}

// Invisible to user error, typically used to signal that the initial block evaluation
// has failed (and already has reported its errors to user).
var RepeatedError = &hcl.Diagnostic{
	Severity: hcl.DiagError,
	Extra:    repeatedError{},
}

// Attach this as an .Extra on diagnostic to enhance the error message for gojq errors.
// Diagnostic must include Subject pointing to the `attribute = "value"` line.
type GoJQError struct {
	Err   error
	Query string
}

func FindByExtra[T any](diags Diag) *hcl.Diagnostic {
	for _, diag := range diags {
		if _, found := hcl.DiagnosticExtra[T](diag); found {
			return diag
		}
	}
	return nil
}

func (d Diag) Error() string {
	return hcl.Diagnostics(d).Error()
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

// Add new diagnostic warning.
func (d *Diag) AddWarn(summary, detail string) {
	*d = append(*d, &hcl.Diagnostic{
		Severity: hcl.DiagWarning,
		Summary:  summary,
		Detail:   detail,
	})
}

// Set the subject for diagnostics if it's not already specified.
func (d *Diag) DefaultSubject(rng *hcl.Range) {
	for _, diag := range *d {
		if diag.Subject == nil {
			diag.Subject = rng
		}
	}
}

// Appends all diags to diagnostics, returns true if the just-appended diagnostics contain an error.
func (d *Diag) Extend(diags []*hcl.Diagnostic) (haveAddedErrors bool) {
	*d = append(*d, diags...)
	return hcl.Diagnostics(diags).HasErrors()
}

// HasErrors returns true if the receiver contains any diagnostics of
// severity DiagError.
func (d Diag) HasErrors() bool {
	return hcl.Diagnostics(d).HasErrors()
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

// AppendErr and appendErr together can't be inlined. We're forbidding go from inlining
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
