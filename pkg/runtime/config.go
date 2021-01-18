/* **********************************
 * Date: 2021-01-07
 * *********************************/

package runtime

import "github.com/Sciebo-RDS/port-reva/pkg/reva"

// Config holds the runtime configuration.
type Config struct {
	WebserverPort uint16

	Reva reva.Config
}
