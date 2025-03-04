// Copyright (c) 2018 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package client

import (
	"time"

	"github.com/m3db/m3/src/x/clock"
	"github.com/m3db/m3/src/x/instrument"
	xio "github.com/m3db/m3/src/x/io"
	xnet "github.com/m3db/m3/src/x/net"
	"github.com/m3db/m3/src/x/retry"
	xtls "github.com/m3db/m3/src/x/tls"
)

const (
	defaultConnectionTimeout            = 1 * time.Second
	defaultConnectionKeepAlive          = true
	defaultReadTimeout                  = 15 * time.Second
	defaultWriteTimeout                 = 15 * time.Second
	defaultInitReconnectThreshold       = 1
	defaultMaxReconnectThreshold        = 4
	defaultReconnectThresholdMultiplier = 2
	defaultMaxReconnectDuration         = 20 * time.Second
	defaultWriteRetryInitialBackoff     = 0
	defaultWriteRetryBackoffFactor      = 2
	defaultWriteRetryMaxBackoff         = time.Second
	defaultWriteRetryMaxRetries         = 1
	defaultWriteRetryJitterEnabled      = true
)

// ConnectionOptions provides a set of options for tcp connections.
type ConnectionOptions interface {
	// SetInstrumentOptions sets the instrument options.
	SetClockOptions(value clock.Options) ConnectionOptions

	// ClockOptions returns the clock options.
	ClockOptions() clock.Options

	// SetInstrumentOptions sets the instrument options.
	SetInstrumentOptions(value instrument.Options) ConnectionOptions

	// InstrumentOptions returns the instrument options.
	InstrumentOptions() instrument.Options

	// SetConnectionTimeout sets the timeout for establishing connections.
	SetConnectionTimeout(value time.Duration) ConnectionOptions

	// ConnectionTimeout returns the timeout for establishing connections.
	ConnectionTimeout() time.Duration

	// SetConnectionKeepAlive sets the keepAlive for the connection.
	SetConnectionKeepAlive(value bool) ConnectionOptions

	// ConnectionKeepAlive returns the keepAlive for the connection.
	ConnectionKeepAlive() bool

	// SetReadTimeout sets the timeout for reading data.
	SetReadTimeout(value time.Duration) ConnectionOptions

	// ReadTimeout returns the timeout for reading data.
	ReadTimeout() time.Duration

	// SetWriteTimeout sets the timeout for writing data.
	SetWriteTimeout(value time.Duration) ConnectionOptions

	// WriteTimeout returns the timeout for writing data.
	WriteTimeout() time.Duration

	// SetInitReconnectThreshold sets the initial threshold for re-establshing connections.
	SetInitReconnectThreshold(value int) ConnectionOptions

	// InitReconnectThreshold returns the initial threshold for re-establishing connections.
	InitReconnectThreshold() int

	// SetMaxReconnectThreshold sets the max threshold for re-establishing connections.
	SetMaxReconnectThreshold(value int) ConnectionOptions

	// MaxReconnectThreshold returns the max threshold for re-establishing connections.
	MaxReconnectThreshold() int

	// SetReconnectThresholdMultiplier sets the threshold multiplier.
	SetReconnectThresholdMultiplier(value int) ConnectionOptions

	// ReconnectThresholdMultiplier returns the threshold multiplier.
	ReconnectThresholdMultiplier() int

	// SetMaxReconnectDuration sets the max duration between attempts to re-establish connections.
	SetMaxReconnectDuration(value time.Duration) ConnectionOptions

	// MaxReconnectDuration returns the max duration between attempts to re-establish connections.
	MaxReconnectDuration() time.Duration

	// SetWriteRetryOptions sets the retry options for retrying failed writes.
	SetWriteRetryOptions(value retry.Options) ConnectionOptions

	// WriteRetryOptions returns the retry options for retrying failed writes.
	WriteRetryOptions() retry.Options

	// SetRWOptions sets RW options.
	SetRWOptions(value xio.Options) ConnectionOptions

	// RWOptions returns the RW options.
	RWOptions() xio.Options

	// SetTLSOptions sets TLS options
	SetTLSOptions(value xtls.Options) ConnectionOptions

	// TLSOptions returns the TLS options
	TLSOptions() xtls.Options

	// ContextDialer allows customizing the way an aggregator client the aggregator, at the TCP layer.
	// By default, this is:
	// (&net.ContextDialer{}).DialContext. This can be used to do a variety of things, such as forwarding a connection
	// over a proxy.
	// NOTE: if your xnet.ContextDialerFn returns anything other a *net.TCPConn, TCP options such as KeepAlivePeriod
	// will *not* be applied automatically. It is your responsibility to make sure these get applied as needed in
	// your custom xnet.ContextDialerFn.
	ContextDialer() xnet.ContextDialerFn
	// SetContextDialer sets ContextDialer() -- see that method.
	SetContextDialer(dialer xnet.ContextDialerFn) ConnectionOptions
}

type connectionOptions struct {
	clockOpts      clock.Options
	instrumentOpts instrument.Options
	writeRetryOpts retry.Options
	rwOpts         xio.Options
	connTimeout    time.Duration
	readTimeout    time.Duration
	writeTimeout   time.Duration
	maxDuration    time.Duration
	initThreshold  int
	maxThreshold   int
	multiplier     int
	connKeepAlive  bool
	tlsOptions     xtls.Options
	dialer         xnet.ContextDialerFn
}

// NewConnectionOptions create a new set of connection options.
func NewConnectionOptions() ConnectionOptions {
	defaultWriteRetryOpts := retry.NewOptions().
		SetInitialBackoff(defaultWriteRetryInitialBackoff).
		SetBackoffFactor(defaultWriteRetryBackoffFactor).
		SetMaxBackoff(defaultWriteRetryMaxBackoff).
		SetMaxRetries(defaultWriteRetryMaxRetries).
		SetJitter(defaultWriteRetryJitterEnabled)
	return &connectionOptions{
		clockOpts:      clock.NewOptions(),
		instrumentOpts: instrument.NewOptions(),
		connTimeout:    defaultConnectionTimeout,
		connKeepAlive:  defaultConnectionKeepAlive,
		readTimeout:    defaultReadTimeout,
		writeTimeout:   defaultWriteTimeout,
		initThreshold:  defaultInitReconnectThreshold,
		maxThreshold:   defaultMaxReconnectThreshold,
		multiplier:     defaultReconnectThresholdMultiplier,
		maxDuration:    defaultMaxReconnectDuration,
		writeRetryOpts: defaultWriteRetryOpts,
		tlsOptions:     xtls.NewOptions(),
		rwOpts:         xio.NewOptions(),
		dialer:         nil, // Will default to net.Dialer{}.DialContext
	}
}

func (o *connectionOptions) SetClockOptions(value clock.Options) ConnectionOptions {
	opts := *o
	opts.clockOpts = value
	return &opts
}

func (o *connectionOptions) ClockOptions() clock.Options {
	return o.clockOpts
}

func (o *connectionOptions) SetInstrumentOptions(value instrument.Options) ConnectionOptions {
	opts := *o
	opts.instrumentOpts = value
	return &opts
}

func (o *connectionOptions) InstrumentOptions() instrument.Options {
	return o.instrumentOpts
}

func (o *connectionOptions) SetConnectionTimeout(value time.Duration) ConnectionOptions {
	opts := *o
	opts.connTimeout = value
	return &opts
}

func (o *connectionOptions) ConnectionTimeout() time.Duration {
	return o.connTimeout
}

func (o *connectionOptions) SetConnectionKeepAlive(value bool) ConnectionOptions {
	opts := *o
	opts.connKeepAlive = value
	return &opts
}

func (o *connectionOptions) ConnectionKeepAlive() bool {
	return o.connKeepAlive
}

func (o *connectionOptions) SetReadTimeout(value time.Duration) ConnectionOptions {
	opts := *o
	opts.readTimeout = value
	return &opts
}

func (o *connectionOptions) ReadTimeout() time.Duration {
	return o.readTimeout
}

func (o *connectionOptions) SetWriteTimeout(value time.Duration) ConnectionOptions {
	opts := *o
	opts.writeTimeout = value
	return &opts
}

func (o *connectionOptions) WriteTimeout() time.Duration {
	return o.writeTimeout
}

func (o *connectionOptions) SetInitReconnectThreshold(value int) ConnectionOptions {
	opts := *o
	opts.initThreshold = value
	return &opts
}

func (o *connectionOptions) InitReconnectThreshold() int {
	return o.initThreshold
}

func (o *connectionOptions) SetMaxReconnectThreshold(value int) ConnectionOptions {
	opts := *o
	opts.maxThreshold = value
	return &opts
}

func (o *connectionOptions) MaxReconnectThreshold() int {
	return o.maxThreshold
}

func (o *connectionOptions) SetReconnectThresholdMultiplier(value int) ConnectionOptions {
	opts := *o
	opts.multiplier = value
	return &opts
}

func (o *connectionOptions) ReconnectThresholdMultiplier() int {
	return o.multiplier
}

func (o *connectionOptions) SetMaxReconnectDuration(value time.Duration) ConnectionOptions {
	opts := *o
	opts.maxDuration = value
	return &opts
}

func (o *connectionOptions) MaxReconnectDuration() time.Duration {
	return o.maxDuration
}

func (o *connectionOptions) SetWriteRetryOptions(value retry.Options) ConnectionOptions {
	opts := *o
	opts.writeRetryOpts = value
	return &opts
}

func (o *connectionOptions) WriteRetryOptions() retry.Options {
	return o.writeRetryOpts
}

func (o *connectionOptions) SetRWOptions(value xio.Options) ConnectionOptions {
	opts := *o
	opts.rwOpts = value
	return &opts
}

func (o *connectionOptions) RWOptions() xio.Options {
	return o.rwOpts
}

func (o *connectionOptions) SetTLSOptions(value xtls.Options) ConnectionOptions {
	opts := *o
	opts.tlsOptions = value
	return &opts
}

func (o *connectionOptions) TLSOptions() xtls.Options {
	return o.tlsOptions
}

func (o *connectionOptions) ContextDialer() xnet.ContextDialerFn {
	return o.dialer
}

// SetContextDialer see ContextDialer.
func (o *connectionOptions) SetContextDialer(dialer xnet.ContextDialerFn) ConnectionOptions {
	opts := *o
	opts.dialer = dialer
	return &opts
}
