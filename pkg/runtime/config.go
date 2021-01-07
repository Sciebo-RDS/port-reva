/* **********************************
 * Date: 2021-01-07
 * *********************************/

package runtime

// Config holds the runtime configuration.
type Config struct {
	WebserverPort uint16

	Reva struct {
		Host     string
		User     string
		Password string
	}
}
