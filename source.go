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
	"time"

	"github.com/conduitio/conduit-commons/opencdc"
	sdk "github.com/conduitio/conduit-connector-sdk"
)

type SourceConfig struct {
	sdk.DefaultSourceMiddleware

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

func (c *SourceConfig) Validate(ctx context.Context) error {
	return Chaos{}.Do(ctx, c.ConfigureMode)
}

type Source struct {
	sdk.UnimplementedSource

	config SourceConfig

	isOpen bool
	chaos  Chaos
}

func NewSource() sdk.Source {
	return &Source{}
}

func (s *Source) Config() sdk.SourceConfig {
	return &s.config
}

func (s *Source) Open(ctx context.Context, _ opencdc.Position) error {
	s.isOpen = true
	return s.chaos.Do(ctx, s.config.OpenMode)
}

func (s *Source) Read(ctx context.Context) (opencdc.Record, error) {
	err := s.chaos.Do(ctx, s.config.ReadMode)
	if err != nil {
		return opencdc.Record{}, err
	}
	if ctx.Err() != nil {
		// TODO add mode that doesn't care about context closing
		return opencdc.Record{}, ctx.Err()
	}
	time.Sleep(time.Second)
	return sdk.Util.Source.NewRecordCreate(
		[]byte("chaos-position"),
		opencdc.Metadata{"chaos.readMode": s.config.ReadMode},
		opencdc.RawData("chaos-key"),
		opencdc.RawData("chaos-payload"),
	), nil
}

func (s *Source) Ack(ctx context.Context, _ opencdc.Position) error {
	return s.chaos.Do(ctx, s.config.AckMode)
}

func (s *Source) Teardown(ctx context.Context) error {
	if s.isOpen {
		// only do if connector is open, teardown also gets called when validating config
		return s.chaos.Do(ctx, s.config.TeardownMode)
	}
	return nil
}
