package LoadBalance

type Balancer interface {
    Get() (int, interface{})
    Put(id string, weight int, obj interface{}) string
    Remove(index int)
}
