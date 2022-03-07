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

package main

import (
	"amqp_go/pkg/amqp"
	"flag"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

var listenAddr = flag.String("listen_host", "0.0.0.0", "amqp listen host")
var listenPort = flag.Int("listen_port", 5672, "amqp listen port")
var multiCore = flag.Bool("multi_core", false, "multi core")
var maxConn = flag.Int("max_conn", 500, "need sasl")

func main() {
	flag.Parse()
	serverConfig := &amqp.ServerConfig{}
	serverConfig.ListenHost = *listenAddr
	serverConfig.ListenPort = *listenPort
	serverConfig.MultiCore = *multiCore
	serverConfig.MaxConn = int32(*maxConn)
	e := &ExampleAmqpImpl{}
	_, err := amqp.Run(serverConfig, e)
	if err != nil {
		logrus.Error(err)
	} else {
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)
		for {
			<-interrupt
		}
	}
}
