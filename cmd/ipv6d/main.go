/*
Copyright 2020 The routerd Authors.

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
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/bombsimon/logrusr"
	"github.com/go-logr/logr"
	"github.com/sirupsen/logrus"

	v1 "routerd.net/ipv6d/api/v1"
	"routerd.net/ipv6d/internal/controllers"
	"routerd.net/ipv6d/internal/machinery/runtime"
	"routerd.net/ipv6d/internal/machinery/state"
)

func main() {
	var configFolder string
	flag.StringVar(&configFolder, "config-folder", "", "config file folder.")
	flag.Parse()

	log := logrusr.NewLogger(logrus.New())

	if err := run(log, configFolder); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(log logr.Logger, configFolder string) error {
	scheme := runtime.NewScheme()
	if err := v1.AddToScheme(scheme); err != nil {
		return err
	}

	metaRepo, err := state.NewMetaRepository(log.WithName("repository"), scheme)
	if err != nil {
		return err
	}
	stopCh := make(chan struct{})
	go metaRepo.Run(stopCh)

	if err := metaRepo.LoadFromFileSystem(configFolder); err != nil {
		return err
	}

	r, err := controllers.NewNPTReconciler(log.WithName("NPTController"), metaRepo, 5*time.Second)
	if err != nil {
		return err
	}

	r.Run(stopCh)
	return nil
}
