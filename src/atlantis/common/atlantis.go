/* Copyright 2014 Ooyala, Inc. All rights reserved.
 *
 * This file is licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
 * except in compliance with the License. You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is
 * distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

package common

// ----------------------------------------------------------------------------------------------------------
// Atlantis Common Types
// ----------------------------------------------------------------------------------------------------------

const (
	StatusOk          = "OK"
	StatusMaintenance = "MAINTENANCE"
	StatusError       = "ERROR"
	StatusDegraded    = "DEGRADED"
	StatusUnknown     = "UNKNOWN"
	StatusDone        = "DONE"
	StatusInit        = "INIT"
	StatusFull        = "FULL" // Supervisor Health Check status when no more containers are available
	ManifestFile      = "manifest.toml"
	DefaultLDAPPort   = uint16(636)
	DefaultRegion     = "dev"
	DefaultZone       = "dev"
	JIRAPrefix        = "CR"
)

// ------------ Version -----------
// used to check manager or supervisor version
type VersionArg struct {
}

type VersionReply struct {
	RPCVersion string
	APIVersion string
}

// ------------ Async -----------
// used to for async requests
type AsyncReply struct {
	ID string
}

// ----------------------------------------------------------------------------------------------------------
// Utility Functions
// ----------------------------------------------------------------------------------------------------------

func DiffSlices(s1, s2 []string) (onlyInS1, onlyInS2 []string) {
	onlyInS1 = []string{}
	onlyInS2 = []string{}
	if s1 == nil && s2 == nil {
		return
	} else if s1 == nil {
		return onlyInS1, s2
	} else if s2 == nil {
		return s1, onlyInS2
	}
	counts := map[string]int{}
	for _, s1str := range s1 {
		counts[s1str]++
	}
	for _, s2str := range s2 {
		if count, present := counts[s2str]; !present || count == 0 {
			onlyInS2 = append(onlyInS2, s2str)
		} else {
			counts[s2str]--
		}
	}
	for s1str, count := range counts {
		for i := count; i > 0; i-- {
			onlyInS1 = append(onlyInS1, s1str)
		}
	}
	return
}
