package bigquery

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type Client struct {
	client *bigquery.Client
}

type QueryOption func(*bigquery.Query)

func WithQueryParameter(name string, value interface{}) QueryOption {
	return func(q *bigquery.Query) {
		q.Parameters = append(q.Parameters, bigquery.QueryParameter{
			Name:  name,
			Value: value,
		})
	}
}

func NewClient(ctx context.Context, projectID string, opts ...option.ClientOption) (*Client, error) {
	client, err := bigquery.NewClient(ctx, projectID, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create bigquery client: %w", err)
	}
	return &Client{client: client}, nil
}

func String(s string) interface{} {
	return s
}

func (c *Client) Query(ctx context.Context, query string, opts ...QueryOption) ([]map[string]bigquery.Value, error) {
	q := c.client.Query(query)
	for _, opt := range opts {
		opt(q)
	}
	it, err := q.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	var rows []map[string]bigquery.Value
	for {
		var row map[string]bigquery.Value
		if err := it.Next(&row); err == iterator.Done {
			break
		} else if err != nil {
			return nil, fmt.Errorf("failed to iterate rows: %w", err)
		}
		rows = append(rows, row)
	}

	return rows, nil
}

func (c *Client) Close() error {
	return c.client.Close()
}
