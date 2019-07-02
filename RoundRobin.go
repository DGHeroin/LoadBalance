package LoadBalance

import (
    "github.com/google/uuid"
    "sync"
)

//
// 带权值的rr算法
type RoundRobin struct {
    index int
    cw    int
    gcd   int
    max   int
    mutex sync.RWMutex
    items map[int]*RoundRobinItem
}

type RoundRobinItem struct {
    id     string
    weight int
    obj    interface{}
}

func NewRoundRobin() Balancer {
    return &RoundRobin{
        items: make(map[int]*RoundRobinItem),
        index: -1,
    }
}

func (r *RoundRobin) Get() (int, interface{}) {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    if len(r.items) == 0 {
        return -1, nil
    }
    for {
        r.index = (r.index + 1) % len(r.items)
        if r.index == 0 {
            r.cw = r.cw - r.gcd
            if r.cw <= 0 {
                r.cw = r.max
            }
        }
        if item, ok := r.items[r.index]; ok && item.weight >= r.cw {
            return r.index, item.obj
        }
    }
}

func (r *RoundRobin) Put(id string, weight int, obj interface{}) (string) {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    if id == "" {
        id = uuid.New().String()
    }
    r.items[len(r.items)] = &RoundRobinItem{
        weight: weight,
        obj:    obj,
    }
    // 计算当前数据的最大公约数
    var weights []int
    for _, v := range r.items {
        weights = append(weights, v.weight)
    }
    l := len(weights)
    a := weights[0]
    t := 0
    for i := 0; i < l-1; i++ {
        t = 0
        for weights[i+1] != 0 {
            t = weights[i+1]
            weights[i+1] = a % weights[i+1]
            a = t
        }
    }
    r.gcd = a
    // 计算最大权值
    r.max = 0
    for _, v := range r.items {
        if weight := v.weight; weight >= r.max {
            r.max = weight
        }
    }
    return id
}

// 当前数量
func (r *RoundRobin) Count() int {
    r.mutex.RLock()
    n := len(r.items)
    r.mutex.RUnlock()
    return n
}

//
func (r *RoundRobin) Remove(index int) {
    r.mutex.Lock()
    _, ok := r.items[index]
    r.mutex.Unlock()
    if !ok {
        return
    }
    delete(r.items, index)
    var temp []*RoundRobinItem
    for _, v := range r.items {
        temp = append(temp, v)
    }
    r.items = make(map[int]*RoundRobinItem)
    // 重新加入
    for _, v := range temp {
        r.Put(v.id, v.weight, v.obj)
    }
}
