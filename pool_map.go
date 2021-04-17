package pool

import "strings"

var (
	poolMap = make(map[string]*Pool)
)

func register(p *Pool) {
	poolMap[strings.ToLower(p.name)] = p
}

func unRegister(name string) {
	delete(poolMap, strings.ToLower(name))
}

func Get(name string) *Pool {
	if p, ok := poolMap[strings.ToLower(name)]; ok {
		return p
	}
	return nil
}
