package utils

import (
	"fmt"
	"reflect"
	"sync"
)

type ServiceScope string

const (
	Singleton ServiceScope = "singleton"
	Transient ServiceScope = "transient"
)

type ServiceDescriptor struct {
	ServiceType reflect.Type
	Factory     interface{}
	Instance    interface{}
	Scope       ServiceScope
}

type Container struct {
	services map[string]*ServiceDescriptor
	mu       sync.RWMutex
}

func NewContainer() *Container {
	return &Container{
		services: make(map[string]*ServiceDescriptor),
	}
}

func (c *Container) RegisterSingleton(serviceType interface{}, factory interface{}) {
	c.register(serviceType, factory, Singleton)
}

func (c *Container) RegisterTransient(serviceType interface{}, factory interface{}) {
	c.register(serviceType, factory, Transient)
}

func (c *Container) register(serviceType interface{}, factory interface{}, scope ServiceScope) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var key string
	var sType reflect.Type

	if t, ok := serviceType.(reflect.Type); ok {
		sType = t
		key = t.String()
	} else {
		sType = reflect.TypeOf(serviceType)
		if sType.Kind() == reflect.Ptr {
			sType = sType.Elem()
		}
		key = sType.String()
	}

	c.services[key] = &ServiceDescriptor{
		ServiceType: sType,
		Factory:     factory,
		Scope:       scope,
	}
}

func (c *Container) Resolve(serviceType interface{}) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var key string
	if t, ok := serviceType.(reflect.Type); ok {
		key = t.String()
	} else {
		sType := reflect.TypeOf(serviceType)
		if sType.Kind() == reflect.Ptr {
			sType = sType.Elem()
		}
		key = sType.String()
	}

	descriptor, exists := c.services[key]
	if !exists {
		return nil, fmt.Errorf("service %s not registered", key)
	}

	if descriptor.Scope == Singleton && descriptor.Instance != nil {
		return descriptor.Instance, nil
	}

	instance, err := c.createInstance(descriptor)
	if err != nil {
		return nil, err
	}

	if descriptor.Scope == Singleton {
		descriptor.Instance = instance
	}

	return instance, nil
}

func (c *Container) createInstance(descriptor *ServiceDescriptor) (interface{}, error) {
	factoryValue := reflect.ValueOf(descriptor.Factory)
	factoryType := factoryValue.Type()

	if factoryType.Kind() != reflect.Func {
		return nil, fmt.Errorf("factory must be a function")
	}

	// Get factory function parameters and resolve dependencies
	numIn := factoryType.NumIn()
	args := make([]reflect.Value, numIn)

	for i := 0; i < numIn; i++ {
		paramType := factoryType.In(i)
		dependency, err := c.Resolve(paramType)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve dependency %s: %v", paramType.String(), err)
		}
		args[i] = reflect.ValueOf(dependency)
	}

	results := factoryValue.Call(args)
	if len(results) == 0 {
		return nil, fmt.Errorf("factory function must return at least one value")
	}

	instance := results[0].Interface()
	
	// Check if factory returns error as second return value
	if len(results) > 1 && !results[1].IsNil() {
		if err, ok := results[1].Interface().(error); ok {
			return nil, err
		}
	}

	return instance, nil
}

func (c *Container) IsRegistered(serviceType interface{}) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var key string
	if t, ok := serviceType.(reflect.Type); ok {
		key = t.String()
	} else {
		sType := reflect.TypeOf(serviceType)
		if sType.Kind() == reflect.Ptr {
			sType = sType.Elem()
		}
		key = sType.String()
	}

	_, exists := c.services[key]
	return exists
}

// Global container instance
var DefaultContainer = NewContainer()