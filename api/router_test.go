package api

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/n4wei/memo/lib/test_helper"
)

type dummyHandler1 struct{}

func (this *dummyHandler1) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

type dummyHandler2 struct{}

func (this *dummyHandler2) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

type dummyHandler3 struct{}

func (this *dummyHandler3) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

type dummyHandler4 struct{}

func (this *dummyHandler4) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

func setup() *router {
	return newRouter()
}

func printTree(root *node) {
	fmt.Println("root")
	cur := []*node{root}
	var next []*node

	for len(cur) > 0 {
		next = []*node{}

		for _, curNode := range cur {
			for key, nextNode := range curNode.keys {
				fmt.Printf("%v", string(key))
				if nextNode.isRoute {
					fmt.Printf("+")
				}

				next = append(next, nextNode)
			}
		}

		fmt.Println("")
		cur = next
	}
}

func TestRouter_InitialState(t *testing.T) {
	router := setup()
	rootNode := router.root
	var exist bool

	test_helper.AssertEqual(t, len(rootNode.keys), 0)
	test_helper.AssertEqual(t, rootNode.isRoute, false)
	test_helper.AssertEqual(t, rootNode.handler, nil)

	_, exist = router.getRoute("")
	test_helper.AssertEqual(t, exist, false)
	_, exist = router.getRoute("a")
	test_helper.AssertEqual(t, exist, false)
	_, exist = router.getRoute("/a")
	test_helper.AssertEqual(t, exist, false)
}

func TestRouter_InvalidRoutes(t *testing.T) {
	router := setup()
	var err error

	err = router.addRoute("", nil)
	test_helper.AssertError(t, err)
	err = router.addRoute("a", nil)
	test_helper.AssertError(t, err)
	err = router.addRoute("/a/", nil)
	test_helper.AssertError(t, err)
	err = router.addRoute("/a//b", nil)
	test_helper.AssertError(t, err)
}

func TestRouter_OverwritesExistingRoute(t *testing.T) {
	router := setup()
	dummyHandler1 := new(dummyHandler1)
	dummyHandler2 := new(dummyHandler2)
	var err error
	var handler http.Handler
	var exist bool

	err = router.addRoute("/a/b/c", dummyHandler1)
	test_helper.AssertNoError(t, err)
	handler, exist = router.getRoute("/a/b/c")
	test_helper.AssertEqual(t, exist, true)
	test_helper.AssertEqual(t, handler, dummyHandler1)

	err = router.addRoute("/a/b/c", dummyHandler2)
	test_helper.AssertNoError(t, err)
	handler, exist = router.getRoute("/a/b/c")
	test_helper.AssertEqual(t, exist, true)
	test_helper.AssertEqual(t, handler, dummyHandler2)
}

func TestRouter_SimpleRoutes(t *testing.T) {
	router := setup()
	dummyHandler1 := new(dummyHandler1)
	dummyHandler2 := new(dummyHandler2)
	dummyHandler3 := new(dummyHandler3)
	var err error
	var handler http.Handler
	var exist bool

	err = router.addRoute("/", dummyHandler1)
	test_helper.AssertNoError(t, err)
	err = router.addRoute("/abc", dummyHandler2)
	test_helper.AssertNoError(t, err)
	err = router.addRoute("/abc/def", dummyHandler3)
	test_helper.AssertNoError(t, err)

	handler, exist = router.getRoute("/")
	test_helper.AssertEqual(t, exist, true)
	test_helper.AssertEqual(t, handler, dummyHandler1)
	handler, exist = router.getRoute("/abc")
	test_helper.AssertEqual(t, exist, true)
	test_helper.AssertEqual(t, handler, dummyHandler2)
	_, exist = router.getRoute("/dummy")
	test_helper.AssertEqual(t, exist, false)
	handler, exist = router.getRoute("/abc/def")
	test_helper.AssertEqual(t, exist, true)
	test_helper.AssertEqual(t, handler, dummyHandler3)
	_, exist = router.getRoute("/abc/dummy")
	test_helper.AssertEqual(t, exist, false)
	_, exist = router.getRoute("/abc/def/dummy")
	test_helper.AssertEqual(t, exist, false)
}

func TestRouter_WildcardRoutes(t *testing.T) {
	router := setup()
	dummyHandler1 := new(dummyHandler1)
	dummyHandler2 := new(dummyHandler2)
	dummyHandler3 := new(dummyHandler3)
	dummyHandler4 := new(dummyHandler4)
	var err error
	var handler http.Handler
	var exist bool

	err = router.addRoute("/abc", dummyHandler1)
	test_helper.AssertEqual(t, err, nil)
	err = router.addRoute("/abc/:guid", dummyHandler2)
	test_helper.AssertEqual(t, err, nil)
	err = router.addRoute("/abc/:guid/def", dummyHandler3)
	test_helper.AssertEqual(t, err, nil)
	err = router.addRoute("/abc/:guid/def/:guid", dummyHandler4)
	test_helper.AssertEqual(t, err, nil)

	_, exist = router.getRoute("/")
	test_helper.AssertEqual(t, exist, false)

	handler, exist = router.getRoute("/abc")
	test_helper.AssertEqual(t, exist, true)
	test_helper.AssertEqual(t, handler, dummyHandler1)

	handler, exist = router.getRoute("/abc/some_guid")
	test_helper.AssertEqual(t, exist, true)
	test_helper.AssertEqual(t, handler, dummyHandler2)

	handler, exist = router.getRoute("/abc/some_guid/def")
	test_helper.AssertEqual(t, exist, true)
	test_helper.AssertEqual(t, handler, dummyHandler3)

	handler, exist = router.getRoute("/abc/some_guid/def/some_other_guid")
	test_helper.AssertEqual(t, exist, true)
	test_helper.AssertEqual(t, handler, dummyHandler4)
}
