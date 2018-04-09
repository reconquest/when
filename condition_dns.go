package main

import (
	resolver "github.com/bogdanovich/dns_resolver"
)

type ConditionDNS struct {
	address  string
	domain   string
	ready    bool
	resolver *resolver.DnsResolver
}

func NewConditionDNS(address string, domain string) *ConditionDNS {
	resolver := resolver.New([]string{address})
	resolver.RetryTimes = 1

	return &ConditionDNS{
		address:  address,
		domain:   domain,
		resolver: resolver,
	}
}

func (dns *ConditionDNS) Ready() bool {
	return dns.ready
}

func (dns *ConditionDNS) Check() error {
	_, err := dns.resolver.LookupHost(dns.domain)
	if err != nil {
		return err
	}

	dns.ready = true

	return nil
}
