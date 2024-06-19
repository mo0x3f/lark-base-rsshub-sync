package cache

import (
	"encoding/json"
	"errors"
	"testing"
	"time"
)

type person struct {
	Name string
	Age  int
}

func (p *person) Serialize() (string, error) {
	data, _ := json.Marshal(&p)
	return string(data), nil
}

func (p *person) Deserialize(val string) error {
	item := &person{}
	_ = json.Unmarshal([]byte(val), &item)
	p.Name = item.Name
	p.Age = item.Age
	return nil
}

func (p *person) Copy(from Serializer) error {
	fromObj, ok := from.(*person)
	if !ok {
		return errors.New("types err")
	}

	p.Name = fromObj.Name
	p.Age = fromObj.Age
	return nil
}

func Test_memoryCache(t *testing.T) {
	cache := &memoryCache{}
	_ = cache.Init()

	notFound := &person{}
	err := cache.Get("p1", notFound)
	if !errors.Is(err, ErrNotExists) {
		t.Errorf("err should be ErrNotExists")
	}

	p2 := &person{
		Name: "p2",
		Age:  2,
	}
	err = cache.Set("p2", p2, 2*time.Second)
	if err != nil {
		t.Errorf("err should not exist")
	}

	found := &person{}
	err = cache.Get("p2", found)
	if found.Name != "p2" {
		t.Errorf("found.Name != p2")
	}
	if found.Age != 2 {
		t.Errorf("found.Age != 2")
	}

	time.Sleep(2 * time.Second)

	fetch := func() (Serializer, error) {
		return &person{
			Name: "p3",
			Age:  3,
		}, nil
	}
	target := &person{}
	err = cache.GetAndRefresh("p3", target, fetch, 2*time.Second)
	if err != nil {
		t.Errorf("err should not exist")
	}
	if target.Name != "p3" {
		t.Errorf("target.Name != p3")
	}
	if target.Age != 3 {
		t.Errorf("target.Age != 3")
	}
}
