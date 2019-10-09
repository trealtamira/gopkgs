package gcplogrus

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/errorreporting"
	joonix "github.com/joonix/log"
	log "github.com/sirupsen/logrus"
)

var (
	googleErrClient          *errorreporting.Client
	googleErrClientInitMutex = &sync.Mutex{}
)

// Opts are the options that must be given to Setup
type Opts struct {
	Levels []log.Level
	Async  bool

	// Error reporting timeout. This is ignored when Async is true.
	Timeout time.Duration

	// Google Cloud project ID
	Project string

	// The service name for error reports
	Service string

	// Set the logrus formatter to be compatible with Stackdriver Logging
	UseJoonixFormatter bool
}

// googleErrHook implements the logrus.Hook interface.
// See https://github.com/sirupsen/logrus/blob/master/hooks.go
type googleErrHook struct {
	opts Opts
}

// Levels returns the log (logrus) levels that will trigger the stackdriver error reporting callback.
func (hook *googleErrHook) Levels() []log.Level {
	list := hook.opts.Levels

	if len(list) == 0 {
		return []log.Level{log.FatalLevel}
	}
	return list
}

// Fire sends the log entry to Stackdriver and (eventually) flush the error client.
// WARNING: do not call logrus functions here.
func (hook *googleErrHook) Fire(e *log.Entry) error {
	if googleErrClient == nil {
		return nil
	}

	entry := errorreporting.Entry{Error: fmt.Errorf(e.Message)}

	if hook.opts.Async {
		googleErrClient.Report(entry)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), hook.opts.Timeout)
	defer cancel()

	err := googleErrClient.ReportSync(ctx, entry)
	defer googleErrClient.Flush()

	// the error returned here will be handled by OnError
	return err
}

// DefaultOpts returns the default options to be used when calling Setup
func DefaultOpts() Opts {
	return Opts{
		Levels:             []log.Level{log.FatalLevel},
		Async:              false,
		Timeout:            5 * time.Second,
		Project:            os.Getenv("GOOGLE_CLOUD_PROJECT"),
		Service:            os.Getenv("STACKDRIVER_SERVICE_NAME"),
		UseJoonixFormatter: true,
	}
}

// Setup intializes the Stackdriver error reporting client and add a hook to logrus.
func Setup(opts Opts) {
	// singleton behaviour
	googleErrClientInitMutex.Lock()
	defer googleErrClientInitMutex.Unlock()

	if googleErrClient != nil {
		return
	}

	ctx := context.Background()
	var err error

	// skip initialization when project or service are empty
	if opts.Project == "" || opts.Service == "" {
		log.Println("error reporting is disabled")
		return
	}

	// create the error reporting client
	googleErrClient, err = errorreporting.NewClient(ctx, opts.Project, errorreporting.Config{
		ServiceName: opts.Service,
		OnError: func(err error) {
			fmt.Fprintf(os.Stderr, "failed to report error: %v", err)
		},
	})

	if err != nil {
		panic(err)
	}

	// init logrus
	if opts.UseJoonixFormatter {
		log.SetFormatter(joonix.NewFormatter())
	}

	// add the logrus hook
	log.AddHook(&googleErrHook{opts: opts})
}
