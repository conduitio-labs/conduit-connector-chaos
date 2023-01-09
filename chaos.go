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
	"errors"
	"fmt"
	"runtime"

	sdk "github.com/conduitio/conduit-connector-sdk"
)

const (
	ModeSuccess     = "success"
	ModeError       = "error"
	ModeBlock       = "block"
	ModeContextDone = "context-done"
	ModePanic       = "panic"
)

type Chaos struct{}

func (c Chaos) Do(ctx context.Context, mode string) error {
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
	case ModeContextDone:
		<-ctx.Done()
		return ctx.Err()
	default:
		panic(fmt.Errorf("invalid mode: %v", mode))
	}
}
