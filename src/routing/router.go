// NOTICE

// Project Name: Cloaq
// Copyright © 2026 Neil Talap and/or its designated Affiliates.

// This software is licensed under the Dragonfly Public License (DPL) 1.0.

// All rights reserved. The names "Neil Talap" and any associated logos or branding
// are trademarks of the Licensor and may not be used without express written permission,
// except as provided in Section 7 of the License.

// For commercial licensing inquiries or permissions beyond the scope of this
// license, please create an issue in github.

package routing

import (
	"fmt"
	"net"
)

type Route struct {
	Prefix *net.IPNet
	OutIf  string
}

type Router struct {
	routes []Route
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) AddRoute(cidr, outIf string) error {
	_, n, err := net.ParseCIDR(cidr)
	if err != nil {
		return err
	}
	r.routes = append(r.routes, Route{Prefix: n, OutIf: outIf})
	return nil
}

func (r *Router) LookupRoute(dst net.IP) (string, error) {
	for _, rt := range r.routes {
		if rt.Prefix.Contains(dst) {
			return rt.OutIf, nil
		}
	}
	return "", fmt.Errorf("no route found for %s", dst)
}
