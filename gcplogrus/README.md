# gcplogrus

[logrus](https://github.com/sirupsen/logrus) integration for [Stackdriver](https://cloud.google.com/stackdriver/).

### Features:

* Logs are compatible with [Stackdriver Logging](https://cloud.google.com/logging/).
* Errors are sent to [Stackdriver Error Reporting](https://cloud.google.com/error-reporting/).

### Usage: 

The logging system can be initialized with a single line of code:

```go
gcplogrus.Setup(gcplogrus.DefaultOpts())
```

As long as the following ENV variables are set:

* ERRORREPORTING_PROJECT
* ERRORREPORTING_SERVICE

Complete example:

```go
  import (
    "github.com/trealtamira/gcplogrus"
    log "github.com/sirupsen/logrus"
  )

  func init() {
    gcplogrus.Setup(gcplogrus.DefaultOpts())
  }

  func main() {
    // The following log will be formatted in a way compatible
    // with Stackdriver Logging
    log.Infof("foobar")

    // The following error will be sent to Stackdriver Error Reporting.
    // By default, only FATAL errors are sent.
    log.Fatalf("this is a test error")
  }
```
