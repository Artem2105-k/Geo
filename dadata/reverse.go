package dadata

import (
	"context"
	"net/url"

	"github.com/ekomobile/dadata/v2/api/suggest"
	"github.com/ekomobile/dadata/v2/client"
	"studentgit.kata.academy/ar.konovalov202_gmail.com/rpc/general"
)

var _ general.GeoProvider = (*Client)(nil)

type Client struct {
	api *suggest.Api
}

func NewClient(apiKey, secretKey string) (*Client, error) {
	endpoint, err := url.Parse("https://suggestions.dadata.ru/suggestions/api/4_1/rs/")
	if err != nil {
		return nil, err
	}

	creds := client.Credentials{
		ApiKeyValue:    apiKey,
		SecretKeyValue: secretKey,
	}

	return &Client{
		api: &suggest.Api{
			Client: client.NewClient(endpoint, client.WithCredentialProvider(&creds)),
		},
	}, nil
}

func (c *Client) AddressSearch(query string) ([]*general.Address, error) {
	rawRes, err := c.api.Address(context.Background(), &suggest.RequestParams{
		Query: query,
		Count: 10,
	})
	if err != nil {
		return nil, err
	}
	var addresses []*general.Address
	for _, r := range rawRes {
		addresses = append(addresses, &general.Address{
			City:   r.Data.City,
			Street: r.Data.Street,
			House:  r.Data.House,
			Lat:    r.Data.GeoLat,
			Lon:    r.Data.GeoLon,
		})
	}
	return addresses, nil
}

func (c *Client) GeoCode(lat, lon string) ([]*general.Address, error) {
	rawRes, err := c.api.GeoLocate(context.Background(), &suggest.GeolocateParams{
		Lat:   lat,
		Lon:   lon,
		Count: 5,
	})
	if err != nil {
		return nil, err
	}
	var addresses []*general.Address
	for _, r := range rawRes {
		addresses = append(addresses, &general.Address{
			City:   r.Data.City,
			Street: r.Data.Street,
			House:  r.Data.House,
			Lat:    r.Data.GeoLat,
			Lon:    r.Data.GeoLon,
		})
	}
	return addresses, nil
}
