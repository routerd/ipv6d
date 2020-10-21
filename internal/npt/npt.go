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

package npt

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/coreos/go-iptables/iptables"
	"github.com/go-logr/logr"

	v1 "routerd.net/ipv6d/api/v1"
	"routerd.net/ipv6d/internal/machinery/controller"
	"routerd.net/ipv6d/internal/machinery/errors"
	"routerd.net/ipv6d/internal/machinery/state"
	"routerd.net/ipv6d/internal/utils/ipv6"
)

// Reconciler ensures private networks are mapped to public networks via Network Prefix Translation.
type Reconciler struct {
	log            logr.Logger
	client         state.Client
	resyncDuration time.Duration
	ip6tables      *iptables.IPTables
	controller     *controller.Controller
}

func NewReconciler(
	log logr.Logger, client state.Client, resyncDuration time.Duration,
) (*Reconciler, error) {
	ip6tables, err := iptables.NewWithProtocol(iptables.ProtocolIPv6)
	if err != nil {
		return nil, fmt.Errorf("ip6tables: %w", err)
	}

	r := &Reconciler{
		log:            log,
		client:         client,
		resyncDuration: resyncDuration,
		ip6tables:      ip6tables,
	}
	r.controller = controller.NewController(log, r)
	return r, nil
}

func (r *Reconciler) Run(stopCh <-chan struct{}) {
	r.controller.Run(stopCh)
}

func (r *Reconciler) Reconcile(key string) (res controller.Result, err error) {
	ctx := context.Background()
	netmap := &v1.NetworkMap{}
	if err := r.client.Get(ctx, key, netmap); err != nil {
		if errors.IsNotFound(err) {
			return res, nil
		}
		return res, err
	}

	res.RequeueAfter = r.resyncDuration

	for _, rule := range r.rules(netmap) {
		exists, err := r.ip6tables.Exists(rule.Table, rule.Chain, rule.Spec...)
		if err != nil {
			return res, fmt.Errorf("checking rule exists: %w", err)
		}

		if exists {
			continue
		}
		if err := r.ip6tables.Append(rule.Table, rule.Chain, rule.Spec...); err != nil {
			return res, fmt.Errorf("append new rule: %w", err)
		}
	}
	return
}

type rule struct {
	Table, Chain string
	Spec         []string
}

func (r *Reconciler) rules(netmap *v1.NetworkMap) []rule {
	var rules []rule

	for i, nm := range netmap.Spec.NetMap {
		inbound, outbound, status, err := r.rule(netmap.Spec.WANInterface, nm)
		if err != nil {
			r.log.Error(err, "iptable rule for .Spec.NetMap [%d]", i)
			continue
		}

		rules = append(rules, *inbound, *outbound)
		netmap.Status.NetMap = append(netmap.Status.NetMap, *status)
	}
	netmap.Status.ObservedGeneration = netmap.Generation
	return rules
}

func (r *Reconciler) rule(wanInterface string, netmap v1.NetMap) (
	inbound, outbound *rule, nmStatus *v1.NetMapStatus, err error) {
	privateNetwork, err := r.lookupNetworkPointer(netmap.Private)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("lookup private ip: %w", err)
	}

	publicNetwork, err := r.lookupNetworkPointer(netmap.Public)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("lookup public ip: %w", err)
	}

	outbound = &rule{
		Table: "nat", Chain: "POSTROUTING",
		Spec: []string{
			"-o", wanInterface,
			"-s", privateNetwork.String(),
			"-j", "NETMAP",
			"--to", publicNetwork.String(),
		},
	}

	inbound = &rule{
		Table: "nat", Chain: "PREROUTING",
		Spec: []string{
			"-i", wanInterface,
			"-d", publicNetwork.String(),
			"-j", "NETMAP",
			"--to", privateNetwork.String(),
		},
	}

	nmStatus = &v1.NetMapStatus{}
	nmStatus.Private = privateNetwork.String()
	nmStatus.Public = publicNetwork.String()
	return
}

func (r *Reconciler) lookupNetworkPointer(pointer v1.NetworkMapNetworkPointer) (*net.IPNet, error) {
	if pointer.Static != "" {
		// Static Network, just parse and GO!
		_, ipNet, err := net.ParseCIDR(pointer.Static)
		return ipNet, err
	}

	// Lookup from interface name
	return r.lookupNetworkFromInterface(pointer.Interface)
}

func (r *Reconciler) lookupNetworkFromInterface(ifaceName string) (*net.IPNet, error) {
	// Interface Addresses
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return nil, fmt.Errorf(
			"get interface by name %s: %w", ifaceName, err)
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}

	// Public Network
	var publicNetwork *net.IPNet
	for _, addr := range addrs {
		ip, ipNet, err := net.ParseCIDR(addr.String())
		if err != nil {
			return nil, fmt.Errorf("parse address %s: %w", addr, err)
		}

		if ipv6.IsIPv6(ip) {
			// ignore IPv4 addresses
			continue
		}
		if ipv6.IsLinkLocal(ip) {
			continue
		}

		publicNetwork = ipNet
	}
	if publicNetwork == nil {
		return nil, fmt.Errorf(
			"interface %s has no public IPv6 address", ifaceName)
	}
	return nil, nil
}
