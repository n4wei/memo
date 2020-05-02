package api

import (
	"errors"
	"net/http"
	"strings"
)

type node struct {
	keys    map[byte]*node
	isRoute bool
	handler http.Handler
}

func newNode() *node {
	return &node{
		keys: map[byte]*node{},
	}
}

type router struct {
	root *node
}

func newRouter() *router {
	return &router{
		root: newNode(),
	}
}

func (this *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, exist := this.getRoute(r.URL.Path)
	if !exist {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	handler.ServeHTTP(w, r)
}

func (this *router) addRoute(path string, handler http.Handler) error {
	if len(path) == 0 {
		return errors.New("invalid path: path is empty")
	}
	if path[0] != '/' {
		return errors.New("invalid path: path must start with '/'")
	}
	if len(path) > 1 && path[len(path)-1] == '/' {
		return errors.New("invalid path: path must not end with '/', unless entire path is '/'")
	}

	path = path[1:]
	parts := strings.Split(path, "/")

	if len(parts) > 1 { // '/' is valid
		for _, part := range parts {
			if len(part) == 0 {
				return errors.New("invalid path: one or more contiguous '/' not allowed")
			}
		}
	}

	curNode := createOrAssign('/', this.root)

	for i, part := range parts {
		if len(part) == 0 {
			continue
		}

		if part[0] == ':' {
			curNode = createOrAssign('*', curNode)
		} else {
			for j := 0; j < len(part); j++ {
				curNode = createOrAssign(part[j], curNode)
			}
		}

		if i != len(parts)-1 {
			curNode = createOrAssign('/', curNode)
		}
	}

	curNode.isRoute = true
	curNode.handler = handler
	return nil
}

func createOrAssign(char byte, curNode *node) *node {
	nextNode, exist := curNode.keys[char]
	if !exist {
		nextNode = newNode()
		curNode.keys[char] = nextNode
	}
	return nextNode
}

func (this *router) getRoute(path string) (http.Handler, bool) {
	if len(this.root.keys) == 0 || len(path) == 0 || path[0] != '/' {
		return nil, false
	}

	curNode := this.root
	curNode = curNode.keys['/']

	if len(path) == 1 {
		return curNode.handler, curNode.isRoute
	}

	path = path[1:]
	parts := strings.Split(path, "/")
	var startNode *node
	var exist bool

	for i, part := range parts {
		startNode = curNode

		for j := 0; j < len(part); j++ {
			curNode, exist = curNode.keys[part[j]]
			if !exist {
				curNode, exist = startNode.keys['*']
				if !exist {
					return nil, false
				}
				break
			}
		}

		if i != len(parts)-1 {
			curNode, exist = curNode.keys['/']
			if !exist {
				return nil, false
			}
		}
	}

	return curNode.handler, curNode.isRoute
}
