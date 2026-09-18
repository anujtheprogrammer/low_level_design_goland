package main

import (
	"errors"
	"fmt"
)

// ============================================================
// 1. THE STATE INTERFACE
// Every concrete state must implement this. Note that every
// method takes the Context (*Order) as a parameter — this is
// what lets a state reach back into the Context and swap
// itself out for the next state.
// ============================================================
type OrderState interface {
	Pay(o *Order) error
	Ship(o *Order) error
	Deliver(o *Order) error
	Cancel(o *Order) error
	Name() string // human redable name for state
}

// ============================================================
// 2. THE CONTEXT
// Order holds a reference to whichever state is "current".
// Crucially, Order itself has NO branching logic — it just
// forwards each call to o.state. This is the whole point of
// the pattern: behaviour lives in the states, not the Context.
// ============================================================
type Order struct {
	ID    string
	state OrderState
}

// NewOrder always starts life in NewState.
func NewOrder(id string) *Order {
	return &Order{ID: id, state: &NewState{}}
}

// SetState is called BY the state objects themselves to move
// the Context forward. It is the only place a transition
// actually happens.
func (o *Order) SetState(s OrderState) {
	fmt.Printf("{transition} order %s : %s -> %s\n", o.ID, o.state.Name(), s.Name())
	o.state = s
}

func (o *Order) CurrentState() string {
	return o.state.Name()
}

// These four methods are pure delegation — Order does not
// decide what is legal, the current state object does.
func (o *Order) Pay() error {
	return o.state.Pay(o)
}

func (o *Order) Ship() error {
	return o.state.Ship(o)
}

func (o *Order) Deliver() error {
	return o.state.Deliver(o)
}

func (o *Order) Cancel() error {
	return o.state.Cancel(o)
}

// ============================================================
// 3. CONCRETE STATES
// Each type below implements OrderState. Each one only knows
// about the transitions that are legal FROM that state.
// ============================================================

// ---- NewState: order created, awaiting payment ----
type NewState struct{}

func (s *NewState) Name() string {
	return "new"
}

func (s *NewState) Pay(o *Order) error {
	o.SetState(&PaidState{})
	return nil
}

func (s *NewState) Ship(o *Order) error {
	return errors.New("cannot ship : order has not been paid yet")
}

func (s *NewState) Deliver(o *Order) error {
	return errors.New("cannot deliver : order has not shipped yet")
}

func (s *NewState) Cancel(o *Order) error {
	o.SetState(&CancelledState{})
	return nil
}

// ---- PaidState: payment confirmed, awaiting shipment ----

type PaidState struct{}

func (s *PaidState) Name() string {
	return "Paid"
}

func (s *PaidState) Pay(o *Order) error {
	return errors.New("Order is laready paid")
}

func (s *PaidState) Ship(o *Order) error {
	o.SetState(&ShippedState{})
	return nil
}

func (s *PaidState) Deliver(o *Order) error {
	return errors.New("the order is o=not shipped yet")
}

func (s *PaidState) Cancel(o *Order) error {
	o.SetState(&CancelledState{})
	return nil
}

// ---- ShippedState: package is with the carrier ----
type ShippedState struct{}

func (s *ShippedState) Name() string {
	return "Shipping"
}

func (s *ShippedState) Pay(o *Order) error {
	return errors.New("The order is already paid")
}

func (s *ShippedState) Ship(o *Order) error {
	return errors.New("the order is already shipped")
}

func (s *ShippedState) Deliver(o *Order) error {
	o.SetState(&DeliveredState{})
	return nil
}

func (s *ShippedState) Cancel(o *Order) error {
	// Business rule: once a package has left the warehouse, it
	// can no longer be cancelled — this rule lives HERE, and
	// nowhere else in the codebase.
	return errors.New("cannot cancel: order has already shipped")
}

// ---- DeliveredState: terminal / happy-path end state ----
type DeliveredState struct{}

func (s *DeliveredState) Name() string           { return "Delivered" }
func (s *DeliveredState) Pay(o *Order) error     { return errors.New("the order is already paid") }
func (s *DeliveredState) Ship(o *Order) error    { return errors.New("The order is already shipped") }
func (s *DeliveredState) Deliver(o *Order) error { return errors.New("The order is delivered") }
func (s *DeliveredState) Cancel(o *Order) error  { return errors.New("cannot cancel : order delivered") }

// ---- CancelledState: terminal / unhappy-path end state ----

type CancelledState struct{}

func (s *CancelledState) Name() string           { return "Cancelled" }
func (s *CancelledState) Pay(o *Order) error     { return errors.New("order is cancelled") }
func (s *CancelledState) Ship(o *Order) error    { return errors.New("order is cancelled") }
func (s *CancelledState) Deliver(o *Order) error { return errors.New("order is cancelled") }
func (s *CancelledState) Cancel(o *Order) error  { return errors.New("order is already cancelled") }
