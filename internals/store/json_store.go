package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Iagobarros211256/rockshop/internals/models"
)

// group products and orders to serialization
type PersistedData struct {
	Products []models.Product `json:"products"`
	Orders   []models.Order   `json:"orders"`
}

// persists data on memory via json

// jsonstore class
type JSONStore struct {
	mu        sync.Mutex
	path      string
	products  map[uint]models.Product
	orders    map[uint]models.Order
	nextProd  uint
	nextOrder uint
}

// jsonstore class constructor
func NewJSONStore(path string) (*JSONStore, error) {
	s := &JSONStore{
		path:      path,
		products:  make(map[uint]models.Product),
		orders:    make(map[uint]models.Order),
		nextProd:  1,
		nextOrder: 1,
	}

	// create a dir if not exists
	//0755 Commonly used on web servers.
	// The owner can read, write, execute.
	// Everyone else can read and execute
	// but not modify the file.
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}

	//verifies if file is created on right place

	if _, err := os.Stat(path); err == nil {
		//read its contents
		b, err := os.ReadFile(path)
		// if file isnt created on right place this code will be activated
		if err != nil {
			return nil, err
		}
		//unpack json pack to form defined on persisted data structure
		if len(b) > 0 {
			var pdata PersistedData
			if err := json.Unmarshal(b, &pdata); err != nil {
				return nil, err
			}
			//looks for another product  or order to marshall
			for _, p := range pdata.Products {
				s.products[p.ID] = p
				if p.ID >= s.nextProd {
					s.nextProd = p.ID + 1
				}
			}

			for _, o := range pdata.Orders {
				s.orders[o.ID] = o
				if o.ID >= s.nextOrder {
					s.nextOrder = o.ID + 1
				}
			}
		}
	}
	//return it to screen
	return s, nil
}

// persisting data on memory
func (s *JSONStore) persist() error {
	pdata := PersistedData{
		Products: make([]models.Product, 0, len(s.products)),
		Orders:   make([]models.Order, 0, len(s.orders)),
	}
	for _, p := range s.products {
		pdata.Products = append(pdata.Products, p)
	}
	for _, o := range s.orders {
		pdata.Orders = append(pdata.Orders, o)
	}
	data, err := json.MarshalIndent(pdata, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)

}

// product  operations crud operations
func (s *JSONStore) ListProducts() []models.Product {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]models.Product, 0, len(s.products))
	for _, p := range s.products {
		out = append(out, p)
	}
	return out
}

func (s *JSONStore) GetProduct(id uint) (models.Product, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.products[id]
	return p, ok
}

func (s *JSONStore) CreateProduct(p models.Product) (models.Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p.ID = s.nextProd
	s.nextProd++
	s.products[p.ID] = p
	if err := s.persist(); err != nil {
		return models.Product{}, err
	}
	return p, nil
}

func (s *JSONStore) UpdateProduct(id uint, p models.Product) (models.Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.products[id]; !ok {
		return models.Product{}, errors.New("not found")
	}
	p.ID = id
	s.products[id] = p
	if err := s.persist(); err != nil {
		return models.Product{}, err
	}
	return p, nil

}
func (s *JSONStore) DeleteProduct(id uint) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.products[id]; !ok {
		return errors.New("not found")
	}

	delete(s.products, id)
	return s.persist()
}

// order operations create reduce stock if avaliable
func (s *JSONStore) ListOrders() []models.Order {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]models.Order, 0, len(s.orders))
	for _, o := range s.orders {
		out = append(out, o)
	}
	return out
}

func (s *JSONStore) GetOrder(id uint) (models.Order, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	o, ok := s.orders[id]
	return o, ok
}

func (s *JSONStore) CreateOrder(o models.Order) (models.Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// check stock and compute total
	var total int64 = 0
	for i, item := range o.Items {
		p, ok := s.products[item.ProductID]
		if !ok {
			return models.Order{}, errors.New("product not found")
		}
		if p.Stock < item.Qty {
			return models.Order{}, errors.New("insufficient stock for the product")
		}
		//set unit price from product at time of the order
		o.Items[i].UnitPrice = p.PriceCents
		total += p.PriceCents * int64(item.Qty)
	}
	//reduce stock
	for _, item := range o.Items {
		p := s.products[item.ProductID]
		p.Stock -= item.Qty
		s.products[item.ProductID] = p
	}
	//finalize order
	o.ID = s.nextOrder
	s.nextOrder++
	o.TotalCents = total
	o.CreatedAt = time.Now()
	s.orders[o.ID] = o
	if err := s.persist(); err != nil {
		return models.Order{}, err
	}
	return o, nil
}
