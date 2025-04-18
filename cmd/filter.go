// Copyright 2024 Fabian `xx4h` Sylvester
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

package cmd

import (
	"slices"
	"strings"

	"github.com/xx4h/hctl/pkg/rest"
)

func filterStates(states []rest.HassState, ignoredStates []string) []rest.HassState {
	if ignoredStates == nil {
		return states
	}

	var filteredStates []rest.HassState
	for _, rel := range states {
		found := false
		for _, ignoredName := range ignoredStates {
			if rel.EntityID == ignoredName {
				found = true
				break
			}
		}
		if !found {
			filteredStates = append(filteredStates, rel)
		}
	}

	return filteredStates
}

func filterAttributes(states []rest.HassState, attributes []string) []rest.HassState {
	if attributes == nil {
		return states
	}

	var filteredStates []rest.HassState

	for rel := range states {
		for _, attr := range attributes {
			if _, ok := states[rel].Attributes[attr]; ok {
				filteredStates = append(filteredStates, states[rel])
				// no need to further iterate attributes
				break
			}
		}
	}

	return filteredStates
}

// Filter states with given service capability and state
func filterCapable(states []rest.HassState, services []rest.HassService, serviceCaps []string, state string) []rest.HassState {
	if serviceCaps == nil {
		return states
	}

	// get all service domains that have "turn_on" as domain service
	// split state.EntryId domain=[0] entity=[1]
	// create list of states that are in a domain having "turn_on" as domain service
	// return only states from the list where state.State == off
	var capableServices []rest.HassService
	for _, rel := range services {
		for name := range rel.Services {
			if slices.Contains(serviceCaps, name) {
				capableServices = append(capableServices, rel)
			}
		}
	}

	// in case we have multiple capableServices, trackList will ensure that we don't add states twice due tue their multiple capabilities
	var trackList []string
	var filteredStates []rest.HassState
	for rel := range states {
		s := strings.Split(states[rel].EntityID, ".")
		stateDomain := s[0]
		for svc := range capableServices {
			if stateDomain == capableServices[svc].Domain {
				if (state == "" || states[rel].State == state) && !slices.Contains(trackList, states[rel].EntityID) {
					filteredStates = append(filteredStates, states[rel])
					trackList = append(trackList, states[rel].EntityID)
				}
			}
		}
	}

	return filteredStates
}
