package storageos

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/storageos/go-api/types"
)

var (

	// UserAPIPrefix is a partial path to the HTTP endpoint.
	UserAPIPrefix = "users"

	// ErrNoSuchUser is the error returned when the user does not exist.
	ErrNoSuchUser = errors.New("no such user")
)

// UserList returns the list of available users.
func (c *Client) UserList(opts types.ListOptions) ([]*types.User, error) {
	listOpts := doOptions{
		fieldSelector: opts.FieldSelector,
		labelSelector: opts.LabelSelector,
		namespace:     opts.Namespace,
		context:       opts.Context,
	}

	if opts.LabelSelector != "" {
		query := url.Values{}
		query.Add("labelSelector", opts.LabelSelector)
		listOpts.values = query
	}

	resp, err := c.do("GET", UserAPIPrefix, listOpts)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var users struct {
		Users []*types.User `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, err
	}
	return users.Users, nil
}

// User returns a user by its username/id.
func (c *Client) User(username string) (*types.User, error) {
	path := fmt.Sprintf("%s/%s", UserAPIPrefix, username)
	resp, err := c.do("GET", path, doOptions{})
	if err != nil {
		if e, ok := err.(*Error); ok && e.Status == http.StatusNotFound {
			return nil, ErrNoSuchUser
		}
		return nil, err
	}
	defer resp.Body.Close()

	var user struct {
		User *types.User `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return user.User, nil
}

// UserCreate creates a user on the server and returns the new object.
func (c *Client) UserCreate(opts types.UserCreateOptions) (*types.User, error) {
	resp, err := c.do("POST", UserAPIPrefix, doOptions{
		data:    opts,
		context: opts.Context,
	})
	if err != nil {
		return nil, err
	}

	var user types.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

// UserUpdate updates a user on the server.
//func (c *Client) UserUpdate(opts types.UserUpdateOptions) (*types.User, error) {
//	ref := opts.Name
//	if IsUUID(opts.ID) {
//		ref = opts.ID
//	}
//	fmt.Printf("%#v\n", opts)
//	path, err := namespacedRefPath(opts.Namespace, UserAPIPrefix, ref)
//	if err != nil {
//		return nil, err
//	}
//	resp, err := c.do("PUT", path, doOptions{
//		data:    opts,
//		context: opts.Context,
//	})
//	if err != nil {
//		return nil, err
//	}
//	defer resp.Body.Close()
//	var user types.User
//	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
//		return nil, err
//	}
//	return &user, nil
//}

// UserDelete removes a user by its reference.
func (c *Client) UserDelete(opts types.DeleteOptions) error {
	resp, err := c.do("DELETE", UserAPIPrefix+"/"+opts.Name, doOptions{})
	if err != nil {
		if e, ok := err.(*Error); ok {
			if e.Status == http.StatusNotFound {
				return ErrNoSuchUser
			}
		}
		return nil
	}
	defer resp.Body.Close()
	return nil
}
