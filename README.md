debug [![Build Status](https://drone.io/github.com/tango-contrib/debug/status.png)](https://drone.io/github.com/tango-contrib/debug/latest) [![](http://gocover.io/_badge/github.com/tango-contrib/debug)](http://gocover.io/github.com/tango-contrib/debug)
======

Middleware debug is a debug middleware for [Tango](https://github.com/lunny/tango). 

## Installation

    go get github.com/tango-contrib/debug

## Simple Example

```Go
type DebugAction struct {
    tango.Ctx
}

func (c *DebugAction) Get() {
    c.Write([]byte("get"))
}

func main() {
    t := tango.Classic()
    t.Use(events.Events())
    t.Get("/", new(EventAction))
    t.Run()
}
```