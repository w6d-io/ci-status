/*
Copyright 2020 WILDCARD

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
Created on 24/01/2021
*/

package pipelinerun

import (
	"context"
	"time"

	"github.com/avast/retry-go"
	"github.com/w6d-io/ci-status/internal/tekton"
	"github.com/w6d-io/x/logx"
)

const KIND = "pipelinerun"

// Scan stars the scan of pipeline run tekton resource
func Scan(ctx context.Context, payload *tekton.PipelineRunPayload) error {
	log := logx.WithName(ctx, "Scan").WithValues("name", payload.NamespacedName())
	log.V(1).Info("start")
	defer log.V(1).Info("stop")
	if err := retry.Do(func() error {
		t := &tekton.Tekton{
			PipelineRun: payload,
		}
		if err := t.PipelineRunSupervise(ctx); err != nil {
			return err
		}
		return nil
	},
		retry.Delay(3*time.Second),
		retry.Attempts(5),
	); err != nil {
		return err
	}
	return nil
}
