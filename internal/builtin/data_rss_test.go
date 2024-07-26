package builtin

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/plugindata"
	"github.com/blackstork-io/fabric/plugin/plugintest"
)

func Test_makeRSSDataSchema(t *testing.T) {
	t.Parallel()
	t.Run("basic", func(t *testing.T) {
		t.Parallel()

		assert := assert.New(t)
		schema := makeRSSDataSource()
		assert.Nil(schema.Config)
		assert.NotNil(schema.Args)
		assert.NotNil(schema.DataFunc)
	})
}

func makeTestRssServer() (baseAddr string, closer func()) {
	mux := http.NewServeMux()
	mux.HandleFunc("/up", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	mux.Handle("/basic-auth-redir", http.RedirectHandler(baseAddr+"basic-auth", http.StatusMovedPermanently))
	mux.HandleFunc("/basic-auth", func(w http.ResponseWriter, r *http.Request) {
		un, pw, ok := r.BasicAuth()
		if !ok || un != "user" || pw != "pass" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		r.RequestURI = strings.ReplaceAll(r.RequestURI, "/basic-auth", "/data/basic.rss")
		var err error
		r.URL, err = url.ParseRequestURI(r.RequestURI)
		if err != nil {
			panic(err)
		}
		mux.ServeHTTP(w, r)
	})

	path, err := filepath.Abs("./testdata/rss/")
	if err != nil {
		panic(err)
	}
	mux.Handle("/data/", http.StripPrefix("/data/", http.FileServerFS(os.DirFS(path))))

	srv := httptest.NewServer(mux)
	close := srv.Close
	defer func() {
		close()
	}()
	baseAddr = srv.URL + "/"
	resp, err := http.Get(baseAddr + "up")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if http.StatusOK != resp.StatusCode {
		panic("resp status code is wrone")
	}

	closer = close
	close = func() {}
	return
}

func Test_fetchRSSData(t *testing.T) {
	t.Parallel()

	addr, close := makeTestRssServer()
	defer close()

	type result struct {
		Data  plugindata.Data
		Diags diagnostics.Diag
	}

	ValidRssData := plugindata.Map{
		"description":   plugindata.String("This is an example of an RSS feed"),
		"link":          plugindata.String("http://www.example.com/main.html"),
		"pub_date":      plugindata.String("Sun, 6 Sep 2009 16:20:12 +0000"),
		"title":         plugindata.String("RSS Title"),
		"pub_timestamp": plugindata.Number(1252254012),
		"items": plugindata.List{
			plugindata.Map{
				"description":   plugindata.String("Here is some text containing an interesting description."),
				"guid":          plugindata.String("7bd204c6-1655-4c27-aeee-53f933c5395f"),
				"link":          plugindata.String("http://www.example.com/blog/post/1"),
				"pub_date":      plugindata.String("Sun, 6 Sep 2009 16:20:23 +0000"),
				"title":         plugindata.String("Example entry"),
				"pub_timestamp": plugindata.Number(1252254023),
			},
		},
	}

	ValidAtomData := plugindata.Map{
		"description": plugindata.String("A subtitle."),
		"link":        plugindata.String("http://example.org/"),
		"pub_date":    plugindata.String(""),
		"title":       plugindata.String("Example Feed"),
		"items": plugindata.List{
			plugindata.Map{
				"description":   plugindata.String("Some text."),
				"guid":          plugindata.String("urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a"),
				"link":          plugindata.String("http://example.org/2003/12/13/atom03"),
				"pub_date":      plugindata.String("2003-11-09T17:23:02Z"),
				"pub_timestamp": plugindata.Number(1068398582),
				"title":         plugindata.String("Atom-Powered Robots Run Amok"),
			},
		},
	}

	tt := []struct {
		name     string
		url      string
		auth     *gofeed.Auth
		expected result
	}{
		{
			name: "valid_rss",
			url:  "data/basic.rss",
			expected: result{
				Data: ValidRssData,
			},
		},
		{
			name: "valid_atom",
			url:  "data/basic.atom",
			expected: result{
				Data: ValidAtomData,
			},
		},
		{
			name: "valid_auth",
			url:  "basic-auth",
			auth: &gofeed.Auth{
				Username: "user",
				Password: "pass",
			},

			expected: result{
				Data: ValidRssData,
			},
		},
		{
			name: "invalid_auth",
			url:  "basic-auth",
			auth: &gofeed.Auth{
				Username: "invalid_user",
				Password: "invalid_pass",
			},

			expected: result{
				Diags: diagnostics.Diag{{
					Severity: hcl.DiagError,
					Summary:  "Failed to fetch the feed",
					Detail:   "http error: 401 Unauthorized",
				}},
			},
		},
		{
			name: "empty_auth",
			url:  "basic-auth",
			auth: &gofeed.Auth{},
			expected: result{
				Diags: diagnostics.Diag{{
					Severity: hcl.DiagError,
					Summary:  "Failed to fetch the feed",
					Detail:   "http error: 401 Unauthorized",
				}},
			},
		},
		{
			name: "incomplete_auth",
			url:  "basic-auth",
			auth: &gofeed.Auth{
				Username: "user",
			},
			expected: result{
				Diags: diagnostics.Diag{{
					Severity: hcl.DiagError,
					Summary:  "Failed to fetch the feed",
					Detail:   "http error: 401 Unauthorized",
				}},
			},
		},
		{
			name: "absent_auth",
			url:  "basic-auth",
			expected: result{
				Diags: diagnostics.Diag{{
					Severity: hcl.DiagError,
					Summary:  "Failed to fetch the feed",
					Detail:   "http error: 401 Unauthorized",
				}},
			},
		},
		{
			name: "valid_auth_redir",
			url:  "basic-auth-redir",
			auth: &gofeed.Auth{
				Username: "user",
				Password: "pass",
			},
			expected: result{
				Data: ValidRssData,
			},
		},
		{
			name: "invalid_path",
			url:  "does_not_exist.rss",
			expected: result{
				Diags: diagnostics.Diag{{
					Severity: hcl.DiagError,
					Summary:  "Failed to fetch the feed",
					Detail:   "http error: 404 Not Found",
				}},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)

			p := &plugin.Schema{
				DataSources: plugin.DataSources{
					"rss": makeRSSDataSource(),
				},
			}

			dec := plugintest.NewTestDecoder(t, p.DataSources["rss"].Args).
				SetAttr("url", cty.StringVal(addr+tc.url))

			if tc.auth != nil {
				dec.AppendBody(fmt.Sprintf(`
					basic_auth {
						username = "%s"
						password = "%s"
					}
					`, tc.auth.Username, tc.auth.Password))
			}

			params := &plugin.RetrieveDataParams{Args: dec.Decode()}

			data, diags := p.RetrieveData(context.Background(), "rss", params)
			assert.Equal(tc.expected, result{data, diags})
		})
	}
}
