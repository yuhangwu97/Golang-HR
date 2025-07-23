package utils

import (
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
)

// ControllerFactory creates controller instances using dependency injection
type ControllerFactory struct {
	container *Container
}

func NewControllerFactory(container *Container) *ControllerFactory {
	return &ControllerFactory{container: container}
}

// CreateHandler creates a Gin handler function that resolves a controller and calls a method
func (cf *ControllerFactory) CreateHandler(controllerType reflect.Type, methodName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Debug output
		fmt.Printf("DEBUG: Attempting to resolve controller type: %s\n", controllerType.String())

		// Resolve controller from container
		controller, err := cf.container.Resolve(controllerType)
		if err != nil {
			fmt.Printf("DEBUG: Failed to resolve controller: %v\n", err)
			ErrorResponse(c, 500, "Failed to resolve controller: "+err.Error())
			return
		}

		// Get method from controller
		controllerValue := reflect.ValueOf(controller)
		method := controllerValue.MethodByName(methodName)
		if !method.IsValid() {
			ErrorResponse(c, 500, "Method "+methodName+" not found on controller")
			return
		}

		// Call the method with gin context
		method.Call([]reflect.Value{reflect.ValueOf(c)})
	}
}

// CreateHandlerFunc creates a handler function for a specific controller method
func CreateHandlerFunc[T any](container *Container, methodName string) gin.HandlerFunc {
	controllerType := reflect.TypeOf((*T)(nil))
	factory := NewControllerFactory(container)
	return factory.CreateHandler(controllerType, methodName)
}

// InjectableHandler creates a handler that automatically injects dependencies
func InjectableHandler(container *Container, handler interface{}) gin.HandlerFunc {
	handlerValue := reflect.ValueOf(handler)
	handlerType := handlerValue.Type()

	if handlerType.Kind() != reflect.Func {
		panic("Handler must be a function")
	}

	return func(c *gin.Context) {
		// Get function parameters and resolve dependencies
		numIn := handlerType.NumIn()
		args := make([]reflect.Value, numIn)

		// First parameter is always gin.Context
		args[0] = reflect.ValueOf(c)

		// Resolve other dependencies
		for i := 1; i < numIn; i++ {
			paramType := handlerType.In(i)
			dependency, err := container.Resolve(paramType)
			if err != nil {
				ErrorResponse(c, 500, "Failed to resolve dependency: "+err.Error())
				return
			}
			args[i] = reflect.ValueOf(dependency)
		}

		// Call the handler
		handlerValue.Call(args)
	}
}
