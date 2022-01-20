//--------------------------------------------------------------------------
// Copyright 2018-2022 infinimesh
// www.infinimesh.io
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
//--------------------------------------------------------------------------

package main

import (
	"errors"
	"fmt"
	"net"

	"github.com/slntopp/infinimesh/pkg/registry/registrypb"
	"github.com/slntopp/mqtt-go/packet"
)

func HandleTCPConnections(tcp net.Listener) {
	for {
		conn, _ := tcp.Accept() // nolint: gosec

		p, err := packet.ReadPacket(conn, 0)
		if err != nil {
			LogErrorAndClose(conn, fmt.Errorf("Error while reading connect packet: %v", err))
			continue
		}
		if debug {
			fmt.Println("ControlPacket", p)
		}

		connectPacket, ok := p.(*packet.ConnectControlPacket)
		if !ok {
			LogErrorAndClose(conn, errors.New("Got wrong packet as first packet..need connect!"))
			continue
		}
		if debug {
			fmt.Println("ConnectPacket", p)
		}

		var fingerprint []byte
		fingerprint, err = verifyBasicAuth(connectPacket)
		if err != nil {
			LogErrorAndClose(conn, fmt.Errorf("Error verifying Basic Auth: %v", err))
			continue
		}

		if debug {
			fmt.Println("Fingerprint", string(fingerprint))
		}

		possibleIDs, err := GetByFingerprintAndVerify(fingerprint, func(device *registrypb.Device) (bool) {
			if device.Name != connectPacket.ConnectPayload.Username {
				fmt.Printf("Failed to verify client as the device name is doesn't match Basic Auth Username. Device ID:%v\n", device.Id)
				return false
			} else if !device.BasicEnabled.Value {
				fmt.Printf("Failed to verify client as the Basic Auth is not enabled for device. Device ID:%v\n", device.Id)
				return false
			} else if !device.Enabled.Value {
				fmt.Printf("Failed to verify client as the device is not enabled. Device ID:%v\n", device.Id)
				return false
			} else {
				fmt.Println(device.Tags)
				return true
			}
		})
		if err != nil {
			LogErrorAndClose(conn, err)
			continue
		}

		fmt.Printf("Client connected, IDs: %v\n", possibleIDs)

		go HandleConn(conn, connectPacket, possibleIDs)
	}
}