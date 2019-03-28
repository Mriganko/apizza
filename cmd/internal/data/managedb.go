package data

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/harrybrwn/apizza/dawg"
	"github.com/harrybrwn/apizza/pkg/cache"
)

// OrderPrefix is the prefix added to user orders when stored in a database.
const OrderPrefix = "user_order_"

// PrintOrders will print all the names of the saved user orders
func PrintOrders(db cache.MapDB, w io.Writer) error {
	all, err := db.Map()
	if err != nil {
		return err
	}
	var orders []string

	for k := range all {
		if strings.Contains(k, OrderPrefix) {
			orders = append(orders, strings.Replace(k, OrderPrefix, "", -1))
		}
	}
	if len(orders) < 1 {
		fmt.Fprintln(w, "No orders saved.")
		return nil
	}
	fmt.Fprintln(w, "Your Orders:")
	for _, o := range orders {
		fmt.Fprintln(w, " ", o)
	}
	return nil
}

func GetOrder(name string, db cache.Getter) (*dawg.Order, error) {
	raw, err := db.Get(OrderPrefix + name)
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, fmt.Errorf("cannot find order %s", name)
	}
	order := &dawg.Order{}
	if err = json.Unmarshal(raw, order); err != nil {
		return nil, err
	}
	order.SetName(name)
	return order, nil
}

func SaveOrder(o *dawg.Order, db cache.Putter) error {
	raw, err := json.Marshal(o)
	if err != nil {
		return err
	}
	return db.Put(OrderPrefix+o.Name(), raw)
}