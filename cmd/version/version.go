/* **********************************
 * Date: 2021-01-06
 * *********************************/

package version

import "fmt"

const (
	Major    = 0
	Minor    = 1
	Revision = 0
)

// GetString returns a formatted version number string.
func GetString() string {
	return fmt.Sprintf("%v.%v.%v", Major, Minor, Revision)
}
