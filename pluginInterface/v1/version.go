package plugin

import "github.com/Masterminds/semver/v3"

// semver.Version doesn't implement MarshalBinary/UnmarshalBinary required for
// gob, which is required by net/rpc, which is required by go-plugin.
// This is a temporary workaround, I expect they are going to add these methods
// when I send a PR.

type Version semver.Version

func (v Version) MarshalBinary() ([]byte, error) {
	return semver.Version(v).MarshalText()
}

func (v *Version) UnmarshalBinary(data []byte) error {
	return (*semver.Version)(v).UnmarshalText(data)
}

// Returns the wrapped `semver.Version`.
func (v *Version) Cast() *semver.Version {
	return (*semver.Version)(v)
}
