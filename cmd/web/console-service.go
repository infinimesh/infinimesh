/*
Copyright © 2021-2023 Infinite Devices GmbH

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
package main

import (
	"encoding/json"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var SERVICES_ENABLED = []byte("{}")
var SERVICES_ENABLED_MAP = map[string]string{}

func init_cs(log *zap.Logger) {
	var enabled = map[string]string{}

	viper.SetDefault("HTTP_FS", "")
	http_fs := viper.GetString("HTTP_FS")
	log.Debug("HTTP_FS", zap.String("value", http_fs))
	if http_fs != "" {
		log.Info("HTTP-FS is set", zap.String("url", http_fs))
		enabled["http_fs"] = http_fs
		SERVICES_ENABLED_MAP["http_fs"] = http_fs
	}

	viper.SetDefault("HANDSFREE", "")
	handsfree := viper.GetString("HANDSFREE")
	log.Debug("HANDSFREE", zap.String("value", handsfree))
	if handsfree != "" {
		log.Info("Handsfree service is enabled")
		enabled["handsfree"] = ""
		SERVICES_ENABLED_MAP["handsfree"] = handsfree
	}

	viper.SetDefault("CHATTING", "")
	chatting := viper.GetString("CHATTING")
	log.Debug("CHATTING", zap.String("value", chatting))
	if chatting != "" {
		log.Info("Chatting service is enabled")
		enabled["chatting"] = chatting
		SERVICES_ENABLED_MAP["chatting"] = chatting
	}

	viper.SetDefault("TIMESERIES", "")
	timeseries := viper.GetString("TIMESERIES")
	log.Debug("TIMESERIES", zap.String("value", timeseries))
	if timeseries != "" {
		log.Info("Timeseries service is enabled")
		enabled["timeseries"] = timeseries
		SERVICES_ENABLED_MAP["timeseries"] = timeseries
	}

	var err error
	SERVICES_ENABLED, err = json.Marshal(enabled)
	if err != nil {
		panic(err)
	}
}

func cs_handler() runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(SERVICES_ENABLED)
	}
}
