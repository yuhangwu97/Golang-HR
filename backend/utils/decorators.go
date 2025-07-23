package utils

import (
	"reflect"
)

// Injectable marks a struct as injectable
type Injectable struct{}

// Service decorator for service classes
type Service struct {
	Scope ServiceScope
}

// Controller decorator for controller classes
type Controller struct{}

// ServiceMetadata stores metadata for services
type ServiceMetadata struct {
	Type         reflect.Type
	Dependencies []reflect.Type
	Scope        ServiceScope
}

var serviceRegistry = make(map[reflect.Type]*ServiceMetadata)

// RegisterService registers a service with its dependencies
func RegisterService(serviceType interface{}, dependencies []interface{}, scope ServiceScope) {
	sType := reflect.TypeOf(serviceType)
	if sType.Kind() == reflect.Ptr {
		sType = sType.Elem()
	}

	deps := make([]reflect.Type, len(dependencies))
	for i, dep := range dependencies {
		depType := reflect.TypeOf(dep)
		if depType.Kind() == reflect.Ptr {
			depType = depType.Elem()
		}
		deps[i] = depType
	}

	serviceRegistry[sType] = &ServiceMetadata{
		Type:         sType,
		Dependencies: deps,
		Scope:        scope,
	}
}

// GetServiceMetadata gets metadata for a service
func GetServiceMetadata(serviceType reflect.Type) *ServiceMetadata {
	if serviceType.Kind() == reflect.Ptr {
		serviceType = serviceType.Elem()
	}
	return serviceRegistry[serviceType]
}

// AutoRegisterServices automatically registers services in the container
func AutoRegisterServices(container *Container) {
	for serviceType, metadata := range serviceRegistry {
		if metadata.Scope == Singleton {
			container.RegisterSingleton(serviceType, createFactory(serviceType, metadata.Dependencies))
		} else {
			container.RegisterTransient(serviceType, createFactory(serviceType, metadata.Dependencies))
		}
	}
}

// createFactory creates a factory function for a service
func createFactory(serviceType reflect.Type, dependencies []reflect.Type) interface{} {
	// Create a function type that takes dependencies and returns the service
	var inTypes []reflect.Type
	for _, dep := range dependencies {
		inTypes = append(inTypes, reflect.PtrTo(dep))
	}
	
	outTypes := []reflect.Type{reflect.PtrTo(serviceType), reflect.TypeOf((*error)(nil)).Elem()}
	funcType := reflect.FuncOf(inTypes, outTypes, false)
	
	factory := reflect.MakeFunc(funcType, func(args []reflect.Value) []reflect.Value {
		// Create new instance
		serviceValue := reflect.New(serviceType)
		
		// Inject dependencies if the service has an injection method
		if injector, ok := serviceValue.Interface().(DependencyInjector); ok {
			deps := make([]interface{}, len(args))
			for i, arg := range args {
				deps[i] = arg.Interface()
			}
			if err := injector.InjectDependencies(deps...); err != nil {
				return []reflect.Value{
					reflect.Zero(reflect.PtrTo(serviceType)),
					reflect.ValueOf(err),
				}
			}
		}
		
		return []reflect.Value{
			serviceValue,
			reflect.Zero(reflect.TypeOf((*error)(nil)).Elem()),
		}
	})
	
	return factory.Interface()
}

// DependencyInjector interface for services that need dependency injection
type DependencyInjector interface {
	InjectDependencies(deps ...interface{}) error
}