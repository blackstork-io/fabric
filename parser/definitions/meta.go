package definitions

type MetaBlock struct {
	Name        string   `hcl:"name,optional"`
	Description string   `hcl:"description,optional"`
	Url         string   `hcl:"url,optional"`
	License     string   `hcl:"license,optional"`
	Author      string   `hcl:"author,optional"`
	Tags        []string `hcl:"tags,optional"`
	UpdatedAt   string   `hcl:"updated_at,optional"`

	// TODO: ?store def range defRange hcl.Range
}
