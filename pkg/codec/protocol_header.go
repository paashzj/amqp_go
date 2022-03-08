// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package codec

import (
	"errors"
	"github.com/sirupsen/logrus"
	"runtime/debug"
)

type ProtocolHeader struct {
	ProtocolId              int8
	ProtocolMajorVersion    int8
	ProtocolMinorVersion    int8
	ProtocolRevisionVersion int8
}

func IsProtocolHeader(bytes []byte) bool {
	if len(bytes) != 8 {
		return false
	}
	return bytes[0] == 0x41 && bytes[1] == 0x4d && bytes[2] == 0x51 && bytes[3] == 0x50
}

func DecodeProtocolHeader(bytes []byte) (protocolHeader *ProtocolHeader, err error) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Info("Recovered in f", r, string(debug.Stack()))
			protocolHeader = nil
			err = errors.New("codec failed")
		}
	}()
	p := &ProtocolHeader{}
	p.ProtocolId = int8(bytes[4])
	p.ProtocolMajorVersion = int8(bytes[5])
	p.ProtocolMinorVersion = int8(bytes[6])
	p.ProtocolRevisionVersion = int8(bytes[7])
	return p, err
}

func (p *ProtocolHeader) BytesLength() int {
	return 8
}

func (p *ProtocolHeader) Bytes() []byte {
	bytes := make([]byte, p.BytesLength())
	bytes[0] = 0x41
	bytes[1] = 0x4d
	bytes[2] = 0x51
	bytes[3] = 0x50
	bytes[4] = byte(p.ProtocolId)
	bytes[5] = byte(p.ProtocolMajorVersion)
	bytes[6] = byte(p.ProtocolMinorVersion)
	bytes[7] = byte(p.ProtocolRevisionVersion)
	return bytes
}
