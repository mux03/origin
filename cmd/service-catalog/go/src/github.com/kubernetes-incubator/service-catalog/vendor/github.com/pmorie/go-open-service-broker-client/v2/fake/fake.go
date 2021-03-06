package fake

import (
	"errors"
	"net/http"
	"sync"

	"github.com/pmorie/go-open-service-broker-client/v2"
)

// NewFakeClientFunc returns a v2.CreateFunc that returns a FakeClient with
// the given FakeClientConfiguration.  It is useful for injecting the
// FakeClient in code that uses the v2.CreateFunc interface.
func NewFakeClientFunc(config FakeClientConfiguration) v2.CreateFunc {
	return func(_ *v2.ClientConfiguration) (v2.Client, error) {
		return NewFakeClient(config), nil
	}
}

// ReturnFakeClientFunc returns a v2.CreateFunc that returns the given
// FakeClient.
func ReturnFakeClientFunc(c *FakeClient) v2.CreateFunc {
	return func(_ *v2.ClientConfiguration) (v2.Client, error) {
		return c, nil
	}
}

// NewFakeClient returns a new fake Client with the given
// FakeClientConfiguration.
func NewFakeClient(config FakeClientConfiguration) *FakeClient {
	return &FakeClient{
		CatalogReaction:           config.CatalogReaction,
		ProvisionReaction:         config.ProvisionReaction,
		UpdateInstanceReaction:    config.UpdateInstanceReaction,
		DeprovisionReaction:       config.DeprovisionReaction,
		PollLastOperationReaction: config.PollLastOperationReaction,
		BindReaction:              config.BindReaction,
		UnbindReaction:            config.UnbindReaction,
	}
}

// FakeClientConfiguration models the configuration of a FakeClient.
type FakeClientConfiguration struct {
	CatalogReaction           *CatalogReaction
	ProvisionReaction         *ProvisionReaction
	UpdateInstanceReaction    *UpdateInstanceReaction
	DeprovisionReaction       *DeprovisionReaction
	PollLastOperationReaction *PollLastOperationReaction
	BindReaction              *BindReaction
	UnbindReaction            *UnbindReaction
}

// Action is a record of a method call on the FakeClient.
type Action struct {
	Type    ActionType
	Request interface{}
}

// ActionType is a typedef over the set of actions that can be taken on a
// FakeClient.
type ActionType string

const (
	GetCatalog          ActionType = "GetCatalog"
	ProvisionInstance   ActionType = "ProvisionInstance"
	UpdateInstance      ActionType = "UpdateInstance"
	DeprovisionInstance ActionType = "DeprovisionInstance"
	PollLastOperation   ActionType = "PollLastOperation"
	Bind                ActionType = "Bind"
	Unbind              ActionType = "Unbind"
)

// FakeClient is a fake implementation of the v2.Client interface. It records
// the actions that are taken on it and runs the appropriate reaction to those
// actions. If an action for which there is no reaction specified occurs, it
// returns an error.  FakeClient is threadsafe.
type FakeClient struct {
	CatalogReaction           *CatalogReaction
	ProvisionReaction         *ProvisionReaction
	UpdateInstanceReaction    *UpdateInstanceReaction
	DeprovisionReaction       *DeprovisionReaction
	PollLastOperationReaction *PollLastOperationReaction
	BindReaction              *BindReaction
	UnbindReaction            *UnbindReaction

	sync.Mutex
	actions []Action
}

var _ v2.Client = &FakeClient{}

func (c *FakeClient) Actions() []Action {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	return c.actions
}

func (c *FakeClient) GetCatalog() (*v2.CatalogResponse, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	c.actions = append(c.actions, Action{Type: GetCatalog})

	if c.CatalogReaction != nil {
		return c.CatalogReaction.Response, c.CatalogReaction.Error
	}

	return nil, UnexpectedActionError()
}

func (c *FakeClient) ProvisionInstance(r *v2.ProvisionRequest) (*v2.ProvisionResponse, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	c.actions = append(c.actions, Action{ProvisionInstance, r})

	if c.ProvisionReaction != nil {
		return c.ProvisionReaction.Response, c.ProvisionReaction.Error
	}

	return nil, UnexpectedActionError()
}

func (c *FakeClient) UpdateInstance(r *v2.UpdateInstanceRequest) (*v2.UpdateInstanceResponse, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	c.actions = append(c.actions, Action{UpdateInstance, r})

	if c.UpdateInstanceReaction != nil {
		return c.UpdateInstanceReaction.Response, c.UpdateInstanceReaction.Error
	}

	return nil, UnexpectedActionError()
}

func (c *FakeClient) DeprovisionInstance(r *v2.DeprovisionRequest) (*v2.DeprovisionResponse, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	c.actions = append(c.actions, Action{DeprovisionInstance, r})

	if c.DeprovisionReaction != nil {
		return c.DeprovisionReaction.Response, c.DeprovisionReaction.Error
	}

	return nil, UnexpectedActionError()
}

func (c *FakeClient) PollLastOperation(r *v2.LastOperationRequest) (*v2.LastOperationResponse, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	c.actions = append(c.actions, Action{PollLastOperation, r})

	if c.PollLastOperationReaction != nil {
		return c.PollLastOperationReaction.Response, c.PollLastOperationReaction.Error
	}

	return nil, UnexpectedActionError()
}

func (c *FakeClient) Bind(r *v2.BindRequest) (*v2.BindResponse, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	c.actions = append(c.actions, Action{Bind, r})

	if c.BindReaction != nil {
		return c.BindReaction.Response, c.BindReaction.Error
	}

	return nil, UnexpectedActionError()
}

func (c *FakeClient) Unbind(r *v2.UnbindRequest) (*v2.UnbindResponse, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	c.actions = append(c.actions, Action{Unbind, r})

	if c.UnbindReaction != nil {
		return c.UnbindReaction.Response, c.UnbindReaction.Error
	}

	return nil, UnexpectedActionError()
}

func UnexpectedActionError() error {
	return errors.New("Unexpected action")
}

type CatalogReaction struct {
	Response *v2.CatalogResponse
	Error    error
}

type ProvisionReaction struct {
	Response *v2.ProvisionResponse
	Error    error
}

type UpdateInstanceReaction struct {
	Response *v2.UpdateInstanceResponse
	Error    error
}

type DeprovisionReaction struct {
	Response *v2.DeprovisionResponse
	Error    error
}

type PollLastOperationReaction struct {
	Response *v2.LastOperationResponse
	Error    error
}

type BindReaction struct {
	Response *v2.BindResponse
	Error    error
}

type UnbindReaction struct {
	Response *v2.UnbindResponse
	Error    error
}

func strPtr(s string) *string {
	return &s
}

// AsyncRequiredError returns error for required asynchronous operations.
func AsyncRequiredError() error {
	return v2.HTTPStatusCodeError{
		StatusCode:   http.StatusUnprocessableEntity,
		ErrorMessage: strPtr(v2.AsyncErrorMessage),
		Description:  strPtr(v2.AsyncErrorDescription),
	}
}

// AppGUIDRequiredError returns error for when app GUID is missing from bind
// request.
func AppGUIDRequiredError() error {
	return v2.HTTPStatusCodeError{
		StatusCode:   http.StatusUnprocessableEntity,
		ErrorMessage: strPtr(v2.AppGUIDRequiredErrorMessage),
		Description:  strPtr(v2.AppGUIDRequiredErrorDescription),
	}
}
