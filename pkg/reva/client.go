/* **********************************
 * Date: 2021-01-07
 * *********************************/

package reva

import (
	"fmt"
	"path"
	"strings"

	storage "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/pkg/sdk"
	"github.com/cs3org/reva/pkg/sdk/action"
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
		return errors.Wrapf(err, "unable to login (u=%v)", username)
	}

	return nil
}

func (cl *Client) ListFolder(filePath string) ([]string, error) {
	if !cl.session.IsValid() {
		return []string{}, errors.Errorf("no valid session")
	}

	filePath = cl.cleanPath(filePath)

	enumFiles, err := action.NewEnumFilesAction(cl.session)
	if err != nil {
		return []string{}, errors.Wrap(err, "unable to create file enumeration action")
	}

	files, err := enumFiles.ListAll(filePath, false)
	if err != nil {
		return []string{}, errors.Wrap(err, fmt.Sprintf("unable to enumerate files in %v", filePath))
	}

	folderContents := make([]string, 0, len(files))
	for _, file := range files {
		// Remove the '/home' root dir
		// TODO: Remove this later
		fileName := strings.TrimPrefix(file.Path, "/home")

		// Ensure that folders end with a slash
		if file.Type == storage.ResourceType_RESOURCE_TYPE_CONTAINER && !strings.HasSuffix(fileName, "/") {
			fileName += "/"
		}

		folderContents = append(folderContents, fileName)
	}

	return folderContents, nil
}

func (cl *Client) DownloadFile(filePath string) ([]byte, error) {
	if !cl.session.IsValid() {
		return []byte{}, errors.Errorf("no valid session")
	}

	filePath = cl.cleanPath(filePath)

	downloadFile, err := action.NewDownloadAction(cl.session)
	if err != nil {
		return []byte{}, errors.Wrap(err, "unable to create download action")
	}

	fileData, err := downloadFile.DownloadFile(filePath)
	if err != nil {
		return []byte{}, errors.Wrap(err, fmt.Sprintf("unable to download file %v", filePath))
	}

	return fileData, nil
}

func (cl *Client) cleanPath(filePath string) string {
	// Clean the file path and prepend the '/home' root dir
	filePath = path.Clean(filePath)
	// TODO: Remove this later
	if !strings.HasPrefix(filePath, "/home") {
		filePath = path.Join("/home", filePath)
	}
	return filePath
}

// New creates a new Client instance.
func New(host string, username string, password string, log *zerolog.Logger) (*Client, error) {
	cl := &Client{}
	if err := cl.initialize(host, username, password, log); err != nil {
		return nil, errors.Wrap(err, "unable to initialize the Reva client")
	}
	return cl, nil
}
