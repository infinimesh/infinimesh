/*
Copyright © 2021-2022 Infinite Devices GmbH, Nikita Ivanovski info@slnt-opp.xyz

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
package schema

import (
	"testing"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	log *zap.Logger
	arangodbHost string
	arangodbCred string
)

func init() {
	viper.AutomaticEnv()
	log = zap.NewExample()

	viper.SetDefault("DB_HOST", "localhost:8529")
	viper.SetDefault("DB_CRED", "root:openSesame")

	arangodbHost = viper.GetString("DB_HOST")
	arangodbCred = viper.GetString("DB_CRED")
}

func TestInitDB(t *testing.T) {
	t.Log("Initializing DB")
	InitDB(log, arangodbHost, arangodbCred, "infinimesh")
}