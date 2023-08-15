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

### Running the examples

#### Cloudflare NS Check
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

----

#### Creating a local database of comics from XKCD
```bash
seq 2500 | go run examples/xkcd_downlaoder/main.go > database
```
seq 2500 generates a sequence from 1 - 2500. 
The sequence represents the IDs of the XKCD comics.

**OUTPUT**
```txt
{"transcript":"((The comic is narrated by an unspecified person. All dialog is shown in boxes overlaid on the comic panels.))\n\n[[The panel background looks like a cloudy sky, with the clouds all running together and appearing as a blue\ngrey smear.]]\nI've always had trouble with the size of clouds.\nI \nknow\n they're huge. I can see their shapes.\nBut I don't really see them as objects on the same scale as trees and buildings.\nThey're a backdrop.\n\n[[A person stands on a flat disk inside a hemispherical dome with the front half cut away. The dome is labelled \"Sky\", and the disk is labelled \"Ground\". The dome is about twice as tall as the person.]]\nStars are the same way.\n\nI know they're scattered through and endless ocean, but my gut insists they're a painting on a domed ceiling.\n\n((The next two lines of dialog are stretched over the following three panels.))\n[[A person stands on a curved surface, looking up.]]\nIf I try hard enough, I get a glimmer of depth, a dizzying sense of space,\n\n[[The perspective of the scene shifts, suddenly the surface the person was standing on is in the top left of the panel. The person is now looking down, leaning back, and waving their arms trying to regain balance.]]\nBut then everything snaps back.\n\n[[The perspective of the scene returns to normal, the person is now semi-crouched, staring at the ground with legs spaced apart to help them balance.]]\n\n[[An american football field is shown, with sections at the tips of the goal posts highlighted and shown as a zoomed view in an insert box. The goal posts each have a webcam mounted on top of them.]]\nSo one summer afternoon\nI set up two HD webcams hundreds of feet apart,\nPointed them at the sky,\n\n((The next two lines of dialog are stretched over two panels each.))\n[[The first panel shows a pair of glasses with the note \"Very strong reading glasses.\" and a smartphone with an attachment designed to clip onto the glasses. The smartphone screen is setup to display two images side by side such that one camera is visible in the left half of the screen, and the other camera is visible in the right half of the screen.]]\nAnd fed one stream to each of my eyes.\n\n[[The next panel shows the completed phone\nglasses assembly.]]\nThe parallax expanded my depth perception by a thousand times,\n\n[[The person stands wearing the phone\nglasses assembly, staring into the sky.]]\nAnd I stood in my living room\nAt the bottom of an abyss\n\n[[The person now stands on the shore of an unidentified coastline (possibly Boston?), a city is near their right foot and the tallest skyscraper appears ankle high. A mountain range is behind them that is also barely ankle high. The person is standing with their head well above cloud level as clouds swim around them.]]\nWatching mountains drift by.\n\n{{Title text: I've looked at clouds from both sides now.}}","img":"https://imgs.xkcd.com/comics/depth_perception.png","title":"Depth Perception","safe_title":"Depth Perception","num":941,"day":"22","month":"8","year":"2011"}
{"transcript":"[[Two people are walking.  The first is wearing a white hat.]]\nSecond person: It just blows my mind. She seemed so genuine. I had no idea she was such a serial liar.\nSecond person: I just wish I had our six months back.\n\n[[The view focuses on the second person.]]\nSecond person: Her exes say the same thing happened to them.\nSecond person: Maybe what we need is a terrible-ex tracking and notification service.\n\n\n[[The second person turns, thoughtfully.]]\nFirst person: But after all the problems with sex offender registries, who would agree to run it?\nSecond person: Maybe one of the state governments more willing to experiment could try it out...\n\nSoon...\n[[Two people are sitting at a table, on which sit wine glasses and plates.  One has glasses and a goatee, and the other has long hair.  A person approaches them carrying a clipboard and a license.]]\nLicense person: Excuse me, ma'am.\nLong hair person: Yes?\nLicense person: This man is known to the state of California to be a total douchebag.\n\n{{Title text: Since the goatee, glasses, and Seltzer \u0026 Friedberg DVD collection didn't tip you off, there will be a $20 negligence charge for this service.}}","img":"https://imgs.xkcd.com/comics/bad_ex.png","title":"Bad Ex","safe_title":"Bad Ex","num":796,"day":"22","month":"9","year":"2010"}
{"transcript":"[[The panel is black with rough-edged white passages running down through it. A stick figure is holding onto a rope, dangling down one of these passages. White text is in the black sections.]]\nYou were afraid that you would disappear, that you would be lost and forgotten.\nI held you tight against the dark and said that I would always come for you.\nThen one day it happened. You were torn from my arms and vanished from this world.\nMaybe you don't remember my promise. But I meant every word.\nI hope you're not afraid, wherever you are.\nYou don't need to be.\nI'm not.\nI will find you.\n{{title-text: I'm like the Terminator, except with love!}}","img":"https://imgs.xkcd.com/comics/find_you.jpg","title":"Find You","safe_title":"Find You","num":104,"day":"19","month":"5","year":"2006"}
```
