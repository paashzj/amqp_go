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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidProtocolHeader(t *testing.T) {
	bytes := testHex2Bytes(t, "414d515003010000")
	assert.True(t, IsProtocolHeader(bytes))
}

func TestDecodeProtocolHeader(t *testing.T) {
	bytes := testHex2Bytes(t, "414d515003010000")
	protocolHeader, err := DecodeProtocolHeader(bytes)
	assert.Nil(t, err)
	var expectedProtocolId int8 = 3
	assert.Equal(t, protocolHeader.ProtocolId, expectedProtocolId)
	var expectedProtocolMajorVersion int8 = 1
	assert.Equal(t, protocolHeader.ProtocolMajorVersion, expectedProtocolMajorVersion)
	var expectedProtocolMinorVersion int8 = 0
	assert.Equal(t, protocolHeader.ProtocolMinorVersion, expectedProtocolMinorVersion)
	var expectedProtocolRevisionVersion int8 = 0
	assert.Equal(t, protocolHeader.ProtocolRevisionVersion, expectedProtocolRevisionVersion)
}

func TestCodeProtocolHeader(t *testing.T) {
	p := &ProtocolHeader{}
	p.ProtocolId = 3
	p.ProtocolMajorVersion = 1
	p.ProtocolMinorVersion = 0
	p.ProtocolRevisionVersion = 0
	bytes := p.Bytes()
	expectBytes := testHex2Bytes(t, "414d515003010000")
	assert.Equal(t, expectBytes, bytes)
}
