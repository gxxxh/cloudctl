package controllers

type Info struct {
	Namespace string
	Name      string
}

type IDCache map[string]Info

func NewIDCache() IDCache {
	return make(map[string]Info)
}

func (c IDCache) Add(id, namespace, name string) {
	if id == "" {
		return
	}
	c[id] = Info{
		Namespace: namespace,
		Name:      name,
	}
}

func (c IDCache) Delete(id string) {
	delete(c, id)
}

func (c IDCache) Get(id string) (namespace, name string) {
	info, ok := c[id]
	if !ok {
		namespace = ""
		name = ""
		return
	}
	namespace = info.Namespace
	name = info.Name
	return
}
