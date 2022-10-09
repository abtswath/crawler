package crawler

import (
	"crawler/pkg/model"
	"crawler/pkg/utils"
	"strings"
)

func longestCommonPrefix(a, b string) int {
	i := 0
	max := utils.Min(len(a), len(b))
	for i < max && a[i] == b[i] {
		i++
	}
	return i
}

type tree struct {
	method string
	root   *node
}

type trees []tree

func (t trees) get(method string) *node {
	for _, tree := range t {
		if tree.method == method {
			return tree.root
		}
	}
	return nil
}

type node struct {
	path     string
	children []*node
	request  model.Request
}

func (n *node) put(path string, request model.Request) {
	if len(n.path) == 0 && len(n.children) == 0 {
		n.path = path
		n.request = request
		return
	}

walk:
	for {
		i := longestCommonPrefix(path, n.path)

		if i < len(n.path) {
			n.children = []*node{
				{
					path:     n.path[i:],
					children: n.children,
					request:  n.request,
				},
			}
			n.path = path[:i]
		}

		if i < len(path) {
			path = path[i:]

			for _, child := range n.children {
				if longestCommonPrefix(child.path, path) > 0 {
					n = child
					continue walk
				}
			}
			child := &node{}

			n.children = append(n.children, child)
			n = child
			n.path = path
			n.request = request
			return
		}
		n.request = request
		return
	}
}

func (n *node) has(path string) bool {
	return n.get(path) != nil
}

func (n *node) get(path string) *node {
	return n.find(n, path)
}

func (n *node) find(currentNode *node, path string) *node {
	if !strings.HasPrefix(path, currentNode.path) {
		return nil
	}

	if path == currentNode.path {
		return currentNode
	}

	path = path[len(currentNode.path):]

	for _, child := range currentNode.children {
		if result := n.find(child, path); result != nil {
			return result
		}
	}
	return nil
}

func (n *node) all() []model.Request {
	return n.getAll(n)
}

func (n *node) getAll(curNode *node) []model.Request {
	if len(curNode.children) <= 0 {
		return []model.Request{curNode.request}
	}
	requests := []model.Request{}
	for _, child := range curNode.children {
		requests = append(requests, n.getAll(child)...)
	}
	return requests
}
