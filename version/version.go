package version

import "fmt"

var Version = Ver{
	Major: 0,
	Minor: 2,
	Patch: 2,
}

// Ver defines the version of fyntrix software.
// It follows the semantic versioning 2.0.0 spec (http://semver.org/).
type Ver struct {
	Major uint // Major version number
	Minor uint // Minor version number
	Patch uint // Patch version number
}

// String returns a string representation of the Version object.
func (v Ver) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}
