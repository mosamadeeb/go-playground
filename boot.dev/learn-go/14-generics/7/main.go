package main

import (
	"fmt"
)

// This is an interface that can be satisfied on a specific type (which should satisfy the customer constraint)
// That type can be either an interface that is a superset of customer (or customer itself), or a concrete type that satisfies customer
// Both userBiller and orgBiller satisfy this
type biller[C customer] interface {
	Charge(C) bill
	Name() string
}

// Using the customer interface directly means that neither userBiller nor orgBiller now satisfy this interface
// This is because they defined their methods on a concrete type that is not "customer"
// type biller[C customer] interface {
// 	Charge(customer) bill
// }

// Using the parametric constraint biller[C] allows us to have types that satisfy
// the interface *only* on a specific constraint
// So userBiller can Charge a user, but cannot charge an org
// vice versa for orgBiller

// don't edit below this line

type userBiller struct {
	Plan string
}

func (ub userBiller) Charge(u user) bill {
	amount := 50.0
	if ub.Plan == "pro" {
		amount = 100.0
	}
	return bill{
		Customer: u,
		Amount:   amount,
	}
}

func (sb userBiller) Name() string {
	return fmt.Sprintf("%s user biller", sb.Plan)
}

type orgBiller struct {
	Plan string
}

func (ob orgBiller) Name() string {
	return fmt.Sprintf("%s org biller", ob.Plan)
}

func (ob orgBiller) Charge(o org) bill {
	amount := 2000.0
	if ob.Plan == "pro" {
		amount = 3000.0
	}
	return bill{
		Customer: o,
		Amount:   amount,
	}
}

type customer interface {
	GetBillingEmail() string
}

type bill struct {
	Customer customer
	Amount   float64
}

type user struct {
	UserEmail string
}

func (u user) GetBillingEmail() string {
	return u.UserEmail
}

type org struct {
	Admin user
	Name  string
}

func (o org) GetBillingEmail() string {
	return o.Admin.GetBillingEmail()
}
