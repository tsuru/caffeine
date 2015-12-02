package main

import (
	"testing"

	"gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

type Suite struct{}

var _ = check.Suite(&Suite{})
