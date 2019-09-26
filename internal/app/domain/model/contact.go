package model

type Contact struct {
	Id    string
	Email string
	Data  string
}

func NewContact(id, email, data string) *Contact {
	return &Contact{
		Id:    id,
		Email: email,
		Data:  data,
	}
}
