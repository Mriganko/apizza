package dawg

import (
	"testing"
)

func TestUserNearestStore(t *testing.T) {
	uname, pass, ok := gettestcreds()
	if !ok {
		t.Skip()
	}
	defer swapclient(10)()

	user, err := getTestUser(uname, pass)
	if err != nil {
		t.Error(err)
	}
	if user == nil {
		t.Fatal("user is nil")
	}
	if user.store != nil {
		t.Error("we should wait for the user to initialize this")
	}
	user.Addresses = []*UserAddress{}
	if user.DefaultAddress() != nil {
		t.Error("we just set this to an empty array, why is it not so")
	}
	user.AddAddress(testAddress())
	if user.DefaultAddress() == nil {
		t.Error("ok, we just added an address, why am i not getting one")
	}
	_, err = user.NearestStore(Delivery)
	if err != nil {
		t.Error(err)
	}
	if user.store == nil {
		t.Error("ok, now this variable should be stored")
	}
	s, err := user.NearestStore(Delivery)
	if err != nil {
		t.Error(err)
	}
	if s != user.store {
		t.Error("user.NearestStore should return the cached store on the second call")
	}
}

func TestUserStoresNearMe(t *testing.T) {
	uname, pass, ok := gettestcreds()
	if !ok {
		t.Skip()
	}
	defer swapclient(10)()

	user, err := getTestUser(uname, pass)
	if err != nil {
		t.Error(err)
	}
	if user == nil {
		t.Fatal("user should not be nil")
	}
	if err = user.SetServiceMethod("not correct"); err == nil {
		t.Error("expected error for an invalid service method")
	}
	if err != ErrBadService {
		t.Error("SetServiceMethod with bad val gave wrong error")
	}
	user.AddAddress(testAddress())
	stores, err := user.StoresNearMe()
	if err == nil {
		t.Error("expedted error")
	}
	if err != errNoServiceMethod {
		t.Error("wrong error")
	}
	if stores != nil {
		t.Error("should not have retured any stores")
	}

	if err = user.SetServiceMethod(Delivery); err != nil {
		t.Error(err)
	}
	addr := user.DefaultAddress()

	stores, err = user.StoresNearMe()
	if err != nil {
		t.Error(err)
	}
	for _, s := range stores {
		if s == nil {
			t.Error("should not have nil store")
		}
		if s.userAddress == nil {
			t.Fatal("nil store.userAddress")
		}
		if s.userService != user.ServiceMethod {
			t.Error("wrong service method")
		}
		if s.userAddress.City() != addr.City() {
			t.Error("wrong city")
		}
		if s.userAddress.LineOne() != addr.LineOne() {
			t.Error("wrong line one")
		}
		if s.userAddress.StateCode() != addr.StateCode() {
			t.Error("wrong state code")
		}
		if s.userAddress.Zip() != addr.Zip() {
			t.Error("wrong zip code")
		}
	}
}
