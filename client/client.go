package client

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/zzztimbo/cq-source-middesk/internal/middesk"
)

type Client struct {
	logger        zerolog.Logger
	Spec          Spec
	MiddeskClient *middesk.Client
}

func (c *Client) ID() string {
	// TODO: Change to either your plugin name or a unique dynamic identifier
	return "middesk"
}

func (c *Client) Logger() *zerolog.Logger {
	return &c.logger
}

func New(ctx context.Context, logger zerolog.Logger, s *Spec) (Client, error) {
	middeskClient, err := middesk.NewClient(s.API_KEY)

	if err != nil {
		return Client{}, err
	}

	c := Client{
		logger:        logger,
		Spec:          *s,
		MiddeskClient: middeskClient,
	}

	return c, nil
}
