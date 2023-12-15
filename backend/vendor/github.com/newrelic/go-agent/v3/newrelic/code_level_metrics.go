// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package newrelic

import (
	"errors"
	"reflect"
	"runtime"
	"strings"
	"sync"
)

//
// defaultAgentProjectRoot is the default filename pattern which is at
// the root of the agent's import path. This is used to identify functions
// on the call stack which are assumed to belong to the agent rather than
// the instrumented application's code.
//
const defaultAgentProjectRoot = "github.com/newrelic/go-agent/"

//
// CodeLocation marks the location of a line of source code for later reference.
//
type CodeLocation struct {
	// LineNo is the line number within the source file.
	LineNo int
	// Function is the function name (note that this may be auto-generated by Go
	// for function literals and the like). This is the fully-qualified name, which
	// includes the package name and other information to unambiguously identify
	// the function.
	Function string
	// FilePath is the absolute pathname on disk of the source file referred to.
	FilePath string
}

//
// CachedCodeLocation provides storage for the code location computed such that
// the discovery of the code location is only done once; thereafter the cached
// value is available for use.
//
// This type includes methods with the same names as some of the basic code location
// reporting functions and TraceOptions. However, when called as methods of a CachedCodeLocation value
// instead of a stand-alone function, the operation will make use of the cache to
// prevent computing the same source location more than once.
//
// A usable CachedCodeLocation value must be obtained via a call to
// NewCachedCodeLocation.
//
type CachedCodeLocation struct {
	location *CodeLocation
	once     *sync.Once
	err      error
}

//
// Err returns the error condition encountered when trying to determine
// the code location being cached, if any.
//
func (c *CachedCodeLocation) Err() error {
	if c == nil {
		return errors.New("nil CachedCodeLocation")
	}
	return c.err
}

//
// IsValid returns true if the cache value was correctly initialized
// (as by, for example, NewCachedCodeLocation), and therefore can be
// used to cache code location values. Otherwise it cannot be used.
//
func (c *CachedCodeLocation) IsValid() bool {
	return c != nil && c.once != nil
}

//
// NewCachedCodeLocation returns a pointer to a newly-created
// CachedCodeLocation value, suitable for use with the methods
// defined for that type.
//
func NewCachedCodeLocation() *CachedCodeLocation {
	return &CachedCodeLocation{
		once: new(sync.Once),
	}
}

type traceOptSet struct {
	LocationOverride *CodeLocation
	SuppressCLM      bool
	DemandCLM        bool
	IgnoredPrefixes  []string
	PathPrefixes     []string
	LocationCallback func() *CodeLocation
}

//
// TraceOption values provide optional parameters to transactions.
//
// (Currently it's only implemented for transactions, but the name TraceOption is
// intentionally generic in case we apply these to other kinds of traces in the future.)
//
type TraceOption func(*traceOptSet)

//
// WithCodeLocation adds an explicit CodeLocation value
// to report for the Code Level Metrics attached to a trace.
// This is probably a value previously obtained by calling
// ThisCodeLocation().
//
// Deprecated: This function requires the caller to do the work
// up-front to calculate the code location, which may be a waste
// of effort if code level metrics happens to be disabled. Instead,
// use the WithCodeLocationCallback function.
//
func WithCodeLocation(loc *CodeLocation) TraceOption {
	return func(o *traceOptSet) {
		o.LocationOverride = loc
	}
}

//
// WithCodeLocationCallback adds a callback function which the agent
// will call if it needs to report the code location with an explicit
// value provided by the caller. This will only be called if code
// level metrics is enabled, saving unnecessary work if those metrics
// are not enabled.
//
// If the callback function value passed here is nil, then no callback
// function will be used (same as if this function were never called).
// If the callback function itself returns nil instead of a pointer to
// a CodeLocation, then it is assumed the callback function was not able
// to determine the code location, and the CLM reporting code's normal
// method for determining the code location is used instead.
//
func WithCodeLocationCallback(locf func() *CodeLocation) TraceOption {
	return func(o *traceOptSet) {
		o.LocationCallback = locf
	}
}

//
// WithIgnoredPrefix indicates that the code location reported
// for Code Level Metrics should be the first function in the
// call stack that does not begin with the given string (or any of the given strings if more than one are given). This
// string is matched against the entire fully-qualified function
// name, which includes the name of the package the function
// comes from. By default, the Go Agent tries to take the first
// function on the call stack that doesn't seem to be internal to
// the agent itself, but you can control this behavior using
// this option.
//
// If all functions in the call stack begin with this prefix,
// the outermost one will be used anyway, since we didn't find
// anything better on the way to the bottom of the stack.
//
// If no prefix strings are passed here, the configured defaults will be used.
//
// Deprecated: New code should use WithIgnoredPrefixes instead.
//
func WithIgnoredPrefix(prefix ...string) TraceOption {
	return func(o *traceOptSet) {
		o.IgnoredPrefixes = prefix
	}
}

//
// WithIgnoredPrefixes indicates that the code location reported
// for Code Level Metrics should be the first function in the
// call stack that does not begin with the given string (or any of the given strings if more than one are given). This
// string is matched against the entire fully-qualified function
// name, which includes the name of the package the function
// comes from. By default, the Go Agent tries to take the first
// function on the call stack that doesn't seem to be internal to
// the agent itself, but you can control this behavior using
// this option.
//
// If all functions in the call stack begin with this prefix,
// the outermost one will be used anyway, since we didn't find
// anything better on the way to the bottom of the stack.
//
// If no prefix strings are passed here, the configured defaults will be used.
//
func WithIgnoredPrefixes(prefix ...string) TraceOption {
	return func(o *traceOptSet) {
		o.IgnoredPrefixes = prefix
	}
}

//
// WithPathPrefix overrides the list of source code path prefixes
// used to trim source file pathnames, providing a new set of one
// or more path prefixes to use for this trace only.
// If no strings are given, the configured defaults will be used.
//
// Deprecated: New code should use WithPathPrefixes instead.
//
func WithPathPrefix(prefix ...string) TraceOption {
	return func(o *traceOptSet) {
		o.PathPrefixes = prefix
	}
}

//
// WithPathPrefixes overrides the list of source code path prefixes
// used to trim source file pathnames, providing a new set of one
// or more path prefixes to use for this trace only.
// If no strings are given, the configured defaults will be used.
//
func WithPathPrefixes(prefix ...string) TraceOption {
	return func(o *traceOptSet) {
		o.PathPrefixes = prefix
	}
}

//
// WithoutCodeLevelMetrics suppresses the collection and reporting
// of Code Level Metrics for this trace. This helps avoid the overhead
// of collecting that information if it's not needed for certain traces.
//
func WithoutCodeLevelMetrics() TraceOption {
	return func(o *traceOptSet) {
		o.SuppressCLM = true
	}
}

//
// WithCodeLevelMetrics includes this trace in code level metrics even if
// it would otherwise not be (for example, if it would be out of the configured
// scope setting). This will never cause code level metrics to be reported if
// CLM were explicitly disabled (e.g. by CLM being globally off or if WithoutCodeLevelMetrics
// is present in the options for this trace).
//
func WithCodeLevelMetrics() TraceOption {
	return func(o *traceOptSet) {
		o.DemandCLM = true
	}
}

//
// WithThisCodeLocation is equivalent to calling WithCodeLocation, referring
// to the point in the code where the WithThisCodeLocation call is being made.
// This can be helpful, for example, when the actual code invocation which starts
// a transaction or other kind of trace is originating from a framework or other
// centralized location, but you want to report this point in your application
// for the Code Level Metrics associated with this trace.
//
func WithThisCodeLocation() TraceOption {
	return WithCodeLocation(ThisCodeLocation())
}

//
// WithThisCodeLocation is equivalent to the standalone WithThisCodeLocation
// TraceOption, but uses the cached value in its receiver to ensure that the
// overhead of computing the code location is only performed the first time
// it is invoked for each instance of the receiver variable.
//
func (c *CachedCodeLocation) WithThisCodeLocation() TraceOption {
	return WithCodeLocation(c.ThisCodeLocation())
}

//
// FunctionLocation is like ThisCodeLocation, but takes as its parameter
// a function value. It will report the code-level metrics information for
// that function if that is possible to do. It returns an error if it
// was not possible to get a code location from the parameter passed to it.
//
// If multiple functions are passed, each will be attempted until one is
// found for which we can successfully find a code location.
//
func FunctionLocation(functions ...interface{}) (*CodeLocation, error) {
	for _, function := range functions {
		if function == nil {
			continue
		}

		v := reflect.ValueOf(function)
		if !v.IsValid() || v.Kind() != reflect.Func {
			continue
		}

		if fInfo := runtime.FuncForPC(v.Pointer()); fInfo != nil {
			var loc CodeLocation

			loc.FilePath, loc.LineNo = fInfo.FileLine(fInfo.Entry())
			loc.Function = fInfo.Name()
			return &loc, nil
		}
	}

	return nil, errors.New("could not find code location for function")
}

//
// FunctionLocation works identically to the stand-alone FunctionLocation function,
// in that it determines the souce code location of the named function, returning
// a pointer to a CodeLocation value which represents that location, or an error value
// if it was unable to find a valid code location for the provided value. However,
// unlike the stand-alone function, this stores the result in the CachedCodeLocation receiver;
// thus, subsequent invocations of FunctionLocation for the same receiver will result in
// immediately repeating the value (or error, if applicable) obtained from the first
// invocation.
//
// This is thread-safe and is intended to allow the same code to run in multiple
// concurrent goroutines without needlessly recalculating the location of the
// function value.
//
func (c *CachedCodeLocation) FunctionLocation(functions ...interface{}) (*CodeLocation, error) {
	if c == nil || !c.IsValid() {
		// The cache is bogus so don't use it
		return FunctionLocation(functions...)
	}

	c.once.Do(func() {
		c.location, c.err = FunctionLocation(functions...)
	})
	return c.location, c.err
}

//
// WithFunctionLocation is like WithThisCodeLocation, but uses the
// function value(s) passed as the location to report. Unlike FunctionLocation,
// this does not report errors explicitly. If it is unable to use the
// value passed to find a code location, it will do nothing.
//
func WithFunctionLocation(functions ...interface{}) TraceOption {
	return func(o *traceOptSet) {
		loc, err := FunctionLocation(functions...)
		if err == nil {
			o.LocationOverride = loc
		}
	}
}

//
// WithFunctionLocation works like the standalone function WithFunctionLocation,
// but it stores a copy of the function's location in its receiver the first time
// it is used. Subsequently that cached value will be used instead of computing
// the source code location every time.
//
// This is thread-safe and is intended to allow the same code to run in multiple
// concurrent goroutines without needlessly recalculating the location of the
// function value.
//
func (c *CachedCodeLocation) WithFunctionLocation(functions ...interface{}) TraceOption {
	return func(o *traceOptSet) {
		loc, err := c.FunctionLocation(functions...)
		if err == nil {
			o.LocationOverride = loc
		}
	}
}

//
// WithDefaultFunctionLocation is like WithFunctionLocation but will only
// evaluate the location of the function if nothing that came before it
// set a code location first. This is useful, for example, if you want to
// provide a default code location value to be used but not pay the overhead
// of resolving that location until it's clear that you will need to. This
// should appear at the end of a TraceOption list (or at least after any
// other options that want to specify the code location).
//
func WithDefaultFunctionLocation(functions ...interface{}) TraceOption {
	return func(o *traceOptSet) {
		if o.LocationOverride == nil {
			WithFunctionLocation(functions...)(o)
		}
	}
}

//
// WithDefaultFunctionLocation works like the standalone WithDefaultFunctionLocation function,
// except that it takes a CachedCodeLocation receiver which will
// be used to cache the source code location of the function value.
//
// Thus, this will arrange for the given function to be reported in Code Level Metrics
// only if no other option that came before it gave an explicit location to use instead,
// but will also cache that answer in the provided CachedCodeLocation receiver variable, so that
// if called again with the same CachedCodeLocation variable, it will avoid the overhead
// of finding the function's location again, using instead the cached answer.
//
// This is thread-safe and is intended to allow the same code to run in multiple
// concurrent goroutines without needlessly recalculating the location of the
// function value.
//
// If an error is encountered when trying to evaluate the source code location of
// the provided function value, WithCachedDefaultFunctionLocation will not set anything
// for the reported code location, and the error will be available as a non-nil value
// in the Err member of the CachedCodeLocation variable.
// In this case, no additional attempts are guaranteed to be made on subsequent executions
// to determine the code location.
//
func (c *CachedCodeLocation) WithDefaultFunctionLocation(functions ...interface{}) TraceOption {
	return func(o *traceOptSet) {
		if o.LocationOverride == nil {
			loc, err := c.FunctionLocation(functions...)
			if err == nil {
				WithCodeLocation(loc)(o)
			}
		}
	}
}

//
// withPreparedOptions copies the option settings from a structure
// which was already set up (probably by executing a set of TraceOption
// functions already).
//
func withPreparedOptions(newOptions *traceOptSet) TraceOption {
	return func(o *traceOptSet) {
		if newOptions != nil {
			if newOptions.LocationOverride != nil {
				o.LocationOverride = newOptions.LocationOverride
			}
			if newOptions.LocationCallback != nil {
				o.LocationCallback = newOptions.LocationCallback
			}
			o.SuppressCLM = newOptions.SuppressCLM
			o.DemandCLM = newOptions.DemandCLM
			if newOptions.IgnoredPrefixes != nil {
				o.IgnoredPrefixes = newOptions.IgnoredPrefixes
			}
			if newOptions.PathPrefixes != nil {
				o.PathPrefixes = newOptions.PathPrefixes
			}
		}
	}
}

//
// ThisCodeLocation returns a CodeLocation value referring to
// the place in your code that it was invoked.
//
// With no arguments (or if passed a 0 value), it returns the location
// of its own caller. However, you may adjust this by passing the number
// of function calls to skip. For example, ThisCodeLocation(1) will return
// the CodeLocation of the place the current function was called from
// (i.e., the caller of the caller of ThisCodeLocation).
//
func ThisCodeLocation(skipLevels ...int) *CodeLocation {
	skip := 0
	if len(skipLevels) > 0 {
		skip += skipLevels[0]
	}

	return thisCodeLocationCommon(skip, false)
}

func thisCodeLocationCommon(skip int, skipInternal bool) *CodeLocation {
	var loc CodeLocation
	pcs := make([]uintptr, 20)

	depth := runtime.Callers(1, pcs)
	if depth > 0 {
		var frame runtime.Frame
		var clmFile string
		stillMore := true
		skipCLM := true

		frames := runtime.CallersFrames(pcs[:depth])
		for stillMore {
			frame, stillMore = frames.Next()
			//
			// We will begin here in our own CLM module code. We don't need to know
			// the IgnoredPrefix value since we can see the actual filename here now.
			// The first function in the call stack will be from here in the CLM code
			// so remember it for later.
			//
			if clmFile == "" {
				clmFile = frame.File
			}
			//
			// We need to skip over all the functions internal to this CLM module
			// to get to the function being reported.
			//
			if skipCLM && frame.File == clmFile {
				continue
			}
			skipCLM = false
			//
			// Now that we're past our CLM code, we might need to skip past an intermediary
			// set of calls that we entered (e.g., sync.(*Once)).
			//
			if skipInternal {
				if frame.File != clmFile {
					continue
				}
				//
				// past that, and back in our CLM code again, which we now also
				// need to skip as well.
				//
				skipInternal = false
				skipCLM = true
				continue
			}
			//
			// This should now be into the user's code. If they asked us to skip
			// over a number of calls, do that.
			//
			if skip > 0 {
				skip--
				continue
			}
			//
			// Finally, we have arrived at the target function.
			//
			stillMore = false
		}
		loc.LineNo = frame.Line
		loc.Function = frame.Function
		loc.FilePath = frame.File
	}

	return &loc
}

//
// ThisCodeLocation works identically to the stand-alone ThisCodeLocation function,
// in that it determines the souce code location from whence it was called, returning
// a pointer to a CodeLocation value which represents that location. However,
// unlike the stand-alone function, this stores the result in the CachedCodeLocation receiver;
// thus, subsequent invocations of ThisCodeLocation for the same receiver will result in
// immediately repeating the value obtained from the first
// invocation.
//
// This is thread-safe and is intended to allow the same code to run in multiple
// concurrent goroutines without needlessly recalculating the location of the
// caller.
//
func (c *CachedCodeLocation) ThisCodeLocation(skiplevels ...int) *CodeLocation {
	var skip int

	if len(skiplevels) > 0 {
		skip = skiplevels[0]
	}

	if c == nil || !c.IsValid() {
		// the cache is bogus so we can't use it.
		return thisCodeLocationCommon(skip, false)
	}

	c.once.Do(func() {
		c.location = thisCodeLocationCommon(skip, true)
		c.err = nil
	})
	return c.location
}

func removeCodeLevelMetrics(remAttr func(string)) {
	remAttr(AttributeCodeLineno)
	remAttr(AttributeCodeNamespace)
	remAttr(AttributeCodeFilepath)
	remAttr(AttributeCodeFunction)
}

//
// Evaluate a set of TraceOptions, returning a pointer to a new traceOptSet struct
// initialized from those options. To avoid any unnecessary performance penalties,
// if we encounter an option that suppresses CLM collection, we stop without evaluating
// anything further.
//
func resolveCLMTraceOptions(options []TraceOption) *traceOptSet {
	optSet := traceOptSet{}
	for _, o := range options {
		o(&optSet)
		if optSet.SuppressCLM {
			break
		}
	}
	return &optSet
}

func reportCodeLevelMetrics(tOpts traceOptSet, run *appRun, setAttr func(string, string, interface{})) {
	var location CodeLocation
	var locationp *CodeLocation

	if tOpts.LocationCallback != nil {
		locationp = tOpts.LocationCallback()
	} else {
		locationp = tOpts.LocationOverride
	}

	if locationp != nil {
		location = *locationp
	} else {
		pcs := make([]uintptr, 20)
		depth := runtime.Callers(2, pcs)
		if depth > 0 {
			frames := runtime.CallersFrames(pcs[:depth])
			moreToRead := true
			var frame runtime.Frame

			if tOpts.IgnoredPrefixes == nil {
				tOpts.IgnoredPrefixes = run.Config.CodeLevelMetrics.IgnoredPrefixes
				// for backward compatibility, add the singleton IgnoredPrefix if there is one
				if run.Config.CodeLevelMetrics.IgnoredPrefix != "" {
					tOpts.IgnoredPrefixes = append(tOpts.IgnoredPrefixes, run.Config.CodeLevelMetrics.IgnoredPrefix)
				}
				if tOpts.IgnoredPrefixes == nil {
					tOpts.IgnoredPrefixes = append(tOpts.IgnoredPrefixes, defaultAgentProjectRoot)
				}
			}

			// skip out to first non-agent frame, unless that IS the top-most frame
			for moreToRead {
				frame, moreToRead = frames.Next()
				if func() bool {
					for _, eachPrefix := range tOpts.IgnoredPrefixes {
						if strings.HasPrefix(frame.Function, eachPrefix) {
							return false
						}
					}
					return true
				}() {
					break
				}
			}

			location.FilePath = frame.File
			location.Function = frame.Function
			location.LineNo = frame.Line
		}
	}

	if tOpts.PathPrefixes == nil {
		tOpts.PathPrefixes = run.Config.CodeLevelMetrics.PathPrefixes
		// bring in a value still lingering in the deprecated PathPrefix field if the user put one there on their own
		if run.Config.CodeLevelMetrics.PathPrefix != "" {
			tOpts.PathPrefixes = append(tOpts.PathPrefixes, run.Config.CodeLevelMetrics.PathPrefix)
		}
	}

	// scan for any requested suppression of leading parts of file pathnames
	if tOpts.PathPrefixes != nil {
		for _, prefix := range tOpts.PathPrefixes {
			if pi := strings.Index(location.FilePath, prefix); pi >= 0 {
				location.FilePath = location.FilePath[pi:]
				break
			}
		}
	}

	ns := strings.LastIndex(location.Function, ".")
	function := location.Function
	namespace := ""

	if ns >= 0 {
		namespace = location.Function[:ns]
		function = location.Function[ns+1:]
	}

	// Impose data value size limits.
	// Report no field over 255 characters in length.
	// Report no CLM data at all if the function name is empty or >255 chars.
	// Report no CLM data at all if both namespace and file path are >255 chars.
	if function != "" && len(function) <= 255 && (len(namespace) <= 255 || len(location.FilePath) <= 255) {
		setAttr(AttributeCodeLineno, "", location.LineNo)
		setAttr(AttributeCodeFunction, function, nil)
		if len(namespace) <= 255 {
			setAttr(AttributeCodeNamespace, namespace, nil)
		}
		if len(location.FilePath) <= 255 {
			setAttr(AttributeCodeFilepath, location.FilePath, nil)
		}
	}
}
