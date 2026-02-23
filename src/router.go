package main

import "net"

type Route struct {
	Prefix *net.IPNet
	OutIf  string
}

type Router struct {
	routes []Route
}

func (r *Router) AddRoute(cidr, outIf string) error {
	_, n, err := net.ParseCIDR(cidr)
	if err != nil {
		return err
	}
	r.routes = append(r.routes, Route{Prefix: n, OutIf: outIf})
	return nil
}

func (r *Router) LookupRoute(dst net.IP) string {
	for _, rt := range r.routes {
		if rt.Prefix.Contains(dst) {
			return rt.OutIf
		}
	}
	return ""
}
