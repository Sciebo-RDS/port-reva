/* **********************************
 * Date: 2021-01-11
 * *********************************/

package cmdline

import "flag"

// IsFlagSet checks whether the given flag was specified via command-line.
func IsFlagSet(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
