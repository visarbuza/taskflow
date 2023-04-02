# TaskFlow

TaskFlow is a simple and efficient library for processing tasks concurrently in Go. It allows you to easily create, process, and print tasks in parallel, taking advantage of Go's goroutines and channels.


### Usage

First, import the TaskFlow package in your Go project:

```go
import "github.com/VisarBuza/taskflow/pkg/taskflow"
```

To use TaskFlow, you need to implement two interfaces:

1. `Factory`: This interface is responsible for creating new `Task` instances from input lines.
2. `Task`: This interface represents the tasks that will be processed concurrently. It includes methods for processing and printing the results.

Here's an example of how to use TaskFlow to check if a list of domain names use Cloudflare's nameservers:

You can also find the full example in the `examples` directory.

#### Define your Factory and Task implementations

```go
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

```

Run TaskFlow with your implementations

```go
func main() {
	f := &lookupFactory{}
	taskflow.Run(f)
}
```
