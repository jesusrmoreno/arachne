package main

type GraphSubject struct {
	Value       Graph
	id          int
	Subscribers map[int]func(Graph)
}

func NewGraphSubject(v Graph) GraphSubject {
	return GraphSubject{
		Subscribers: map[int]func(Graph){},
		Value:       v,
		id:          0,
	}
}

func (gs *GraphSubject) notify() {
	for _, fn := range gs.Subscribers {
		fn(gs.Value)
	}
}

func (gs *GraphSubject) Next(v Graph) {
	gs.Value = v
	gs.notify()
}

func (gs *GraphSubject) Subscribe(fn func(g Graph)) func() {
	id := gs.id
	gs.Subscribers[gs.id] = fn
	gs.id += 1
	fn(gs.Value)
	return func() {
		delete(gs.Subscribers, id)
	}
}
