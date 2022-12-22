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

//go:generate paramgen -output=paramgen_dest.go DestinationConfig

package chaos

import (
	"context"
	"errors"
	"fmt"
	"runtime"

	sdk "github.com/conduitio/conduit-connector-sdk"
)

type Destination struct {
	sdk.UnimplementedDestination

	Config DestinationConfig
	isOpen bool
}

func NewDestination() sdk.Destination {
	return &Destination{}
}

type DestinationConfig struct {
	// ConfigureMode controls what the Configure method should do.
	ConfigureMode string `validate:"inclusion=success|error|block|panic" default:"success"`
	// OpenMode controls what the Open method should do.
	OpenMode string `validate:"inclusion=success|error|block|panic" default:"success"`
	// WriteMode controls what the Write method should do.
	WriteMode string `validate:"inclusion=success|error|block|panic" default:"success"`
	// TeardownMode controls what the Teardown method should do.
	TeardownMode string `validate:"inclusion=success|error|block|panic" default:"success"`
}

const (
	ModeSuccess = "success"
	ModeError   = "error"
	ModeBlock   = "block"
	ModePanic   = "panic"
)

func (d *Destination) Parameters() map[string]sdk.Parameter {
	return d.Config.Parameters()
}

func (d *Destination) Configure(ctx context.Context, cfg map[string]string) error {
	err := sdk.Util.ParseConfig(cfg, &d.Config)
	if err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}
	return d.do(ctx, d.Config.ConfigureMode)
}

func (d *Destination) Open(ctx context.Context) error {
	d.isOpen = true
	return d.do(ctx, d.Config.OpenMode)
}

func (d *Destination) Write(ctx context.Context, records []sdk.Record) (int, error) {
	return 0, d.do(ctx, d.Config.WriteMode)
}

func (d *Destination) Teardown(ctx context.Context) error {
	if d.isOpen {
		// only do if connector is open, teardown also gets called when validating config
		return d.do(ctx, d.Config.TeardownMode)
	}
	return nil
}

func (d *Destination) do(ctx context.Context, mode string) error {
	var callingFunc string
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok {
		callingFunc = details.Name()
	}
	sdk.Logger(ctx).Info().Str("func", callingFunc).Str("mode", mode).Send()

	switch mode {
	case ModeSuccess:
		return nil
	case ModeError:
		return errors.New("chaos")
	case ModeBlock:
		<-make(chan struct{}) // block forever
		return nil
	case ModePanic:
		panic("chaos")
	default:
		panic(fmt.Errorf("invalid mode: %v", mode))
	}
}
