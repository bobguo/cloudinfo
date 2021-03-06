// Copyright © 2019 Banzai Cloud
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

package cloudinfo

import (
	"github.com/banzaicloud/cloudinfo/pkg/cloudinfo"
)

// InMemoryProviderStore keeps instance types in the memory.
// Use it in tests or for development/demo purposes.
type InMemoryProviderStore struct {
	providers []cloudinfo.Provider
	services  map[string][]cloudinfo.Service
}

// NewInMemoryProviderStore returns a new InMemoryProviderStore.
func NewInMemoryProviderStore() *InMemoryProviderStore {
	return &InMemoryProviderStore{
		providers: []cloudinfo.Provider{},
		services:  make(map[string][]cloudinfo.Service),
	}
}

func (s *InMemoryProviderStore) GetProviders() ([]cloudinfo.Provider, error) {
	return s.providers, nil
}

func (s *InMemoryProviderStore) GetServices(provider string) ([]cloudinfo.Service, error) {
	return s.services[provider], nil
}
