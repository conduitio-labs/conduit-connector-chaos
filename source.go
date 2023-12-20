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

//go:generate paramgen -output=paramgen_src.go SourceConfig

package chaos

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/conduitio/conduit-connector-sdk"
)

type Source struct {
	sdk.UnimplementedSource

	Config SourceConfig

	isOpen bool
	chaos  Chaos
}

func NewSource() sdk.Source {
	return &Source{}
}

type SourceConfig struct {
	// ConfigureMode controls what the Configure method should do.
	ConfigureMode string `validate:"inclusion=success|error|context-done|block|panic" default:"success"`
	// OpenMode controls what the Open method should do.
	OpenMode string `validate:"inclusion=success|error|context-done|block|panic" default:"success"`
	// ReadMode controls what the Read method should do.
	ReadMode string `validate:"inclusion=success|error|context-done|block|panic" default:"success"`
	// AckMode controls what the Ack method should do.
	AckMode string `validate:"inclusion=success|error|context-done|block|panic" default:"success"`
	// TeardownMode controls what the Teardown method should do.
	TeardownMode string `validate:"inclusion=success|error|context-done|block|panic" default:"success"`
}

func (d *Source) Parameters() map[string]sdk.Parameter {
	return d.Config.Parameters()
}

func (d *Source) Configure(ctx context.Context, cfg map[string]string) error {
	err := sdk.Util.ParseConfig(cfg, &d.Config)
	if err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}
	return d.chaos.Do(ctx, d.Config.ConfigureMode)
}

func (d *Source) Open(ctx context.Context, _ sdk.Position) error {
	d.isOpen = true
	return d.chaos.Do(ctx, d.Config.OpenMode)
}

func (d *Source) Read(ctx context.Context) (sdk.Record, error) {
	err := d.chaos.Do(ctx, d.Config.ReadMode)
	if err != nil {
		return sdk.Record{}, err
	}
	if ctx.Err() != nil {
		// TODO add mode that doesn't care about context closing
		return sdk.Record{}, ctx.Err()
	}
	time.Sleep(time.Second)
	return sdk.Util.Source.NewRecordCreate(
		[]byte("chaos-position"),
		sdk.Metadata{"chaos.readMode": d.Config.ReadMode},
		sdk.RawData("chaos-key"),
		sdk.RawData("chaos-payload"),
	), nil
}

func (d *Source) Ack(ctx context.Context, _ sdk.Position) error {
	return d.chaos.Do(ctx, d.Config.AckMode)
}

func (d *Source) Teardown(ctx context.Context) error {
	if d.isOpen {
		// only do if connector is open, teardown also gets called when validating config
		return d.chaos.Do(ctx, d.Config.TeardownMode)
	}
	return nil
}
