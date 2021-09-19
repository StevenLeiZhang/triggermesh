/*
Copyright 2020 TriggerMesh Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package webhooksource

import (
	"knative.dev/eventing/pkg/adapter/v2"
)

// EnvAccessor for configuration parameters
func EnvAccessor() adapter.EnvConfigAccessor {
	return &envAccessor{}
}

type envAccessor struct {
	adapter.EnvConfig

	EventType         string `envconfig:"WEBHOOK_EVENT_TYPE" required:"true"`
	EventSource       string `envconfig:"WEBHOOK_EVENT_SOURCE" required:"true"`
	BasicAuthUsername string `envconfig:"WEBHOOK_BASICAUTH_USERNAME"`
	BasicAuthPassword string `envconfig:"WEBHOOK_BASICAUTH_PASSWORD"`
}
