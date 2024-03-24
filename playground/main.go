package main

import (
	"fmt"
)

// Middleware defines the interface for middleware
type Middleware interface {
	Process(data string) string
}

// ConcreteMiddlewareA is a concrete implementation of Middleware that holds state
type ConcreteMiddlewareA struct {
	state int
}

// Process implements the Middleware interface with a pointer receiver
// allowing the method to mutate the struct's state
func (m *ConcreteMiddlewareA) Process(data string) string {
	m.state++ // Mutate state
	return fmt.Sprintf("A processed: %s, state: %d", data, m.state)
}

// ConcreteMiddlewareB another implementation of Middleware
type ConcreteMiddlewareB struct {
	prefix string
}

// Process for ConcreteMiddlewareB
func (m *ConcreteMiddlewareB) Process(data string) string {
	return fmt.Sprintf("%s B processed: %s", m.prefix, data)
}

// MiddlewareChain holds a sequence of Middleware
type MiddlewareChain struct {
	middlewares []Middleware
}

// AddMiddleware adds a new Middleware to the chain
func (chain *MiddlewareChain) AddMiddleware(m Middleware) {
	chain.middlewares = append(chain.middlewares, m)
}

// ProcessChain processes the data through the chain of middlewares
func (chain *MiddlewareChain) ProcessChain(data string) string {
	result := data
	for _, middleware := range chain.middlewares {
		result = middleware.Process(result)
	}
	return result
}

func main() {
	// Initialize middleware instances
	a := &ConcreteMiddlewareA{}
	b := &ConcreteMiddlewareB{prefix: "==> "}

	// Create a middleware chain and add middleware to it
	chain := MiddlewareChain{}
	chain.AddMiddleware(a)
	chain.AddMiddleware(b)

	// Process data through the chain
	processed := chain.ProcessChain("initial data")
	fmt.Println(processed)
}
