package api

import "net/http"

type node struct {
	chars   map[byte]*node
	isRoute bool
	handler http.Handler
}

type router struct {
	root *node
}

func New() *router {
	return &router{}
}

func (this *router) RegisterRoute(route string, handler http.Handler) {
	cur := this.root
	var char byte

	for i := 0; i < len(route); i++ {
		char = route[i]
		next, exist := cur.chars[char]
		if !exist {
			next = &node{
				chars: map[byte]*node{},
			}
			cur.chars[char] = next
		}
		cur = next
	}

	cur.isRoute = true
	cur.handler = handler
}

func (this *router) FindRoute(route string) (http.Handler, bool) {
	cur := this.root
	var char byte

	for i := 0; i < len(route); i++ {
		char = route[i]
		next, exist := cur.chars[char]
		if !exist {
			return nil, false
		}
		cur = next
	}

	if !cur.isRoute {
		return nil, false
	}
	return cur.handler, true
}
