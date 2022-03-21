package service

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/nayyara-cropsey/jwtmock"
)

// ClientRepo is a repo for storing registred clients and generating tokens for clients.
type ClientRepo struct {
	clients map[string]*jwtmock.ClientRegistration

	m sync.Mutex
}

// NewClientRepo is the preferred way to instantiate a client repo.
func NewClientRepo() *ClientRepo {
	return &ClientRepo{
		clients: make(map[string]*jwtmock.ClientRegistration),
	}
}

// Register registers a new client and returns any errors
func (c *ClientRepo) Register(registration jwtmock.ClientRegistration) error {
	c.m.Lock()
	defer c.m.Unlock()

	if _, ok := c.clients[registration.ID]; ok {
		return errors.New("duplicate client registration")
	}

	c.clients[registration.ID] = &registration

	return nil
}

// GenerateToken generates a token response from the given client request
func (c *ClientRepo) GenerateToken(request jwtmock.ClientTokenRequest,
	key *jwtmock.SigningKey) (*jwtmock.ClientTokenResponse, error) {
	client, ok := c.clients[request.ClientID]
	if !ok {
		return nil, errors.New("client does not exist")
	}

	if client.Secret != request.ClientSecret {
		return nil, errors.New("client secret is wrong")
	}

	if request.GrantType != jwtmock.ClientCredentials {
		return nil, errors.New("invalid grant type")
	}

	exp := time.Now().Add(time.Hour).Unix()
	claims, err := jwtmock.ClaimsFrom(jwtmock.ClientTokenClaims{
		Issuer:          "https://jwtmock.co",
		Subject:         fmt.Sprintf("%v@clients", client.ID),
		Audience:        request.Audience,
		IssuedAt:        time.Now().Unix(),
		Expires:         exp,
		AuthorizedParty: client.ID,
		Scope:           client.Scope,
		GrantType:       jwtmock.ClientCredentials,
	})
	if err != nil {
		return nil, fmt.Errorf("claims generation: %w", err)
	}

	token, err := claims.CreateJWT(key)
	if err != nil {
		return nil, fmt.Errorf("JWT generation: %w", err)
	}

	return &jwtmock.ClientTokenResponse{
		AccessToken: token,
		Scope:       client.Scope,
		ExpiresIn:   exp,
		TokenType:   jwtmock.Bearer,
	}, nil
}
