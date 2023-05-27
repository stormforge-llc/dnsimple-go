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
	LinkedSecondaryZones []string `json:"linked_secondary_zones"`
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

func secondaryDNSZone(accountID string, zoneIdentifier string) (path string) {
	path = fmt.Sprintf("/%v/secondary_dns/zones", accountID)
	if zoneIdentifier != "" {
		path += fmt.Sprintf("/%v", zoneIdentifier)
	}
	return
}

func secondaryDNSLink(accountID string, serverIdentifier string) (path string) {
	return fmt.Sprintf("/%v/secondary_dns/primaries/%v/link", accountID, serverIdentifier)
}

func secondaryDNSUnlink(accountID string, serverIdentifier string) (path string) {
	return fmt.Sprintf("/%v/secondary_dns/primaries/%v/unlink", accountID, serverIdentifier)
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

// GetPrimaryServer
//
// See https://developer.dnsimple.com/v2/secondary-dns/#getPrimaryServer
func (s *SecondaryService) GetPrimaryServer(ctx context.Context, accountID string, serverIdentifier string) (*SecondaryServerResponse, error) {
	path := versioned(secondaryDNSPrimary(accountID, serverIdentifier))
	serverResponse := &SecondaryServerResponse{}

        resp, err := s.client.get(ctx, path, serverResponse)
        if err != nil {
                return nil, err
        }

        serverResponse.HTTPResponse = resp
        return serverResponse, nil
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

// SecondaryZone
type SecondaryZone struct {
	Secondary bool `json:"secondary"`
	LastTransferredAt string `json:"last_transferred_at,omitempty"`

	Zone
}

// SecondaryZoneResponse represents a response from an API method that returns a SecondaryZone struct.
type SecondaryZoneResponse struct {
	Response
	Data *SecondaryZone `json:"data"`
}

// CreateSecondaryZone creates a secondary zone for the account.
//
// See https://developer.dnsimple.com/v2/secondary-dns/#createSecondaryZone
func (s *SecondaryService) CreateSecondaryZone(ctx context.Context, accountID string, secondaryZoneAttributes SecondaryZone) (*SecondaryZoneResponse, error) {
	path := versioned(secondaryDNSZone(accountID, ""))
	zoneResponse := &SecondaryZoneResponse{}

	resp, err := s.client.post(ctx, path, secondaryZoneAttributes, zoneResponse)
	if err != nil {
		return nil, err
	}

	zoneResponse.HTTPResponse = resp
	return zoneResponse, nil
}

type LinkPrimaryServerToSecondaryZoneRequest struct {
	ZoneIdentifier string `json:"zone"`
}

// LinkPrimaryServerToSecondaryZone
//
// See https://developer.dnsimple.com/v2/secondary-dns/#linkPrimaryServer
func (s *SecondaryService) LinkPrimaryServerToSecondaryZone(ctx context.Context, accountID string, serverIdentifier string, secondaryIdentifier string) (*SecondaryServerResponse, error) {
	path := versioned(secondaryDNSLink(accountID, serverIdentifier))
	linkReq := &LinkPrimaryServerToSecondaryZoneRequest{
		ZoneIdentifier: secondaryIdentifier,
	}
	serverResponse := &SecondaryServerResponse{}

	resp, err := s.client.put(ctx, path, linkReq, serverResponse)
	if err != nil {
		return nil, err
	}

	serverResponse.HTTPResponse = resp
	return serverResponse, nil
}

// UnlinkPrimaryServerToSecondaryZone
//
// See https://developer.dnsimple.com/v2/secondary-dns/#unlinkPrimaryServer
func (s *SecondaryService) UnlinkPrimaryServerToSecondaryZone(ctx context.Context, accountID string, serverIdentifier string, secondaryIdentifier string) (*SecondaryServerResponse, error) {
	path := versioned(secondaryDNSUnlink(accountID, serverIdentifier))
	unlinkReq := &LinkPrimaryServerToSecondaryZoneRequest{
		ZoneIdentifier: secondaryIdentifier,
	}
	serverResponse := &SecondaryServerResponse{}

	resp, err := s.client.put(ctx, path, unlinkReq, serverResponse)
	if err != nil {
		return nil, err
	}

	serverResponse.HTTPResponse = resp
	return serverResponse, nil
}
