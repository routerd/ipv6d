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
package ipv6

import "net"

var (
	LinkLocalNet *net.IPNet
	PrivateNet   *net.IPNet
)

// IsLinkLocal checks whether the given IP is in the fe80::/10 network.
func IsLinkLocal(ip net.IP) bool {
	return LinkLocalNet.Contains(ip)
}

// IsPrivate checks wether the given IP is in the fd00::/8 private network range.
func IsPrivate(ip net.IP) bool {
	return PrivateNet.Contains(ip)
}

// IsIPv6 checks wether the given IP is IPv6
func IsIPv6(ip net.IP) bool {
	return ip.To4() == nil
}

func init() {
	var err error
	_, LinkLocalNet, err = net.ParseCIDR("fe80::/10")
	if err != nil {
		panic(err)
	}

	_, PrivateNet, err = net.ParseCIDR("fd00::/8")
	if err != nil {
		panic(err)
	}
}
