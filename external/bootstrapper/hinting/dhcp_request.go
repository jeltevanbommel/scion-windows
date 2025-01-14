// Copyright 2020 Anapaya Systems
// Copyright 2021 ETH Zurich
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build linux
// +build linux

// Package github.com/insomniacslk/dhcp/client4 has u-root as dependency, which does not support windows,

package hinting

import (
	"fmt"

	"golang.org/x/sys/unix"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/client4"
)

func (g *DHCPHintGenerator) sendReceive(p *dhcpv4.DHCPv4) (*dhcpv4.DHCPv4, error) {
	p.SetBroadcast()
	client := client4.NewClient()
	sender, err := client4.MakeBroadcastSocket(g.iface.Name)
	if err != nil {
		return nil, fmt.Errorf("DHCP hinter failed to open broadcast sender socket: %w", err)
	}
	defer unix.Close(sender)
	receiver, err := client4.MakeListeningSocket(g.iface.Name)
	if err != nil {
		return nil, fmt.Errorf("DHCP hinter failed to open receiver socket: %w", err)
	}
	defer unix.Close(receiver)
	ack, err := client.SendReceive(sender, receiver, p, dhcpv4.MessageTypeAck)
	if err != nil {
		return nil, fmt.Errorf("DHCP hinter failed to send inform request: %w", err)
	}
	return ack, nil
}
