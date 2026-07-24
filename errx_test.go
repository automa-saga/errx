// SPDX-License-Identifier: Apache-2.0

package errx_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/joomcode/errorx"

	"github.com/automa-saga/errx"
)

var (
	testNS   = errorx.NewNamespace("errx_test")
	testType = testNS.NewType("boom")
)

func TestWithReasonAndHints(t *testing.T) {
	err := errx.WithHints(
		errx.WithReason(testType.New("something broke"), "SomethingBroke"),
		"Try turning it off and on again",
		"Check the logs",
	)

	reason, ok := errx.ReasonOf(err)
	if !ok || reason != "SomethingBroke" {
		t.Fatalf("Reason() = %q, %v; want SomethingBroke, true", reason, ok)
	}
	hints, ok := errx.Hints(err)
	if !ok || len(hints) != 2 || hints[0] != "Try turning it off and on again" {
		t.Fatalf("Hints() = %v, %v", hints, ok)
	}
}

func TestReasonWalksCauseChain(t *testing.T) {
	inner := errx.WithReason(testType.New("inner"), "InnerReason")
	outer := testType.Wrap(inner, "outer context")

	reason, ok := errx.ReasonOf(outer)
	if !ok || reason != "InnerReason" {
		t.Fatalf("Reason() through chain = %q, %v; want InnerReason, true", reason, ok)
	}
}

func TestReasonSurvivesFmtErrorfWrapper(t *testing.T) {
	inner := errx.WithHints(errx.WithReason(testType.New("inner"), "InnerReason"), "do the thing")
	wrapped := fmt.Errorf("adding context: %w", inner)

	reason, ok := errx.ReasonOf(wrapped)
	if !ok || reason != "InnerReason" {
		t.Fatalf("Reason() through fmt.Errorf = %q, %v; want InnerReason, true", reason, ok)
	}
	hints, ok := errx.Hints(wrapped)
	if !ok || len(hints) != 1 || hints[0] != "do the thing" {
		t.Fatalf("Hints() through fmt.Errorf = %v, %v", hints, ok)
	}
}

func TestNilAndAbsent(t *testing.T) {
	if errx.WithReason(nil, "x") != nil {
		t.Fatal("WithReason(nil) should be nil")
	}
	if errx.WithHints(nil) != nil {
		t.Fatal("WithHints(nil) should be nil")
	}
	plain := errors.New("plain")
	if _, ok := errx.ReasonOf(plain); ok {
		t.Fatal("plain error should have no reason")
	}
	if _, ok := errx.Hints(plain); ok {
		t.Fatal("plain error should have no hints")
	}
}

func TestWithHintsNoHintsIsPassthrough(t *testing.T) {
	base := testType.New("x")
	if got := errx.WithHints(base); got != error(base) {
		t.Fatal("WithHints with no hints should return the error unchanged")
	}
}

func TestDecorate(t *testing.T) {
	err := errx.Decorate(testType.New("fail"), "DecorateReason", "hint one", "hint two")

	reason, ok := errx.ReasonOf(err)
	if !ok || reason != "DecorateReason" {
		t.Fatalf("Decorate Reason = %q, %v; want DecorateReason, true", reason, ok)
	}
	hints, ok := errx.Hints(err)
	if !ok || len(hints) != 2 {
		t.Fatalf("Decorate Hints = %v, %v; want 2 hints", hints, ok)
	}
}

func TestDecorateNil(t *testing.T) {
	if errx.Decorate(nil, "x", "hint") != nil {
		t.Fatal("Decorate(nil) should be nil")
	}
}

func TestFormat(t *testing.T) {
	if errx.Format(nil) != "" {
		t.Fatal("Format(nil) should be empty")
	}

	plain := testType.New("boom")
	if got := errx.Format(plain); got != plain.Error() {
		t.Fatalf("Format without hints = %q; want %q", got, plain.Error())
	}

	withHints := errx.WithHints(testType.New("boom"), "restart", "check config")
	out := errx.Format(withHints)
	if !contains(out, "Resolution:") || !contains(out, "restart") || !contains(out, "check config") {
		t.Fatalf("Format with hints = %q; want Resolution block", out)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
