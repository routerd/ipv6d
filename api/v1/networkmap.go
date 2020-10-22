/*
Copyright 2020 The routerd Authors.

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

package v1

import (
	"github.com/jinzhu/copier"

	"routerd.net/ipv6d/internal/machinery/runtime"
)

type NetworkMapList struct {
	TypeMeta `json:",inline"`
	Items    []NetworkMap `json:"items"`
}

func (l *NetworkMapList) DeepCopy() *NetworkMapList {
	new := &NetworkMapList{}
	if err := copier.Copy(new, l); err != nil {
		panic(err)
	}
	return new
}

func (l *NetworkMapList) DeepCopyObject() runtime.Object {
	return l.DeepCopy()
}

// NetworkMap uses Network Prefix Translation to map one IPv6 Network to another.
type NetworkMap struct {
	TypeMeta   `json:",inline"`
	ObjectMeta `json:"metadata"`
	Spec       NetworkMapSpec   `json:"spec"`
	Status     NetworkMapStatus `json:"status"`
}

type NetworkMapSpec struct {
	// WAN Interface
	// will be used as -i <IFACE> / -o <IFACE>
	// in the iptables rule spec
	WANInterface string `json:"wanInterface"`

	// Network Maps
	// Configures how a private network is mapped to a public network.
	NetMap []NetMap `json:"netmap,omitempty"`
}

type NetMap struct {
	// Private Network in the fd00::/8 range
	Private NetworkMapNetworkPointer `json:"private"`
	// Public Network
	Public NetworkMapNetworkPointer `json:"public"`
}

// NetworkMapNetworkPointer specifies a network,
// or tells the controller how to look it up.
type NetworkMapNetworkPointer struct {
	// Static network configuration
	Static string `json:"static,omitempty"`
	// Interface to lookup, will take the first non-link-local IPv6 Net
	Interface string `json:"interface,omitempty"`
}

type NetworkMapStatus struct {
	ObservedGeneration int64          `json:"observedGeneration"`
	NetMap             []NetMapStatus `json:"netmap,omitempty"`
}

type NetMapStatus struct {
	Private string `json:"private"`
	Public  string `json:"public"`
}

func (nm *NetworkMap) DeepCopy() *NetworkMap {
	new := &NetworkMap{}
	if err := copier.Copy(new, nm); err != nil {
		panic(err)
	}
	return new
}

func (nm *NetworkMap) DeepCopyObject() runtime.Object {
	return nm.DeepCopy()
}

func init() {
	SchemeBuilder.Register(func(scheme *runtime.Scheme) error {
		scheme.AddKnownTypes(Version, &NetworkMap{}, &NetworkMapList{})
		return nil
	})
}
