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

package network

import (
	"amqp_go/pkg/service"
	"fmt"
	"github.com/panjf2000/gnet"
	"github.com/sirupsen/logrus"
	"sync"
	"sync/atomic"
)

// connCount amqp connection count
var connCount int32

func Run(config *Config, impl service.AmqpServer) (*Server, error) {
	server := &Server{
		EventServer: nil,
		amqpImpl:    impl,
	}
	go func() {
		err := gnet.Serve(server, fmt.Sprintf("tcp://%s:%d", config.ListenHost, config.ListenPort), gnet.WithMulticore(config.MultiCore))
		logrus.Error("amqp broker started error ", err)
	}()
	return server, nil
}

type Server struct {
	*gnet.EventServer
	config   Config
	ConnMap  sync.Map
	amqpImpl service.AmqpServer
}

func (s *Server) OnInitComplete(server gnet.Server) (action gnet.Action) {
	logrus.Info("Amqp Server started")
	return
}

func (s *Server) React(frame []byte, c gnet.Conn) ([]byte, gnet.Action) {
	return nil, gnet.Close
}

func (s *Server) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	if atomic.LoadInt32(&connCount) > s.config.MaxConn {
		logrus.Error("connection reach max, refused to connect ", c.RemoteAddr())
		return nil, gnet.Close
	}
	connCount := atomic.AddInt32(&connCount, 1)
	s.ConnMap.Store(c.RemoteAddr(), c)
	logrus.Info("new connection connected ", connCount, " from ", c.RemoteAddr())
	return
}

func (s *Server) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	logrus.Info("connection closed from ", c.RemoteAddr())
	s.ConnMap.Delete(c.RemoteAddr())
	atomic.AddInt32(&connCount, -1)
	return
}
