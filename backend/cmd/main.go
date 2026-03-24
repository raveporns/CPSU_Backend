package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"cpsu/internal/config"
	"cpsu/internal/connectdb"

	authHandler "cpsu/internal/auth/handler"
	authRepo "cpsu/internal/auth/repository"
	authService "cpsu/internal/auth/service"

	middlewares "cpsu/internal/auth/middlewares"

	userHandler "cpsu/internal/auth/handler"
	userRepo "cpsu/internal/auth/repository"
	userService "cpsu/internal/auth/service"

	auditLogHandler "cpsu/internal/auth/handler"
	auditLogRepo "cpsu/internal/auth/repository"
	auditLogService "cpsu/internal/auth/service"

	newsHandler "cpsu/internal/news/handler"
	newsRepo "cpsu/internal/news/repository"
	newsService "cpsu/internal/news/service"

	courseHandler "cpsu/internal/course/handler"
	courseRepo "cpsu/internal/course/repository"
	courseService "cpsu/internal/course/service"

	structureHandler "cpsu/internal/course_structure/handler"
	structureRepo "cpsu/internal/course_structure/repository"
	structureService "cpsu/internal/course_structure/service"

	roadmapHandler "cpsu/internal/roadmap/handler"
	roadmapRepo "cpsu/internal/roadmap/repository"
	roadmapService "cpsu/internal/roadmap/service"

	subjectHandler "cpsu/internal/subject/handler"
	subjectRepo "cpsu/internal/subject/repository"
	subjectService "cpsu/internal/subject/service"

	personnelHandler "cpsu/internal/personnel/handler"
	personnelRepo "cpsu/internal/personnel/repository"
	"cpsu/internal/personnel/service"
	personnelService "cpsu/internal/personnel/service"

	admissionHandler "cpsu/internal/admission/handler"
	admissionRepo "cpsu/internal/admission/repository"
	admissionService "cpsu/internal/admission/service"

	calendarHandler "cpsu/internal/calendar/handler"
	calendarRepo "cpsu/internal/calendar/repository"
	calendarService "cpsu/internal/calendar/service"
)

func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := connectdb.NewPostgresDatabase(cfg.GetConnectionString())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	authUserRepo := authRepo.NewUserRepository(db.GetDB())
	authRoleRepo := authRepo.NewRoleRepository(db.GetDB())
	authPermissionRepo := authRepo.NewPermissionRepository(db.GetDB())
	authTokenRepo := authRepo.NewTokenRepository(db.GetDB())

	auditLogRepo := auditLogRepo.NewAuditRepository(db.GetDB())
	auditLogService := auditLogService.NewAuditService(auditLogRepo)
	auditLogHandler := auditLogHandler.NewAuditHandler(auditLogService)

	authService := authService.NewAuthService(authUserRepo, authRoleRepo, authTokenRepo, auditLogRepo)
	authHandler := authHandler.NewAuthHandler(authService)
	permissionMiddleware := middlewares.NewPermissionMiddleware(authPermissionRepo)

	roleRepo := userRepo.NewRoleRepository(db.GetDB())
	roleService := userService.NewRoleService(roleRepo, auditLogRepo)
	roleHandler := userHandler.NewRoleHandler(roleService)

	userRepo := userRepo.NewUserRepository(db.GetDB())
	userService := userService.NewUserService(userRepo, auditLogRepo)
	userHandler := userHandler.NewUserHandler(userService)

	newsRepo := newsRepo.NewNewsRepository(db.GetDB())
	newsService := newsService.NewNewsService(newsRepo, auditLogRepo, cfg.MinioEndpoint, cfg.MinioAccessKey, cfg.MinioSecretKey, cfg.MinioBucket, cfg.MinioUseSSL, cfg.MinioPublicBaseURL)
	newsHandler := newsHandler.NewNewsHandler(newsService)

	courseRepo := courseRepo.NewCourseRepository(db.GetDB())
	courseService := courseService.NewCourseService(courseRepo, auditLogRepo)
	courseHandler := courseHandler.NewCourseHandler(courseService)

	structureRepo := structureRepo.NewCourseStructureRepository(db.GetDB())
	structureService := structureService.NewCourseStructureService(structureRepo)
	structureHandler := structureHandler.NewCourseStructureHandler(structureService)

	roadmapRepo := roadmapRepo.NewRoadmapRepository(db.GetDB())
	roadmapService := roadmapService.NewRoadmapService(roadmapRepo, cfg.MinioEndpoint, cfg.MinioAccessKey, cfg.MinioSecretKey, cfg.MinioBucket, cfg.MinioUseSSL, cfg.MinioPublicBaseURL)
	roadmapHandler := roadmapHandler.NewRoadmapHandler(roadmapService)

	subjectRepo := subjectRepo.NewSubjectRepository(db.GetDB())
	subjectService := subjectService.NewSubjectService(subjectRepo, auditLogRepo)
	subjectHandler := subjectHandler.NewSubjectHandler(subjectService)

	personnelRepo := personnelRepo.NewPersonnelRepository(db.GetDB())
	personnelService := personnelService.NewPersonnelService(personnelRepo, auditLogRepo, cfg.MinioEndpoint, cfg.MinioAccessKey, cfg.MinioSecretKey, cfg.MinioBucket, cfg.MinioUseSSL, cfg.MinioPublicBaseURL)
	personnelHandler := personnelHandler.NewPersonnelHandler(personnelService)
	service.SyncScopus(personnelService)

	admissionRepo := admissionRepo.NewAdmissionRepository(db.GetDB())
	admissionService := admissionService.NewAdmissionService(admissionRepo, auditLogRepo, cfg.MinioEndpoint, cfg.MinioAccessKey, cfg.MinioSecretKey, cfg.MinioBucket, cfg.MinioUseSSL, cfg.MinioPublicBaseURL)
	admissionHandler := admissionHandler.NewAdmissionHandler(admissionService)

	calendarRepo := calendarRepo.NewCalendarRepository(db.GetDB())
	calendarService := calendarService.NewCalendarService(calendarRepo, auditLogRepo)
	calendarHandler := calendarHandler.NewCalendarHandler(calendarService)

	go func() {
		for {
			time.Sleep(10 * time.Second)
			if err := db.Ping(); err != nil {
				log.Printf("Database connection lost: %v", err)
				if reconnErr := db.Reconnect(cfg.GetConnectionString()); reconnErr != nil {
					log.Printf("Failed to reconnect: %v", reconnErr)
				} else {
					log.Printf("Successfully reconnected to the database")
				}
			}
		}
	}()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(TimeoutMiddleware(5 * time.Second))

	r.GET("/health", func(c *gin.Context) {
		if err := connectdb.CheckDBConnection(db.GetDB()); err != nil {
			c.JSON(503, gin.H{"detail": "Database connection failed"})
			return
		}
		c.JSON(200, gin.H{"status": "healthy", "database": "connected"})
	})

	public := r.Group("/api/v1")
	{
		public.POST("/auth/login", authHandler.Login)
		public.POST("/auth/refresh", authHandler.RefreshToken)

		public.GET("/news", newsHandler.GetAllNews)
		public.GET("/news/:id", newsHandler.GetNewsByID)

		public.GET("/course", courseHandler.GetAllCourses)
		public.GET("/course/:id", courseHandler.GetCourseByID)

		public.GET("/structure", structureHandler.GetAllCourseStructure)
		public.GET("/structure/:id", structureHandler.GetCourseStructureByID)

		public.GET("/roadmap", roadmapHandler.GetAllRoadmap)
		public.GET("/roadmap/:id", roadmapHandler.GetRoadmapByID)

		public.GET("/subject", subjectHandler.GetAllSubjects)
		public.GET("/subject/:id", subjectHandler.GetSubjectByID)

		public.GET("/personnel", personnelHandler.GetAllPersonnels)
		public.GET("/personnel/:id", personnelHandler.GetPersonnelByID)
		public.GET("/personnel/research", personnelHandler.GetAllResearch)

		public.GET("/admission", admissionHandler.GetAllAdmission)
		public.GET("/admission/:id", admissionHandler.GetAdmissionByID)

		public.GET("/calendar", calendarHandler.GetAllCalendars)
		public.GET("/calendar/:id", calendarHandler.GetCalendarByID)
	}

	protected := r.Group("/api/v1")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.POST("/auth/logout", authHandler.Logout)
	}

	admin := protected.Group("/admin")
	{

		userAdmin := admin.Group("/user")
		{
			userAdmin.GET("", permissionMiddleware.RequirePermission("users:read"), userHandler.GetAllUser)
			userAdmin.POST("", permissionMiddleware.RequirePermission("users:create"), userHandler.CreateUser)
			userAdmin.DELETE("/:id", permissionMiddleware.RequirePermission("users:delete"), userHandler.DeleteUser)
		}

		permissionAdmin := admin.Group("/permission/user")
		{
			permissionAdmin.POST("/:id", permissionMiddleware.RequirePermission("roles:assign"), roleHandler.AssignRole)
		}

		newsAdmin := admin.Group("/news")
		{
			newsAdmin.GET("", permissionMiddleware.RequirePermission("news:read"), newsHandler.GetAllNews)
			newsAdmin.GET("/:id", permissionMiddleware.RequirePermission("news:read_id"), newsHandler.GetNewsByID)
			newsAdmin.POST("", permissionMiddleware.RequirePermission("news:create"), newsHandler.CreateNews)
			newsAdmin.PUT("/:id", permissionMiddleware.RequirePermission("news:update"), newsHandler.UpdateNews)
			newsAdmin.DELETE("/:id", permissionMiddleware.RequirePermission("news:delete"), newsHandler.DeleteNews)
		}

		courseAdmin := admin.Group("/course")
		{
			courseAdmin.GET("", permissionMiddleware.RequirePermission("courses:read"), courseHandler.GetAllCourses)
			courseAdmin.GET("/:id", permissionMiddleware.RequirePermission("courses:read_id"), courseHandler.GetCourseByID)
			courseAdmin.POST("", permissionMiddleware.RequirePermission("courses:create"), courseHandler.CreateCourse)
			courseAdmin.PUT("/:id", permissionMiddleware.RequirePermission("courses:update"), courseHandler.UpdateCourse)
			courseAdmin.DELETE("/:id", permissionMiddleware.RequirePermission("courses:delete"), courseHandler.DeleteCourse)
		}

		structureAdmin := admin.Group("/structure")
		{
			structureAdmin.GET("", permissionMiddleware.RequirePermission("course_structure:read"), structureHandler.GetAllCourseStructure)
			structureAdmin.GET("/:id", permissionMiddleware.RequirePermission("course_structure:read_id"), structureHandler.GetCourseStructureByID)
			structureAdmin.POST("", permissionMiddleware.RequirePermission("course_structure:create"), structureHandler.CreateCourseStructure)
			structureAdmin.PUT("/:id", permissionMiddleware.RequirePermission("course_structure:update"), structureHandler.UpdateCourseStructure)
			structureAdmin.DELETE("/:id", permissionMiddleware.RequirePermission("course_structure:delete"), structureHandler.DeleteCourseStructure)
		}

		roadmapAdmin := admin.Group("/roadmap")
		{
			roadmapAdmin.GET("", permissionMiddleware.RequirePermission("roadmap:read"), roadmapHandler.GetAllRoadmap)
			roadmapAdmin.GET("/:id", permissionMiddleware.RequirePermission("roadmap:read_id"), roadmapHandler.GetRoadmapByID)
			roadmapAdmin.POST("", permissionMiddleware.RequirePermission("roadmap:create"), roadmapHandler.CreateRoadmap)
			roadmapAdmin.DELETE("/:id", permissionMiddleware.RequirePermission("roadmap:delete"), roadmapHandler.DeleteRoadmap)
		}

		subjectAdmin := admin.Group("/subject")
		{
			subjectAdmin.GET("", permissionMiddleware.RequirePermission("subject:read"), subjectHandler.GetAllSubjects)
			subjectAdmin.GET("/:id", permissionMiddleware.RequirePermission("subject:read_id"), subjectHandler.GetSubjectByID)
			subjectAdmin.POST("", permissionMiddleware.RequirePermission("subject:create"), subjectHandler.CreateSubject)
			subjectAdmin.PUT("/:id", permissionMiddleware.RequirePermission("subject:update"), subjectHandler.UpdateSubject)
			subjectAdmin.DELETE("/:id", permissionMiddleware.RequirePermission("subject:delete"), subjectHandler.DeleteSubject)
		}

		personnelAdmin := admin.Group("/personnel")
		{
			personnelAdmin.GET("", permissionMiddleware.RequirePermission("personnel:read"), personnelHandler.GetAllPersonnels)
			personnelAdmin.GET("/:id", permissionMiddleware.RequirePermission("personnel:read_id"), personnelHandler.GetPersonnelByID)
			personnelAdmin.POST("", permissionMiddleware.RequirePermission("personnel:create"), personnelHandler.CreatePersonnel)
			personnelAdmin.PUT("/:id", permissionMiddleware.RequirePermission("personnel:update"), personnelHandler.UpdatePersonnel)
			personnelAdmin.DELETE("/:id", permissionMiddleware.RequirePermission("personnel:delete"), personnelHandler.DeletePersonnel)
			personnelAdmin.GET("/scopus", permissionMiddleware.RequirePermission("scopus:read"), personnelHandler.GetResearchfromScopus)
			personnelAdmin.GET("/research", permissionMiddleware.RequirePermission("research:read"), personnelHandler.GetAllResearch)
		}

		admission := admin.Group("/admission")
		{
			admission.GET("", permissionMiddleware.RequirePermission("admission:read"), admissionHandler.GetAllAdmission)
			admission.GET("/:id", permissionMiddleware.RequirePermission("admission:read_id"), admissionHandler.GetAdmissionByID)
			admission.POST("", permissionMiddleware.RequirePermission("admission:create"), admissionHandler.CreateAdmission)
			admission.PUT("/:id", permissionMiddleware.RequirePermission("admission:update"), admissionHandler.UpdateAdmission)
			admission.DELETE("/:id", permissionMiddleware.RequirePermission("admission:delete"), admissionHandler.DeleteAdmission)
		}

		calendarAdmin := admin.Group("/calendar")
		{
			calendarAdmin.GET("", permissionMiddleware.RequirePermission("calendar:read"), calendarHandler.GetAllCalendars)
			calendarAdmin.GET("/:id", permissionMiddleware.RequirePermission("calendar:read_id"), calendarHandler.GetCalendarByID)
			calendarAdmin.POST("", permissionMiddleware.RequirePermission("calendar:create"), calendarHandler.CreateCalendar)
			calendarAdmin.PUT("/:id", permissionMiddleware.RequirePermission("calendar:update"), calendarHandler.UpdateCalendar)
			calendarAdmin.DELETE("/:id", permissionMiddleware.RequirePermission("calendar:delete"), calendarHandler.DeleteCalendar)
		}

		AuditLogAdmin := admin.Group("/logs")
		{
			AuditLogAdmin.GET("", permissionMiddleware.RequirePermission("logs:read"), auditLogHandler.GetAllAuditLog)
		}
	}

	teacher := protected.Group("/teacher")
	{
		teacherPersonnel := teacher.Group("/personnel")
		{
			teacherPersonnel.PUT("/:id", permissionMiddleware.RequirePermission("your_personnel:update"), personnelHandler.UpdateTeacher)
		}
	}

	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
