package dnsimple

import (
	"context"
	"fmt"
)

type SecondaryService struct {
	client *Client
}

// SecondaryServer represents a primary DNS server within the Secondary DNS service in DNSimple.
type SecondaryServer struct {
	ID           int64  `json:"id,omitempty"`
	AccountID    int64  `json:"account_id,omitempty"`
	Name         string `json:"name,omitempty"`
	IP           string `json:"ip,omitempty"`
	Port         uint64 `json:"port"`
	CreatedAt    string `json:"created_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
}

func secondaryDNSPrimary(accountID string, serverIdentifier string) (path string) {
	path = fmt.Sprintf("/%v/secondary_dns/primaries", accountID)
	if serverIdentifier != "" {
		path += fmt.Sprintf("/%v", serverIdentifier)
	}
	return
}

// SecondaryServerResponse represents a response from an API method that returns a SecondaryServer struct.
type SecondaryServerResponse struct {
	Response
	Data *SecondaryServer `json:"data"`
}

// SecondaryServersResponse represents a response from an API method that returns a collection of SecondaryServer struct.
type SecondaryServersResponse struct {
	Response
	Data []SecondaryServer `json:"data"`
}

// CreatePrimaryServer creates a new primary server in the account.
//
// See https://developer.dnsimple.com/v2/secondary-dns/#createPrimaryServer
func (s *SecondaryService) CreatePrimaryServer(ctx context.Context, accountID string, serverAttributes SecondaryServer) (*SecondaryServerResponse, error) {
	path := versioned(secondaryDNSPrimary(accountID, ""))
	serverResponse := &SecondaryServerResponse{}

	resp, err := s.client.post(ctx, path, serverAttributes, serverResponse)
	if err != nil {
		return nil, err
	}

	serverResponse.HTTPResponse = resp
	return serverResponse, nil
}

// SecondaryServerListOptions specifies the optional parameters you can provide
// to customize the SecondaryService.ListPrimaryServers method.
type SecondaryServerListOptions struct {
        ListOptions
}

// ListPrimaryServers
//
// See https://developer.dnsimple.com/v2/secondary-dns/#listPrimaryServers
func (s *SecondaryService) ListPrimaryServers(ctx context.Context, accountID string, options *SecondaryServerListOptions) (*SecondaryServersResponse, error) {
	path := versioned(secondaryDNSPrimary(accountID, ""))
	serversResponse := &SecondaryServersResponse{}

        path, err := addURLQueryOptions(path, options)
        if err != nil {
                return nil, err
        }

        resp, err := s.client.get(ctx, path, serversResponse)
        if err != nil {
                return nil, err
        }

        serversResponse.HTTPResponse = resp
        return serversResponse, nil
}

// DeletePrimaryServer
//
// See https://developer.dnsimple.com/v2/secondary-dns/#removePrimaryServer
func (s *SecondaryService) DeletePrimaryServer(ctx context.Context, accountID string, serverIdentifier string) (*SecondaryServerResponse, error) {
        path := versioned(secondaryDNSPrimary(accountID, serverIdentifier))
        serverResponse := &SecondaryServerResponse{}

        resp, err := s.client.delete(ctx, path, nil, nil)
        if err != nil {
                return nil, err
        }

        serverResponse.HTTPResponse = resp
        return serverResponse, nil
}
