package forward_must_edns0

import (
	"github.com/miekg/dns"
	"net"
	"time"
)

type Upstream struct {
	addr string
}

func NewUpstream(addr string) *Upstream {
	if _, _, err := net.SplitHostPort(addr); err != nil {
		addr = net.JoinHostPort(addr, "53")
	}
	return &Upstream{addr: addr}
}

func (u *Upstream) Exchange(m *dns.Msg) (*dns.Msg, error) {
	if m.IsEdns0() != nil {
		return u.exchangeOPTM(m)
	}
	mc := m.Copy()
	mc.SetEdns0(512, false)
	return u.exchangeOPTM(mc)
}

func (u *Upstream) exchangeOPTM(m *dns.Msg) (*dns.Msg, error) {
	c, err := dns.Dial("udp", u.addr)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	c.SetDeadline(time.Now().Add(time.Second * 3))
	if err := c.WriteMsg(m); err != nil {
		return nil, err
	}

	for {
		r, err := c.ReadMsg()
		if err != nil {
			return nil, err
		}
		if r.IsEdns0() == nil {
			continue
		}
		return r, nil
	}
}
