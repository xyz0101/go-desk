package lru

type Node struct {
	Key   string
	Value interface{}
	pre   *Node
	next  *Node
}

func (n *Node) Init(key string, value string) {
	n.Key = key
	n.Value = value
}

var head *Node
var end *Node

var limit int

type LRUCache struct {
	limit   int
	HashMap map[string]*Node
}

func GetLRUCache(limit int) *LRUCache {
	lruCache := LRUCache{limit: limit}
	lruCache.HashMap = make(map[string]*Node, limit)
	return &lruCache
}

func (l *LRUCache) Get(key string) interface{} {
	if v, ok := l.HashMap[key]; ok {
		l.refreshNode(v)
		return v.Value
	} else {
		return ""
	}
}

func (l *LRUCache) Put(key string, value interface{}) {
	if v, ok := l.HashMap[key]; !ok {
		if len(l.HashMap) >= l.limit {
			oldKey := l.removeNode(head)
			delete(l.HashMap, oldKey)
		}
		node := Node{Key: key, Value: value}
		l.addNode(&node)
		l.HashMap[key] = &node
	} else {
		v.Value = value
		l.refreshNode(v)
	}
}

func (l *LRUCache) refreshNode(node *Node) {
	if node == end {
		return
	}
	l.removeNode(node)
	l.addNode(node)
}

func (l *LRUCache) removeNode(node *Node) string {
	if node == end {
		end = end.pre
	} else if node == head {
		head = head.next
	} else {
		node.pre.next = node.next
		node.next.pre = node.pre
	}
	return node.Key
}

func (l *LRUCache) addNode(node *Node) {
	if end != nil {
		end.next = node
		node.pre = end
		node.next = nil
	}
	end = node
	if head == nil {
		head = node
	}
}
