package model

import (
	"crawler/pkg/utils"
	"net/url"
	"strings"
	"sync"

	"github.com/go-rod/rod/lib/proto"
)

func longestCommonPrefix(a, b string) int {
	i := 0
	max := utils.Min(len(a), len(b))
	for i < max && a[i] == b[i] {
		i++
	}
	return i
}

type Tree struct {
	Method string `json:"method"`
	Root   *Node  `json:"root"`
}

type Trees []Tree

func (t Trees) Get(method string) *Node {
	for _, tree := range t {
		if tree.Method == method {
			return tree.Root
		}
	}
	return nil
}

func NewTree(method string, root *Node) Tree {
	return Tree{
		Method: method,
		Root:   root,
	}
}

type Node struct {
	path         string
	children     []*Node
	lock         sync.Mutex
	URL          url.URL
	ResourceType proto.NetworkResourceType
}

func newNode(path string, children []*Node, u url.URL, resourceType proto.NetworkResourceType) *Node {
	return &Node{
		path:         path,
		children:     children,
		URL:          u,
		ResourceType: resourceType,
	}
}

func (n *Node) Put(path string, u url.URL, resourceType proto.NetworkResourceType) {
	n.lock.Lock()
	defer n.lock.Unlock()
	if len(n.path) == 0 && len(n.children) == 0 {
		n.path = path
		n.URL = u
		n.ResourceType = resourceType
		return
	}

walk:
	for {
		i := longestCommonPrefix(path, n.path)

		if i < len(n.path) {
			n.children = []*Node{
				newNode(n.path[i:], n.children, n.URL, n.ResourceType),
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
			child := &Node{}

			n.children = append(n.children, child)
			n = child
			n.path = path
			n.URL = u
			n.ResourceType = resourceType
			return
		}
		n.URL = u
		n.ResourceType = resourceType
		return
	}
}

func (n *Node) Get(path string) *Node {
	return n.Find(n, path)
}

func (n *Node) Find(currentNode *Node, path string) *Node {
	if !strings.HasPrefix(path, currentNode.path) {
		return nil
	}

	if path == currentNode.path {
		return currentNode
	}

	path = path[len(currentNode.path):]

	for _, child := range currentNode.children {
		if result := n.Find(child, path); result != nil {
			return result
		}
	}
	return nil
}

func (n *Node) Children() []*Node {
	return n.children
}
