//     Copyright (C) 2020-2021, IrineSistiana
//
//     This file is part of mosdns.
//
//     mosdns is free software: you can redistribute it and/or modify
//     it under the terms of the GNU General Public License as published by
//     the Free Software Foundation, either version 3 of the License, or
//     (at your option) any later version.
//
//     mosdns is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//     GNU General Public License for more details.
//
//     You should have received a copy of the GNU General Public License
//     along with this program.  If not, see <https://www.gnu.org/licenses/>.

package forward_must_edns0

import (
	"context"
	"errors"
	"github.com/IrineSistiana/mosdns/v3/dispatcher/handler"
	"github.com/IrineSistiana/mosdns/v3/dispatcher/pkg/bundled_upstream"
	"github.com/miekg/dns"
)

const PluginType = "forward_must_edns0"

func init() {
	handler.RegInitFunc(PluginType, Init, func() interface{} { return new(Args) })
}

var _ handler.ExecutablePlugin = (*forwardMustEDNS0)(nil)

type forwardMustEDNS0 struct {
	*handler.BP

	upstreams *bundled_upstream.BundledUpstream
}

type Args struct {
	UpstreamConfig []UpstreamConfig `yaml:"upstream"`
}

type UpstreamConfig struct {
	Addr    string `yaml:"addr"`
	Trusted bool   `yaml:"trusted"`
}

func Init(bp *handler.BP, args interface{}) (p handler.Plugin, err error) {
	return newForwarder(bp, args.(*Args))
}

func newForwarder(bp *handler.BP, args *Args) (*forwardMustEDNS0, error) {
	if len(args.UpstreamConfig) == 0 {
		return nil, errors.New("no upstream is configured")
	}

	f := new(forwardMustEDNS0)
	f.BP = bp

	bu := make([]bundled_upstream.Upstream, 0)
	for i, conf := range args.UpstreamConfig {
		if len(conf.Addr) == 0 {
			return nil, errors.New("missing upstream address")
		}

		if i == 0 { // Set first upstream as trusted upstream.
			conf.Trusted = true
		}

		bu = append(bu, &upstreamWrapper{
			u:       NewUpstream(conf.Addr),
			addr:    conf.Addr,
			trusted: conf.Trusted,
		})
	}

	f.upstreams = bundled_upstream.NewBundledUpstream(bu, bp.L())
	return f, nil
}

type upstreamWrapper struct {
	u       *Upstream
	addr    string
	trusted bool
}

func (u *upstreamWrapper) Address() string {
	return u.addr
}

func (u *upstreamWrapper) Exchange(ctx context.Context, q *dns.Msg) (*dns.Msg, error) {
	type res struct {
		r   *dns.Msg
		err error
	}
	resChan := make(chan *res, 1)

	qCopy := q.Copy()
	go func() {
		r := new(res)
		r.r, r.err = u.u.Exchange(qCopy)
		resChan <- r
	}()

	select {
	case r := <-resChan:
		return r.r, r.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (u *upstreamWrapper) Trusted() bool {
	return u.trusted
}

func (f *forwardMustEDNS0) Exec(ctx context.Context, qCtx *handler.Context, next handler.ExecutableChainNode) error {
	err := f.exec(ctx, qCtx)
	if err != nil {
		return err
	}

	return handler.ExecChainNode(ctx, qCtx, next)
}

func (f *forwardMustEDNS0) exec(ctx context.Context, qCtx *handler.Context) error {
	r, err := f.upstreams.ExchangeParallel(ctx, qCtx)
	if err != nil {
		qCtx.SetResponse(nil, handler.ContextStatusServerFailed)
		return err
	}
	qCtx.SetResponse(r, handler.ContextStatusResponded)
	return nil
}
