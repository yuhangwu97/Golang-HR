package utils

import (
	"reflect"
)

// Module represents a module that can register services and controllers
type Module interface {
	ConfigureServices(container *Container)
	GetControllers() []reflect.Type
	GetServices() []reflect.Type
}

// BaseModule provides basic module functionality
type BaseModule struct {
	controllers []reflect.Type
	services    []reflect.Type
}

func NewBaseModule() *BaseModule {
	return &BaseModule{
		controllers: make([]reflect.Type, 0),
		services:    make([]reflect.Type, 0),
	}
}

func (m *BaseModule) RegisterController(controllerType reflect.Type) {
	m.controllers = append(m.controllers, controllerType)
}

func (m *BaseModule) RegisterService(serviceType reflect.Type) {
	m.services = append(m.services, serviceType)
}

func (m *BaseModule) GetControllers() []reflect.Type {
	return m.controllers
}

func (m *BaseModule) GetServices() []reflect.Type {
	return m.services
}

func (m *BaseModule) ConfigureServices(container *Container) {
	// Default implementation does nothing
	// Override in specific modules to configure services
}

// ModuleBuilder helps build modules with fluent API
type ModuleBuilder struct {
	module *BaseModule
	container *Container
}

func NewModuleBuilder(container *Container) *ModuleBuilder {
	return &ModuleBuilder{
		module:    NewBaseModule(),
		container: container,
	}
}

func (mb *ModuleBuilder) AddController(controllerType reflect.Type, factory interface{}) *ModuleBuilder {
	mb.module.RegisterController(controllerType)
	mb.container.RegisterTransient(controllerType, factory)
	return mb
}

func (mb *ModuleBuilder) AddService(serviceType reflect.Type, factory interface{}) *ModuleBuilder {
	mb.module.RegisterService(serviceType)
	mb.container.RegisterSingleton(serviceType, factory)
	return mb
}

func (mb *ModuleBuilder) Build() Module {
	return mb.module
}

// ApplicationBuilder builds the entire application
type ApplicationBuilder struct {
	container *Container
	modules   []Module
}

func NewApplicationBuilder() *ApplicationBuilder {
	return &ApplicationBuilder{
		container: NewContainer(),
		modules:   make([]Module, 0),
	}
}

func (ab *ApplicationBuilder) AddModule(module Module) *ApplicationBuilder {
	ab.modules = append(ab.modules, module)
	module.ConfigureServices(ab.container)
	return ab
}

func (ab *ApplicationBuilder) GetContainer() *Container {
	return ab.container
}

func (ab *ApplicationBuilder) Build() *Container {
	return ab.container
}