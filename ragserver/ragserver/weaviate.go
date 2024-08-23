// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// Utilities for working with Weaviate.

import (
	"cmp"
	"context"
	"fmt"
	"os"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate/entities/models"
)

// initWeaviate initializes a weaviate client for our application.
func initWeaviate(ctx context.Context) (*weaviate.Client, error) {
	client, err := weaviate.NewClient(weaviate.Config{
		Host:   "localhost:" + cmp.Or(os.Getenv("WVPORT"), "9035"),
		Scheme: "http",
	})
	if err != nil {
		return nil, fmt.Errorf("initializing weaviate: %w", err)
	}

	// Create a new class (collection) in weaviate if it doesn't exist yet.
	cls := &models.Class{
		Class:      "Document",
		Vectorizer: "none",
	}
	exists, err := client.Schema().ClassExistenceChecker().WithClassName(cls.Class).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("weaviate error: %w", err)
	}
	if !exists {
		err = client.Schema().ClassCreator().WithClass(cls).Do(ctx)
		if err != nil {
			return nil, fmt.Errorf("weaviate error: %w", err)
		}
	}

	return client, nil
}

// combinedWeaviateError generates an error if err is non-nil or result has
// errors, and returns an error (or nil if there's no error). It's useful for
// the results of the Weaviate GraphQL API's "Do" calls.
func combinedWeaviateError(result *models.GraphQLResponse, err error) error {
	if err != nil {
		return err
	}
	if len(result.Errors) != 0 {
		var ss []string
		for _, e := range result.Errors {
			ss = append(ss, e.Message)
		}
		return fmt.Errorf("weaviate error: %v", ss)
	}
	return nil
}
