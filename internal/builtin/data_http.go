package builtin

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/internal/builtin/utils"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/constraint"
)

func makeHTTPDataSource(version string) *plugin.DataSource {
	return &plugin.DataSource{
		DataFunc: fetchHTTPDataWrapper(version),
		Args: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name:        "url",
				Type:        cty.String,
				ExampleVal:  cty.StringVal("https://example.localhost/file.json"),
				Constraints: constraint.RequiredMeaningful,
				Doc:         "URL to fetch data from. Supported schemas are `http` and `https`",
			},
			&dataspec.AttrSpec{
				Name:       "method",
				Type:       cty.String,
				DefaultVal: cty.StringVal("GET"),
				OneOf: []cty.Value{
					cty.StringVal("GET"),
					cty.StringVal("POST"),
					cty.StringVal("HEAD"),
				},
				Doc: "HTTP method for the request. Allowed methods are `GET`, `POST` and `HEAD`",
			},
			&dataspec.AttrSpec{
				Name:       "basic_auth_username",
				Type:       cty.String,
				DefaultVal: cty.NullVal(cty.String),
				Doc:        `A basic authentication username`,
			},
			&dataspec.AttrSpec{
				Name:       "basic_auth_password",
				Type:       cty.String,
				Secret:     true,
				DefaultVal: cty.NullVal(cty.String),
				Doc:        `A basic authentication password`,
			},
			&dataspec.AttrSpec{
				Name:       "insecure",
				Type:       cty.Bool,
				DefaultVal: cty.BoolVal(false),
				Doc:        "If set to `true`, disabled verification of the server's certificate.",
			},
			&dataspec.AttrSpec{
				Name:        "timeout_secs",
				Type:        cty.Number,
				DefaultVal:  cty.NumberIntVal(15),
				Constraints: constraint.NonNull,
				Doc:         "The timeout for the request in seconds. Accepts fractional values.",
			},
			&dataspec.AttrSpec{
				Name:       "headers",
				Type:       cty.Map(cty.String),
				DefaultVal: cty.NullVal(cty.Map(cty.String)),
				Doc:        `The headers to be set in a request`,
			},
			&dataspec.AttrSpec{
				Name:       "body",
				Type:       cty.Map(cty.String),
				DefaultVal: cty.NullVal(cty.Map(cty.String)),
				Doc:        `Request body`,
			},
		},
		Doc: `
		Loads data from a URL.

		At the moment, the data source accepts only responses with UTF-8 charset and parses only responses with MIME types ` + "`text/csv`" + ` or ` + "`application/json`" + `. 
		If MIME type of the response is ` + "`text/csv`" + ` or ` + "`application/json`" + `, the response content will be parsed and returned as a JSON structure (similar to the behaviour of CSV and JSON data sources). Otherwise, the response content will be returned as text`,
	}
}

func StringPtr(s string) *string {
	return &s
}

type Request struct {
	Url               string
	Method            string
	Timeout           time.Duration
	SkipVerify        bool
	Headers           map[string]string
	Body              *string
	BasicAuthUsername *string
	BasicAuthPassword *string
}

type Response struct {
	Body     []byte
	MimeType string
}

func SendRequest(ctx context.Context, r *Request) (*Response, error) {
	var u *url.URL
	var err error

	u, err = url.Parse(r.Url)
	if err != nil {
		return nil, err
	}

	var reqBody io.Reader
	if r.Body != nil {
		reqBody = strings.NewReader(*r.Body)
	}

	request, err := http.NewRequestWithContext(ctx, r.Method, u.String(), reqBody)
	if err != nil {
		return nil, err
	}

	if r.BasicAuthUsername != nil && r.BasicAuthPassword != nil {
		request.SetBasicAuth(*r.BasicAuthUsername, *r.BasicAuthPassword)
	}

	if r.Headers != nil {
		for k, v := range r.Headers {
			request.Header.Set(k, v)
		}
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: r.SkipVerify}
	client := &http.Client{Transport: transport, Timeout: r.Timeout}

	slog.Debug("Sending a HTTP request", "url", r.Url, "method", r.Method, "insecure", r.SkipVerify)

	res, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("the server responded with status code %d", res.StatusCode)
	}

	contentType := res.Header.Get("Content-Type")
	var mimeType string
	if contentType == "" {
		mimeType = "text/plain" // assume `text/plain` if no content type set
	} else {
		mimeType, _, err = mime.ParseMediaType(contentType)
		if err != nil {
			return nil, err
		}
	}

	defer res.Body.Close()

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading response body: %s", err)
	}

	if !utf8.Valid(bytes) {
		return nil, fmt.Errorf("response body is not recognized as UTF-8: %s", err)
	}
	return &Response{Body: bytes, MimeType: mimeType}, nil
}

func fetchHTTPDataWrapper(version string) plugin.RetrieveDataFunc {
	return func(ctx context.Context, params *plugin.RetrieveDataParams) (plugin.Data, diagnostics.Diag) {
		return fetchHTTPData(ctx, params, version)
	}
}

func fetchHTTPData(ctx context.Context, params *plugin.RetrieveDataParams, version string) (plugin.Data, diagnostics.Diag) {
	url := params.Args.GetAttr("url").AsString()
	method := params.Args.GetAttr("method").AsString()
	authUsername := params.Args.GetAttr("basic_auth_username")
	authPassword := params.Args.GetAttr("basic_auth_username")
	insecure := params.Args.GetAttr("insecure").True()
	timeoutSecs, _ := params.Args.GetAttr("timeout_secs").AsBigFloat().Float64()
	headers := params.Args.GetAttr("headers")
	body := params.Args.GetAttr("body")

	var req = Request{
		Url:        url,
		Method:     method,
		Timeout:    time.Duration(timeoutSecs * float64(time.Second)),
		SkipVerify: insecure,
		Headers:    make(map[string]string),
		Body:       nil,
	}

	if !authUsername.IsNull() && authUsername.AsString() != "" {
		req.BasicAuthUsername = StringPtr(authUsername.AsString())
	}

	if !authPassword.IsNull() && authPassword.AsString() != "" {
		req.BasicAuthPassword = StringPtr(authPassword.AsString())
	}

	if !headers.IsNull() {
		for k, v := range headers.AsValueMap() {
			req.Headers[k] = v.AsString()
		}
	}

	if !body.IsNull() && body.AsString() != "" {
		req.Body = StringPtr(body.AsString())
	}

	req.Headers["User-Agent"] = fmt.Sprintf("fabric-data-http/%s", version)

	response, err := SendRequest(ctx, &req)
	if err != nil {
		return nil, diagnostics.Diag{
			{
				Severity: hcl.DiagError,
				Summary:  "Failed to fetch data with HTTP request",
				Detail:   err.Error(),
			},
		}
	}
	slog.Debug("Response received", "mime_type", response.MimeType, "body_bytes_count", len(response.Body))

	var result plugin.Data

	if response.MimeType == "text/csv" {
		reader := csv.NewReader(bytes.NewBuffer(response.Body))
		reader.Comma = ',' // Use `,` as a delimiter by default
		result, err = utils.ParseCSVContent(ctx, reader)
		if err != nil {
			return nil, diagnostics.Diag{
				{
					Severity: hcl.DiagError,
					Summary:  "Failed to parse CSV content",
					Detail:   err.Error(),
				},
			}
		}
	} else if response.MimeType == "application/json" {
		result, err = plugin.UnmarshalJSONData(response.Body)
		if err != nil {
			return nil, diagnostics.Diag{
				{
					Severity: hcl.DiagError,
					Summary:  "Failed to parse JSON content",
					Detail:   err.Error(),
				},
			}
		}
	} else {
		result = plugin.StringData(response.Body)
	}
	return result, nil
}
