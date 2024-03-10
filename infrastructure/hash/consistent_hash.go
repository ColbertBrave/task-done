package hash

import (
	"crypto/md5"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"sync"

	"github.com/spaolacci/murmur3"
)

const (
	// 最大权重
	topWeight = 100

	// 默认真实节点对应虚拟节点的最大副本数
	minReplicas = 100

	// 质数
	prime = 16777619
)

// HashFunc 哈希函数
type HashFunc func(data []byte) uint64

// ConsistentHash 一致性哈希环
type ConsistentHash struct {
	// 哈希函数
	hashFunc HashFunc

	// 虚拟节点放大因子，确定真实节点的虚拟节点数量
	maxReplicas int

	// 虚拟节点列表，哈希环
	virtualNodes []uint64

	// 虚拟节点到物理节点的映射
	virtualNodeMap map[uint64][]interface{}

	// 物理节点映射，快速判断是否存在node
	realNodesSet map[string]struct{}

	//读写锁
	lock sync.RWMutex
}

// 默认构造器
func NewConsistentHash() *ConsistentHash {
	return NewCustomConsistentHash(minReplicas, Hash)
}

// 有参构造器
func NewCustomConsistentHash(replicas int, hashFunc HashFunc) *ConsistentHash {
	if replicas < minReplicas {
		replicas = minReplicas
	}

	if hashFunc == nil {
		hashFunc = Hash
	}

	return &ConsistentHash{
		hashFunc:       hashFunc,
		maxReplicas:    replicas,
		virtualNodeMap: make(map[uint64][]interface{}),
		realNodesSet:   make(map[string]struct{}),
	}
}

// Get 根据v顺时针找到最近的虚拟节点，再通过虚拟节点映射找到真实节点
func (h *ConsistentHash) Get(v interface{}) (interface{}, bool) {
	h.lock.RLock()
	defer h.lock.RUnlock()

	// 当前哈希还没有物理节点
	if len(h.virtualNodeMap) == 0 {
		return nil, false
	}

	// 计算给定参数的哈希值
	hash := h.hashFunc([]byte(repr(v)))

	// 采用二分查找，查找大于给定参数哈希值的第一个虚拟节点
	// 因为每次添加节点后虚拟节点都会重新排序
	// 所以查询到的第一个节点就是我们的目标节点
	// 取余则可以实现环形列表效果，顺时针查找节点
	index := sort.Search(len(h.virtualNodes), func(i int) bool {
		return h.virtualNodes[i] >= hash
	}) % len(h.virtualNodes)

	// 将虚拟节点->物理节点映射
	realNodes := h.virtualNodeMap[h.virtualNodes[index]]
	switch len(realNodes) {
	// 不存在真实节点
	case 0:
		return nil, false
	// 只有一个真实节点，直接返回
	case 1:
		return realNodes[0], true
	// 存在多个真实节点，意味这出现哈希冲突
	default:
		// 此时我们对v重新进行哈希计算，对nodes长度取余得到一个新的index
		innerIndex := h.hashFunc([]byte(innerRepr(v)))
		pos := int(innerIndex % uint64(len(realNodes)))
		return realNodes[pos], true
	}
}

// Remove 删除物理节点
func (h *ConsistentHash) Remove(realNode interface{}) {
	h.lock.Lock()
	defer h.lock.Unlock()

	// 首先判断物理节点是否存在
	nodeRepr := repr(realNode)
	if !h.containsNode(nodeRepr) {
		return
	}

	// 移除虚拟节点映射
	for i := 0; i < h.maxReplicas; i++ {
		// 计算哈希值
		virtualNodeHash := h.hashFunc([]byte(nodeRepr + strconv.Itoa(i)))

		// 二分查找到第一个虚拟节点
		index := sort.Search(len(h.virtualNodes), func(i int) bool {
			return h.virtualNodes[i] >= virtualNodeHash
		})

		//切片删除对应的元素
		if index < len(h.virtualNodes) && h.virtualNodes[index] == virtualNodeHash {
			//定位到切片index之前的元素
			//将index之后的元素（index+1）前移覆盖index
			h.virtualNodes = append(h.virtualNodes[:index], h.virtualNodes[index+1:]...)
		}

		//虚拟节点删除映射
		h.removeVirtualNode(virtualNodeHash, nodeRepr)
	}

	// 删除真实节点
	h.removeNode(nodeRepr)
}

// removeVirtualNode 删除虚拟-真实节点映射关系
func (h *ConsistentHash) removeVirtualNode(virtualNode uint64, realNodeRepr string) {
	if realNodes, exist := h.virtualNodeMap[virtualNode]; exist {
		// 新建一个空的切片,容量与nodes保持一致
		newNodes := realNodes[:0]
		// 遍历nodes
		for _, realNode := range realNodes {
			//如果序列化值不相同，x是其他节点，不能删除
			if repr(realNode) != realNodeRepr {
				newNodes = append(newNodes, realNode)
			}
		}

		// 剩余节点不为空则重新绑定映射关系
		if len(newNodes) > 0 {
			h.virtualNodeMap[virtualNode] = newNodes
		} else {
			// 否则删除即可
			delete(h.virtualNodeMap, virtualNode)
		}
	}
}

// Add 扩容操作，增加物理节点
func (h *ConsistentHash) Add(node interface{}) {
	h.AddWithReplicas(node, h.maxReplicas)
}

// AddWithReplicas 扩容操作，增加物理节点
func (h *ConsistentHash) AddWithReplicas(node interface{}, replicas int) {
	// 支持可重复添加
	// 先执行删除操作
	h.Remove(node)

	// 不能超过放大因子上限
	replicas = max(h.maxReplicas, replicas)
	nodeRepr := repr(node)

	h.lock.Lock()
	defer h.lock.Unlock()

	// 添加node map映射
	h.addNode(nodeRepr)
	for i := 0; i < replicas; i++ {
		//创建虚拟节点，并添加至虚拟节点列表
		virtualNodeHash := h.hashFunc([]byte(nodeRepr + strconv.Itoa(i)))
		h.virtualNodes = append(h.virtualNodes, virtualNodeHash)

		// 映射虚拟节点-真实节点
		// 注意hashFunc可能会出现哈希冲突，所以采用的是追加操作
		// 虚拟节点-真实节点的映射对应的其实是个数组
		// 一个虚拟节点可能对应多个真实节点，当然概率非常小
		h.virtualNodeMap[virtualNodeHash] = append(h.virtualNodeMap[virtualNodeHash], node)
	}

	// 排序，后面会使用二分查找虚拟节点
	sort.Slice(h.virtualNodes, func(i, j int) bool {
		return h.virtualNodes[i] < h.virtualNodes[j]
	})
}

// AddWithWeight 按权重添加节点，通过权重来计算方法因子，最终控制虚拟节点的数量
// 权重越高，虚拟节点数量越多
func (h *ConsistentHash) AddWithWeight(node interface{}, weight int) {
	// don't need to make sure weight not larger than TopWeight,
	// because AddWithReplicas makes sure replicas cannot be larger than h.replicas
	replicas := h.maxReplicas * weight / topWeight
	h.AddWithReplicas(node, replicas)
}

func (h *ConsistentHash) addNode(nodeRepr string) {
	h.realNodesSet[nodeRepr] = struct{}{}
}

// 节点是否已存在
func (h *ConsistentHash) containsNode(nodeRepr string) bool {
	// realNodesSet 本身是一个map，时间复杂度是O1效率非常高
	_, ok := h.realNodesSet[nodeRepr]
	return ok
}

// 删除node
func (h *ConsistentHash) removeNode(nodeRepr string) {
	delete(h.realNodesSet, nodeRepr)
}

// 返回node的string值
// 在遇到哈希冲突时需要重新对key进行哈希计算
// 为了减少冲突的概率前面追加了一个质数 prime来减小冲突的概率
func innerRepr(v interface{}) string {
	return fmt.Sprintf("%d:%v", prime, v)
}

// Hash 返回data的哈希值
func Hash(data []byte) uint64 {
	return murmur3.Sum64(data)
}

// Md5 返回数据的MD5值
func Md5(data []byte) []byte {
	digest := md5.New()
	digest.Write(data)
	return digest.Sum(nil)
}

// Md5Hex 返回数据的MD5 HEX值
func Md5Hex(data []byte) string {
	return fmt.Sprintf("%x", Md5(data))
}

// Repr 将任意类型的值，以字符串表示
func repr(node interface{}) string {
	if node == nil {
		return ""
	}

	// if func (v *Type) String() string, we can't use Elem()
	switch vt := node.(type) {
	case fmt.Stringer:
		return vt.String()
	}

	val := reflect.ValueOf(node)
	if val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}

	return reprOfValue(val)
}

func reprOfValue(val reflect.Value) string {
	switch vt := val.Interface().(type) {
	case bool:
		return strconv.FormatBool(vt)
	case error:
		return vt.Error()
	case float32:
		return strconv.FormatFloat(float64(vt), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(vt, 'f', -1, 64)
	case fmt.Stringer:
		return vt.String()
	case int:
		return strconv.Itoa(vt)
	case int8:
		return strconv.Itoa(int(vt))
	case int16:
		return strconv.Itoa(int(vt))
	case int32:
		return strconv.Itoa(int(vt))
	case int64:
		return strconv.FormatInt(vt, 10)
	case string:
		return vt
	case uint:
		return strconv.FormatUint(uint64(vt), 10)
	case uint8:
		return strconv.FormatUint(uint64(vt), 10)
	case uint16:
		return strconv.FormatUint(uint64(vt), 10)
	case uint32:
		return strconv.FormatUint(uint64(vt), 10)
	case uint64:
		return strconv.FormatUint(vt, 10)
	case []byte:
		return string(vt)
	default:
		return fmt.Sprint(val.Interface())
	}
}
