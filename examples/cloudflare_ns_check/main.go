package main

import (
	"fmt"
	"github.com/VisarBuza/taskflow/pkg/taskflow"
	"net"
	"strings"
)

type lookupFactory struct {
}

func (f *lookupFactory) Make(line string) taskflow.Task {
	return &lookup{name: line}
}

type lookup struct {
	name       string
	err        error
	cloudflare bool
}

func (t *lookup) Process() {
	nss, err := net.LookupNS(t.name)
	if err != nil {
		t.err = err
	} else {
		for _, ns := range nss {
			if strings.HasSuffix(ns.Host, ".ns.cloudflare.com.") {
				t.cloudflare = true
			}
		}
	}
}

func (t *lookup) Print() {
	state := "other"
	switch {
	case t.err != nil:
		state = "error"
	case t.cloudflare:
		state = "cloudflare"
	}

	fmt.Printf("%s: %s\n", t.name, state)
}

func main() {
	f := &lookupFactory{}
	taskflow.Run(f)
}
