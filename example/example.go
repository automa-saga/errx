// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"

	"github.com/joomcode/errorx"

	"github.com/automa-saga/errx"
)

var (
	errNS         = errorx.NewNamespace("myapp")
	errConfigLoad = errNS.NewType("config_load")
	errDBConnect  = errNS.NewType("db_connect", errorx.Temporary())
)

const (
	ReasonConfigMissing errx.Reason = "ConfigMissing"
	ReasonDBTimeout     errx.Reason = "DBTimeout"
)

func main() {
	// Create a structured error with a reason code and remediation hints.
	err := errx.Decorate(
		errConfigLoad.New("config file not found: app.yaml"),
		ReasonConfigMissing,
		"Create app.yaml in the working directory",
		"Copy from app.yaml.example and fill in required values",
	)

	// Extract reason for programmatic branching.
	if reason, ok := errx.ReasonOf(err); ok {
		fmt.Printf("reason = %s\n", reason)
	}

	// Render the full error with resolution block for a human.
	fmt.Fprintln(os.Stderr, errx.Format(err))
}
