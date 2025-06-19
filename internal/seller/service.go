package seller

type Service interface {
    GetByID(id string) (Seller, bool)
}

type service struct {
    store map[string]Seller
}

func NewService() Service {
    return &service{
        store: map[string]Seller{
            "seller1": {ID: "seller1", Name: "Ravi's Kitchen", Email: "ravi@kitchen.com"},
        },
    }
}

func (s *service) GetByID(id string) (Seller, bool) {
    seller, found := s.store[id]
    return seller, found
}
