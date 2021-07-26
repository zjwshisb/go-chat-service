package repositories

type repository interface {
	get() []interface{}
}

type Message struct {


}

func (Message *Message) get() []interface{}  {
	return []interface{}{
		1,
	}
}