// Copyright © 2022 Meroxa, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package chaos

import (
	"context"

	"github.com/conduitio/conduit-commons/opencdc"
	sdk "github.com/conduitio/conduit-connector-sdk"
)

type DestinationConfig struct {
	sdk.DefaultDestinationMiddleware

	// ConfigureMode controls what the Configure method should do.
	ConfigureMode string `validate:"inclusion=success|error|block|context-done|panic" default:"success"`
	// OpenMode controls what the Open method should do.
	OpenMode string `validate:"inclusion=success|error|block|context-done|panic" default:"success"`
	// WriteMode controls what the Write method should do.
	WriteMode string `validate:"inclusion=success|error|block|context-done|panic" default:"success"`
	// TeardownMode controls what the Teardown method should do.
	TeardownMode string `validate:"inclusion=success|error|block|context-done|panic" default:"success"`
}

func (c *DestinationConfig) Validate(ctx context.Context) error {
	return Chaos{}.Do(ctx, c.ConfigureMode)
}

type Destination struct {
	sdk.UnimplementedDestination

	config DestinationConfig

	isOpen bool
	chaos  Chaos
}

func NewDestination() sdk.Destination {
	return &Destination{}
}

func (d *Destination) Config() sdk.DestinationConfig {
	return &d.config
}

func (d *Destination) Open(ctx context.Context) error {
	d.isOpen = true
	return d.chaos.Do(ctx, d.config.OpenMode)
}

func (d *Destination) Write(ctx context.Context, records []opencdc.Record) (int, error) {
	err := d.chaos.Do(ctx, d.config.WriteMode)
	if err != nil {
		return 0, err
	}
	return len(records), nil
}

func (d *Destination) Teardown(ctx context.Context) error {
	if d.isOpen {
		// only do if connector is open, teardown also gets called when validating config
		return d.chaos.Do(ctx, d.config.TeardownMode)
	}
	return nil
}
