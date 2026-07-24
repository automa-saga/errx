# errx

`errx` is a structured error handling library for Go built on top of [errorx](https://github.com/joomcode/errorx). It provides two cross-cutting error properties — a machine-readable **reason code** and a list of human **remediation hints** — plus helpers to attach and extract them.

## Features

- Machine-readable `Reason` codes for programmatic error branching
- Human-readable remediation hints attached to errors
- Cause-chain traversal through both errorx and `errors.Unwrap` layers
- `Format()` renderer for CLI/diagnostic output with resolution blocks

## Installation

```sh
go get github.com/automa-saga/errx
```

## Usage

```go
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
)

const ReasonConfigMissing errx.Reason = "ConfigMissing"

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
```

Output:
```
reason = ConfigMissing
myapp.config_load: config file not found: app.yaml {reason: ConfigMissing}

Resolution:
  - Create app.yaml in the working directory
  - Copy from app.yaml.example and fill in required values
```

## API

| Function | Description |
|---|---|
| `WithReason(err, reason)` | Attach a `Reason` code to an error |
| `WithHints(err, hints...)` | Attach remediation hints to an error |
| `Decorate(err, reason, hints...)` | Attach both reason and hints in one call |
| `ReasonOf(err)` | Extract the `Reason` from anywhere in the cause chain |
| `Hints(err)` | Extract hints from anywhere in the cause chain |
| `Format(err)` | Render error + resolution block for human display |

## License

Apache License 2.0. See the `LICENSE` file for details.
