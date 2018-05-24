/*
Copyright 2018 Slack Inc.

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
	"vitess.io/vitess/go/slack"
	"vitess.io/vitess/go/streamlog"
	"vitess.io/vitess/go/vt/servenv"
	"vitess.io/vitess/go/vt/vtgate"
)

func init() {
	servenv.OnRun(func() {
		if slack.EnableMurronLogging() {
			initMurronLogger()
		}
	})
}

func initMurronLogger() {
	logger := slack.InitMurronLogger()
	logChan := vtgate.QueryLogger.Subscribe("Murron")
	formatParams := map[string][]string{"full": {}}
	formatter := streamlog.GetFormatter(vtgate.QueryLogger)

	go func() {
		for {
			record := <-logChan
			message := formatter(formatParams, record)
			logger.SendMessage(message)
		}
	}()
}
