/**
 * Copyright (c) 2015 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package main

import (
	"log"
	"os"
	"strings"
)

var args = Config{}

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)

	// load from env variables

	args.ID = GetEnvVarAsString("GATEWAY_ID", "g1")
	args.Index = GetEnvVarAsInt("GATEWAY_INDEX", 0)
	args.Trace = GetEnvVarAsBool("GATEWAY_TRACE", false)
	
	args.Server.Root = GetEnvVarAsString("GATEWAY_SERVER_ROOT", "/ws")
	args.Server.Host = GetEnvVarAsString("GATEWAY_SERVER_HOST", "0.0.0.0")
	args.Server.Port = GetEnvVarAsInt("GATEWAY_SERVER_PORT", 8080)
	args.Server.Token = GetEnvVarAsString("GATEWAY_SERVER_TOKEN", "")
	args.Server.AuthMethod = GetEnvVarAsString("GATEWAY_SERVER_AUTHMETHOD", "none")
	args.Server.DeviceKeysURI = GetEnvVarAsString("GATEWAY_SERVER_DEVICEKEYSURI", "")
	args.Server.TolerableJWTAge = GetEnvVarAsInt("GATEWAY_SERVER_TOLERABLEJWTAGE", 5)

        var kafkaNodes string = GetEnvVarAsString("GATEWAY_PUB_URI", "docker:9091,docker:9092")

        if len(kafkaNodes) > 0 {
                args.Pub.URI = strings.Split(kafkaNodes, ",")
        }

	args.Pub.Topic = GetEnvVarAsString("GATEWAY_PUB_TOPIC", "messages")
	args.Pub.Ack = GetEnvVarAsBool("GATEWAY_PUB_ACK", false)
	args.Pub.Compress = GetEnvVarAsBool("GATEWAY_PUB_COMPRESS", true)
	args.Pub.FlushFreq = GetEnvVarAsInt("GATEWAY_PUB_FLUSHFREQ", 1)

	args.Server.AuthMethod = strings.ToLower(args.Server.AuthMethod)

	switch args.Server.AuthMethod {
	case "none":
	case "simple":
		if len(args.Server.Token) == 0 {
			log.Panicf("Simple auth requires a token")
		}
	case "jwt":
		if len(args.Server.DeviceKeysURI) == 0 {
			log.Panicf("JWT auth requires an API URI for public key retrieval")
		}
	default:
		log.Panicf("Invalid gateway authentication method: %v", args.Server.AuthMethod)
	}

	Trace("config", args)
}

// ServerConfig represents the Web server configuration holder
type ServerConfig struct {
	Root            string `json:"root,omitempty"`
	Host            string `json:"host,omitempty"`
	Port            int    `json:"port,omitempty"`
	Token           string `json:"token,omitempty"`
	AuthMethod      string `json:"auth_method"`
	DeviceKeysURI   string `json:"device_keys_uri,omitempty"`
	TolerableJWTAge int    `json:"tolerable_jwt_age,omitempty"`
}

// PubConfig represents the publisher configuration holder
type PubConfig struct {
	URI       []string `json:"uri,omitempty"`
	Topic     string   `json:"topic,omitempty"`
	Ack       bool     `json:"args,acks"`
	Compress  bool     `json:"args,compress"`
	FlushFreq int      `json:"args,flushevery"`
}

// Config represents the root object configuraiton holder
type Config struct {
	ID     string       `json:"id,omitempty"`
	Index  int          `json:"index,omitempty"`
	Trace  bool         `json:"trace,omitempty"`
	Server ServerConfig `json:"server,omitempty"`
	Pub    PubConfig    `json:"publisher,omitempty"`
}
