package definitions

type MetaBlock struct {
	// XXX: is empty sting enougth or use a proper ptr-nil-if-missing?
	Author string   `hcl:"author,optional"`
	Tags   []string `hcl:"tags,optional"`
}
