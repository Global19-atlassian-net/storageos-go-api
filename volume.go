package storageos

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/storageos/go-api/types"
)

var (

	// VolumeAPIPrefix is a partial path to the HTTP endpoint.
	VolumeAPIPrefix = "volumes"

	// ErrNoSuchVolume is the error returned when the volume does not exist.
	ErrNoSuchVolume = errors.New("no such volume")

	// ErrVolumeInUse is the error returned when the volume requested to be removed is still in use.
	ErrVolumeInUse = errors.New("volume in use and cannot be removed")
)

// VolumeList returns the list of available volumes.
func (c *Client) VolumeList(opts types.ListOptions) ([]types.Volume, error) {
	path, err := namespacedPath(opts.Namespace, VolumeAPIPrefix)
	if err != nil {
		return nil, err
	}
	resp, err := c.do("GET", path+"?"+queryString(opts), doOptions{context: opts.Context})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var volumes []types.Volume
	if err := json.NewDecoder(resp.Body).Decode(&volumes); err != nil {
		return nil, err
	}
	return volumes, nil
}

// Volume returns a volume by its reference.
func (c *Client) Volume(namespace string, ref string) (*types.Volume, error) {
	path, err := namespacedRefPath(namespace, VolumeAPIPrefix, ref)
	if err != nil {
		return nil, err
	}
	resp, err := c.do("GET", path, doOptions{})
	if err != nil {
		if e, ok := err.(*Error); ok && e.Status == http.StatusNotFound {
			return nil, ErrNoSuchVolume
		}
		return nil, err
	}
	defer resp.Body.Close()
	var volume types.Volume
	if err := json.NewDecoder(resp.Body).Decode(&volume); err != nil {
		return nil, err
	}
	return &volume, nil
}

// VolumeCreate creates a volume on the server and returns its unique id.
func (c *Client) VolumeCreate(opts types.VolumeCreateOptions) (*types.Volume, error) {
	path, err := namespacedPath(opts.Namespace, VolumeAPIPrefix)
	if err != nil {
		return nil, err
	}
	resp, err := c.do("POST", path, doOptions{
		data:    opts,
		context: opts.Context,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// out, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return "", err
	// }
	// return strconv.Unquote(string(out))
	var volume types.Volume
	if err := json.NewDecoder(resp.Body).Decode(&volume); err != nil {
		return nil, err
	}
	return &volume, nil
}

// VolumeUpdate updates a volume on the server.
func (c *Client) VolumeUpdate(opts types.VolumeUpdateOptions) (*types.Volume, error) {
	ref := opts.Name
	if IsUUID(opts.ID) {
		ref = opts.ID
	}
	path, err := namespacedRefPath(opts.Namespace, VolumeAPIPrefix, ref)
	if err != nil {
		return nil, err
	}
	resp, err := c.do("PUT", path, doOptions{
		data:    opts,
		context: opts.Context,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var volume types.Volume
	if err := json.NewDecoder(resp.Body).Decode(&volume); err != nil {
		return nil, err
	}
	return &volume, nil
}

// VolumeDelete removes a volume by its reference.
func (c *Client) VolumeDelete(namespace string, ref string) error {
	path, err := namespacedRefPath(namespace, VolumeAPIPrefix, ref)
	if err != nil {
		return err
	}
	resp, err := c.do("DELETE", path, doOptions{})
	if err != nil {
		if e, ok := err.(*Error); ok {
			if e.Status == http.StatusNotFound {
				return ErrNoSuchVolume
			}
			if e.Status == http.StatusConflict {
				return ErrVolumeInUse
			}
		}
		return nil
	}
	defer resp.Body.Close()
	return nil
}
