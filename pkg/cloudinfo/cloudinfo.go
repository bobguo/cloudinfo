// Copyright © 2018 Banzai Cloud
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
	"strings"

	"github.com/goph/emperror"
	"github.com/goph/logur"
	"github.com/pkg/errors"
)

// cachingCloudInfo is the module struct, holds configuration and cache
// It's the entry point for the product info retrieval and management subsystem
// It's also responsible for delegating to the cloud provider specific implementations
type cachingCloudInfo struct {
	log            logur.Logger
	cloudInfoers   map[string]CloudInfoer
	cloudInfoStore CloudInfoStore
}

func (v AttrValues) floatValues() []float64 {
	floatValues := make([]float64, len(v))
	for i, av := range v {
		floatValues[i] = av.Value
	}
	return floatValues
}

// SpotPriceInfo represents different prices per availability zones
type SpotPriceInfo map[string]float64

// Price describes the on demand price and spot prices per availability zones
type Price struct {
	OnDemandPrice float64       `json:"onDemandPrice"`
	SpotPrice     SpotPriceInfo `json:"spotPrice"`
}

// VmInfo representation of a virtual machine
type VmInfo struct {
	Category      string            `json:"category"`
	Type          string            `json:"type"`
	OnDemandPrice float64           `json:"onDemandPrice"`
	SpotPrice     []ZonePrice       `json:"spotPrice"`
	Cpus          float64           `json:"cpusPerVm"`
	Mem           float64           `json:"memPerVm"`
	Gpus          float64           `json:"gpusPerVm"`
	NtwPerf       string            `json:"ntwPerf"`
	NtwPerfCat    string            `json:"ntwPerfCategory"`
	Zones         []string          `json:"zones"`
	Attributes    map[string]string `json:"attributes"`
	// CurrentGen signals whether the instance type generation is the current one. Only applies for amazon
	CurrentGen bool `json:"currentGen"`
}

// IsBurst returns true if the EC2 instance vCPU is burst type
// the decision is made based on the instance type
func (vm VmInfo) IsBurst() bool {
	return strings.HasPrefix(strings.ToUpper(vm.Type), "T")
}

// NewCachingCloudInfo creates a new cachingCloudInfo instance
func NewCachingCloudInfo(infoers map[string]CloudInfoer, ciStore CloudInfoStore, logger logur.Logger) (*cachingCloudInfo, error) {
	if infoers == nil || ciStore == nil {
		return nil, errors.New("could not create product infoer")
	}

	pi := cachingCloudInfo{
		cloudInfoers:   infoers,
		cloudInfoStore: ciStore,
		log:            logur.WithFields(logger, map[string]interface{}{"component": "cachingCloudInfo"}),
	}
	return &pi, nil
}

// GetProviders returns the supported providers
func (cpi *cachingCloudInfo) GetProviders() ([]Provider, error) {
	var (
		providers []Provider
		provider  Provider
		err       error
	)

	// iterate over supported provider names only
	for pn := range cpi.cloudInfoers {
		if provider, err = cpi.GetProvider(pn); err != nil {
			return nil, err
		}

		providers = append(providers, provider)
	}

	return providers, nil
}

// GetProvider returns the supported provider
func (cpi *cachingCloudInfo) GetProvider(provider string) (Provider, error) {
	var (
		srvcs []Service
		err   error
	)

	if _, ok := cpi.cloudInfoers[provider]; !ok {
		return Provider{}, emperror.With(errors.New("unsupported provider"), "provider", provider)
	}

	if srvcs, err = cpi.GetServices(provider); err != nil {
		return Provider{}, emperror.With(errors.New("no supported services for provider"), "provider", provider)
	}

	// decorate the provider with service information
	p := NewProvider(provider)
	p.Services = srvcs

	return p, nil
}

// GetZones returns the availability zones in a region
func (cpi *cachingCloudInfo) GetZones(provider, service, region string) ([]string, error) {
	if cachedVal, ok := cpi.cloudInfoStore.GetZones(provider, service, region); ok {
		return cachedVal, nil
	}

	return nil, emperror.With(errors.New("zones not yet cached"), "provider", provider, "region", region)
}

// GetRegions gets the regions for the provided provider
func (cpi *cachingCloudInfo) GetRegions(provider, service string) (map[string]string, error) {
	if cachedVal, ok := cpi.cloudInfoStore.GetRegions(provider, service); ok {
		return cachedVal, nil
	}

	return nil, emperror.With(errors.New("regions not yet cached"), "provider", provider, "services", service)
}

func (cpi *cachingCloudInfo) GetServices(provider string) ([]Service, error) {
	if cachedVal, ok := cpi.cloudInfoStore.GetServices(provider); ok {
		return cachedVal, nil
	}

	return nil, emperror.With(errors.New("services not yet cached"), "provider", provider)
}

// GetProductDetails retrieves product details form the given provider and region
func (cpi *cachingCloudInfo) GetProductDetails(provider, service, region string) ([]ProductDetails, error) {
	var (
		vms interface{}
		ok  bool
	)

	if vms, ok = cpi.cloudInfoStore.GetVm(provider, service, region); !ok {
		return nil, emperror.With(errors.New("vms not yet cached"),
			"provider", provider, "service", service, "region", region)
	}

	var details []ProductDetails

	for _, vm := range vms.([]VmInfo) {
		pd := newProductDetails(vm)
		if cachedVal, ok := cpi.cloudInfoStore.GetPrice(provider, region, vm.Type); ok {
			for zone, price := range cachedVal.SpotPrice {
				pd.SpotPrice = append(pd.SpotPrice, *newZonePrice(zone, price))
			}
		} else {
			cpi.log.Debug("price info not yet cached", map[string]interface{}{"instanceType": vm.Type})
		}

		details = append(details, *pd)
	}

	return details, nil
}

// GetStatus retrieves status form the given provider
func (cpi *cachingCloudInfo) GetStatus(provider string) (string, error) {
	if cachedStatus, ok := cpi.cloudInfoStore.GetStatus(provider); ok {
		return cachedStatus, nil
	}
	return "", emperror.With(errors.New("status not yet cached"), "provider", provider)
}

// GetServiceImages retrieves available images for the given provider, service and region
func (cpi *cachingCloudInfo) GetServiceImages(provider, service, region string) ([]Image, error) {
	if cachedImages, ok := cpi.cloudInfoStore.GetImage(provider, service, region); ok {
		return cachedImages, nil
	}

	return nil, emperror.With(errors.New("images not yet cached"), "provider", provider,
		"service", service, "region", region)
}

// GetVersions retrieves available versions for the given provider, service and region
func (cpi *cachingCloudInfo) GetVersions(provider, service, region string) ([]LocationVersion, error) {
	if cachedVersions, ok := cpi.cloudInfoStore.GetVersion(provider, service, region); ok {
		return cachedVersions, nil
	}
	return nil, emperror.With(errors.New("versions not yet cached"),
		"provider", provider, "service", service, "region", region)
}

// GetContinents retrieves available continents
func (cpi *cachingCloudInfo) GetContinents() []string {
	return []string{ContinentAsia, ContinentAustralia, ContinentEurope, ContinentNorthAmerica, ContinentSouthAmerica}
}

// GetContinents gets the continents and regions for the provided provider
func (cpi *cachingCloudInfo) GetContinentsData(provider, service string) (map[string][]Region, error) {
	if cachedVal, ok := cpi.cloudInfoStore.GetRegions(provider, service); ok {
		var continents = make(map[string][]Region)
		for id, name := range cachedVal {
			continent := getContinent(id)
			continents[continent] = append(continents[continent], Region{
				Id:   id,
				Name: name,
			})
		}
		return continents, nil
	}

	return nil, emperror.With(errors.New("regions not yet cached"), "provider", provider, "services", service)
}

// getContinent categorizes regions by continents
func getContinent(region string) string {
	switch {
	case checkContinent(region, []string{"ap-southeast-2", "australia"}):
		return ContinentAustralia
	case checkContinent(region, []string{"cn-", "ap-", "me-", "asia", "japan", "india", "korea"}):
		return ContinentAsia
	case checkContinent(region, []string{"eu", "uk", "france"}):
		return ContinentEurope
	case checkContinent(region, []string{"us", "ca-central-1", "canada", "northamerica"}):
		return ContinentNorthAmerica
	case checkContinent(region, []string{"southamerica", "brazil", "sa-"}):
		return ContinentSouthAmerica
	case checkContinent(region, []string{"africa"}):
		return ContinentAfrica
	default:
		return "unknown"
	}
}

func checkContinent(region string, substrs []string) bool {
	for _, substr := range substrs {
		if strings.Contains(region, substr) {
			return true
		}
	}
	return false
}

// Contains is a helper function to check if a slice contains a string
func Contains(slice []string, s string) bool {
	for _, e := range slice {
		if e == s {
			return true
		}
	}
	return false
}

// Attributes create a map with the specified parameters
func Attributes(cpu, memory, ntwPerfCat, vmCategory string) map[string]string {
	var attributes = make(map[string]string)

	attributes[Cpu] = cpu
	attributes[Memory] = memory
	attributes["networkPerfCategory"] = ntwPerfCat
	attributes["instanceTypeCategory"] = vmCategory

	return attributes
}

// CreateString builds strings effectively
func CreateString(keys ...string) string {
	var b strings.Builder
	for _, key := range keys {
		b.WriteString(key)
	}
	return b.String()
}
