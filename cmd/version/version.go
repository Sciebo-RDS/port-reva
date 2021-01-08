/* **********************************
 * Date: 2021-01-06
 * *********************************/

package version

import "fmt"

const (
	Major    = 0
	Minor    = 3
	Revision = 0
)

// GetString returns a formatted version number string.
func GetString() string {
	return fmt.Sprintf("%v.%v.%v", Major, Minor, Revision)
}
