# TaskFlow

TaskFlow is a simple and efficient library for processing tasks concurrently in Go. It allows you to easily create, process, and print tasks in parallel, taking advantage of Go's goroutines and channels.

### Inspired by

This library was inspired by the following talk:

- [dotGo 2014 - John Graham-Cumming](https://www.youtube.com/watch?v=woCg2zaIVzQ&ab_channel=dotconferences)

### Architecture Diagram

![Architecture of the program](assets/architecture.png)

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

### Running the example

```bash
go run examples/cloudflare_ns_check/main.go < examples/cloudflare_ns_check/records
```

**OUTPUT**
```txt
facebook.com: other
youtube.com: other
google.com: other
instagram.com: other
gjirafa.com: cloudflare
twitter.com: other
linkedin.com: other
vpapps.cloud: cloudflare
```

If you want to redirect the output to a file, you can use the following command:

```bash
go run examples/cloudflare_ns_check/main.go < examples/cloudflare_ns_check/records > output.txt
```
