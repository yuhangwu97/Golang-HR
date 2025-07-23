package config

import (
	"reflect"

	"gin-project/controllers"
	"gin-project/services"
	"gin-project/utils"

	"gorm.io/gorm"
)

// SetupContainer configures the dependency injection container
func SetupContainer() *utils.Container {
	container := utils.NewContainer()

	// Register repository as singleton (pointer type for dependency injection)
	container.RegisterSingleton(
		reflect.TypeOf((*utils.MySQLRepository)(nil)),
		func() *utils.MySQLRepository {
			if DB == nil || RedisClient == nil {
				panic("Database connections not initialized before container setup")
			}
			return utils.NewMySQLRepository(DB, RedisClient)
		},
	)

	// Register legacy services as singletons
	container.RegisterSingleton(
		reflect.TypeOf((*services.AuthService)(nil)),
		func(repo *utils.MySQLRepository) *services.AuthService {
			service := &services.AuthService{}
			service.InjectDependencies(repo)
			return service
		},
	)

	container.RegisterSingleton(
		reflect.TypeOf((*services.DepartmentServiceInterface)(nil)).Elem(),
		func(db *gorm.DB) services.DepartmentServiceInterface {
			service := &services.DepartmentService{}
			service.InjectDependencies(db)
			return service
		},
	)

	container.RegisterSingleton(
		reflect.TypeOf((*services.EmployeeServiceInterface)(nil)).Elem(),
		func(db *gorm.DB) services.EmployeeServiceInterface {
			service := &services.EmployeeService{}
			service.InjectDependencies(db)
			return service
		},
	)

	container.RegisterSingleton(
		reflect.TypeOf((*services.SalaryServiceInterface)(nil)).Elem(),
		func(db *gorm.DB) services.SalaryServiceInterface {
			service := &services.SalaryService{}
			service.InjectDependencies(db)
			return service
		},
	)

	container.RegisterSingleton(
		reflect.TypeOf((*services.AttendanceServiceInterface)(nil)).Elem(),
		func(db *gorm.DB) services.AttendanceServiceInterface {
			service := &services.AttendanceService{}
			service.InjectDependencies(db)
			return service
		},
	)

	container.RegisterSingleton(
		reflect.TypeOf((*services.PositionServiceInterface)(nil)).Elem(),
		func(db *gorm.DB) services.PositionServiceInterface {
			service := &services.PositionService{}
			service.InjectDependencies(db)
			return service
		},
	)

	container.RegisterSingleton(
		reflect.TypeOf((*services.JobLevelServiceInterface)(nil)).Elem(),
		func(db *gorm.DB) services.JobLevelServiceInterface {
			return services.NewJobLevelService(db)
		},
	)

	container.RegisterSingleton(
		reflect.TypeOf((*services.UserServiceInterface)(nil)).Elem(),
		func(db *gorm.DB) services.UserServiceInterface {
			return services.NewUserService(db)
		},
	)

	container.RegisterSingleton(
		reflect.TypeOf((*services.RoleServiceInterface)(nil)).Elem(),
		func(db *gorm.DB) services.RoleServiceInterface {
			return services.NewRoleService(db)
		},
	)

	container.RegisterSingleton(
		reflect.TypeOf((*services.PermissionServiceInterface)(nil)).Elem(),
		func(db *gorm.DB) services.PermissionServiceInterface {
			return services.NewPermissionService(db)
		},
	)

	container.RegisterSingleton(
		reflect.TypeOf((*services.WebSocketService)(nil)),
		func() *services.WebSocketService {
			return services.NewWebSocketService()
		},
	)

	container.RegisterSingleton(
		reflect.TypeOf((*services.NotificationService)(nil)),
		func(wsService *services.WebSocketService) *services.NotificationService {
			return services.NewNotificationService(wsService)
		},
	)

	container.RegisterSingleton(
		reflect.TypeOf((*services.OrganizationServiceInterface)(nil)).Elem(),
		func(db *gorm.DB, deptService services.DepartmentServiceInterface, empService services.EmployeeServiceInterface, posService services.PositionServiceInterface, jobService services.JobLevelServiceInterface) services.OrganizationServiceInterface {
			return services.NewOrganizationService(db, deptService, empService, posService, jobService)
		},
	)

	// Register database as singleton
	container.RegisterSingleton(
		reflect.TypeOf((*gorm.DB)(nil)),
		func() *gorm.DB {
			return DB
		},
	)

	// Register controllers as transient
	container.RegisterTransient(
		reflect.TypeOf((*controllers.UserController)(nil)),
		func(userService services.UserServiceInterface) *controllers.UserController {
			controller := &controllers.UserController{}
			// Note: This assumes UserController has been updated to use UserServiceInterface
			// If not, this will need to be adjusted based on actual implementation
			return controller
		},
	)

	container.RegisterTransient(
		reflect.TypeOf((*controllers.AuthController)(nil)),
		func(authService *services.AuthService) *controllers.AuthController {
			controller := &controllers.AuthController{}
			controller.InjectDependencies(authService)
			return controller
		},
	)

	container.RegisterTransient(
		reflect.TypeOf((*controllers.DepartmentController)(nil)),
		func(departmentService services.DepartmentServiceInterface) *controllers.DepartmentController {
			return controllers.NewDepartmentController(departmentService)
		},
	)

	container.RegisterTransient(
		reflect.TypeOf((*controllers.EmployeeController)(nil)),
		func(employeeService services.EmployeeServiceInterface) *controllers.EmployeeController {
			return controllers.NewEmployeeController(employeeService)
		},
	)

	container.RegisterTransient(
		reflect.TypeOf((*controllers.SalaryController)(nil)),
		func(salaryService services.SalaryServiceInterface) *controllers.SalaryController {
			return controllers.NewSalaryController(salaryService)
		},
	)

	container.RegisterTransient(
		reflect.TypeOf((*controllers.AttendanceController)(nil)),
		func(attendanceService services.AttendanceServiceInterface) *controllers.AttendanceController {
			return controllers.NewAttendanceController(attendanceService)
		},
	)

	container.RegisterTransient(
		reflect.TypeOf((*controllers.PositionController)(nil)),
		func(positionService services.PositionServiceInterface) *controllers.PositionController {
			return controllers.NewPositionController(positionService)
		},
	)

	container.RegisterTransient(
		reflect.TypeOf((*controllers.JobLevelController)(nil)),
		func(jobLevelService services.JobLevelServiceInterface) *controllers.JobLevelController {
			return controllers.NewJobLevelController(jobLevelService)
		},
	)

	container.RegisterTransient(
		reflect.TypeOf((*controllers.SystemUserController)(nil)),
		func(userService services.UserServiceInterface) *controllers.SystemUserController {
			return controllers.NewSystemUserController(userService)
		},
	)

	container.RegisterTransient(
		reflect.TypeOf((*controllers.RoleController)(nil)),
		func(roleService services.RoleServiceInterface, permissionService services.PermissionServiceInterface) *controllers.RoleController {
			return controllers.NewRoleController(roleService, permissionService)
		},
	)

	container.RegisterTransient(
		reflect.TypeOf((*controllers.PermissionController)(nil)),
		func(permissionService services.PermissionServiceInterface) *controllers.PermissionController {
			return controllers.NewPermissionController(permissionService)
		},
	)

	container.RegisterTransient(
		reflect.TypeOf((*controllers.WebSocketController)(nil)),
		func(wsService *services.WebSocketService) *controllers.WebSocketController {
			return controllers.NewWebSocketController(wsService)
		},
	)

	container.RegisterTransient(
		reflect.TypeOf((*controllers.OrganizationController)(nil)),
		func(organizationService services.OrganizationServiceInterface) *controllers.OrganizationController {
			return controllers.NewOrganizationController(organizationService)
		},
	)

	return container
}

// Global container instance
var Container *utils.Container

// InitContainer initializes the container after database connections are established
func InitContainer() {
	Container = SetupContainer()
}

// GetContainer returns the global container instance
func GetContainer() *ServiceContainer {
	return &ServiceContainer{Container}
}

// ServiceContainer provides easy access to services
type ServiceContainer struct {
	container *utils.Container
}

func (sc *ServiceContainer) UserService() services.UserServiceInterface {
	service, _ := sc.container.Resolve(reflect.TypeOf((*services.UserServiceInterface)(nil)).Elem())
	return service.(services.UserServiceInterface)
}

func (sc *ServiceContainer) RoleService() services.RoleServiceInterface {
	service, _ := sc.container.Resolve(reflect.TypeOf((*services.RoleServiceInterface)(nil)).Elem())
	return service.(services.RoleServiceInterface)
}

func (sc *ServiceContainer) PermissionService() services.PermissionServiceInterface {
	service, _ := sc.container.Resolve(reflect.TypeOf((*services.PermissionServiceInterface)(nil)).Elem())
	return service.(services.PermissionServiceInterface)
}

func (sc *ServiceContainer) SalaryService() services.SalaryServiceInterface {
	service, _ := sc.container.Resolve(reflect.TypeOf((*services.SalaryServiceInterface)(nil)).Elem())
	return service.(services.SalaryServiceInterface)
}
