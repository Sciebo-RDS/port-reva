/* **********************************
 * Date: 2021-01-07
 * *********************************/

package reva

import (
	"github.com/cs3org/reva/pkg/sdk"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Client represents the Reva client.
type Client struct {
	log *zerolog.Logger

	session *sdk.Session
}

func (cl *Client) initialize(host string, username string, password string, log *zerolog.Logger) error {
	if log == nil {
		return errors.Errorf("no logger specified")
	}
	cl.log = log

	session, err := sdk.NewSession()
	if err != nil {
		return errors.Wrap(err, "unable to create session")
	}
	cl.session = session

	if err := session.Initiate(host, false); err != nil {
		return errors.Wrapf(err, "unable to initiate session to host %v", host)
	}

	if err := session.BasicLogin(username, password); err != nil {
		return errors.Wrapf(err, "unable to login (u=%v, p=%v)", username, password)
	}

	return nil
}

// New creates a new Client instance.
func New(host string, username string, password string, log *zerolog.Logger) (*Client, error) {
	cl := &Client{}
	if err := cl.initialize(host, username, password, log); err != nil {
		return nil, errors.Wrap(err, "unable to initialize the Reva client")
	}
	return cl, nil
}
