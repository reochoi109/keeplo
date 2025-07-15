package router

import (
	"context"
	"keeplo/internal/adapter/repository/monitor_repo"
	"keeplo/internal/adapter/repository/user_repo"
	"keeplo/internal/adapter/rest/handler"
	"keeplo/internal/adapter/rest/middleware"
	"keeplo/internal/application/monitor"
	"keeplo/internal/application/user"
	"keeplo/pkg/db/postgresql"
	"net/http"
	"time"

	_ "keeplo/docs" // 반드시 swag init 후 docs 디렉토리를 import

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Run(ctx context.Context) {
	r := gin.Default()
	api := r.Group("/api/v1")

	// middleware

	// cors

	// https

	// --- TEMP
	userRepo := user_repo.NewGormUserRepo(postgresql.GetDB())
	monitorRepo := monitor_repo.NewGormMonitorRepo(postgresql.GetDB())
	userService := user.NewUserService(userRepo)
	monitorService := monitor.NewMonitorService(monitorRepo, userRepo)
	handlerService := handler.NewHandler(userService, monitorService)
	// --- TEMP

	api.GET("/me", middleware.AuthMiddleware(), func(c *gin.Context) {
		userID, ok := c.Get(middleware.ContextUserIDKey)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user_id": userID})
	})

	registerUserHandler(api, handlerService)
	registerMonitorHandler(api, handlerService)
	registerLogHandler(api)

	srv := &http.Server{
		Addr:              ":8888",
		Handler:           r,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	// swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

func registerUserHandler(api *gin.RouterGroup, handlerService *handler.Handler) {
	auth := api.Group("/auth")

	auth.POST("/signup", handlerService.SignupHandler)                              // 회원가입
	auth.POST("/login", handlerService.LoginHandler)                                // 로그인
	auth.GET("/me", middleware.AuthMiddleware(), handlerService.GetUserInfoHandler) // 로그인 정보 조회
	auth.PUT("/me/nickname", middleware.AuthMiddleware(), handlerService.UpdateNicknameHandler)
	auth.PUT("/me/password", middleware.AuthMiddleware(), handlerService.UpdatePasswordHandler)
	auth.DELETE("/me/logout", middleware.AuthMiddleware(), handlerService.LogoutHandler) // 로그아웃
	auth.DELETE("/me/resign", middleware.AuthMiddleware(), handlerService.ReSignHandler) // 회원 탈퇴 요청
	auth.POST("/password", middleware.AuthMiddleware(), handlerService.CheckPassword)    // 비밀번호 검사
	auth.GET("/duplicate", handlerService.DuplicateEmail)                                // 이메일 중복 검사
}

func registerMonitorHandler(api *gin.RouterGroup, handlerService *handler.Handler) {
	monitor := api.Group("/monitor")

	monitor.GET("/:id", handlerService.GetMonitorHandler)       // 모니터링 주소 정보 조회
	monitor.PUT("/:id", handlerService.UpdateMonitorHandler)    // 모니터링 주소 정보 수정
	monitor.POST("", handlerService.RegisterMonitorHandler)     // 모니터링 주소 추가
	monitor.DELETE("/:id", handlerService.RemoveMonitorHandler) // 모니터링 주소 삭제
	monitor.GET("/list", handlerService.GetMonitorListHandler)  // 모니터링 주소 목록 조회
}

func registerLogHandler(api *gin.RouterGroup) {
	logg := api.Group("/log")

	logg.GET("/health") // 모니터링 헬스 기록 정보 조회
	logg.GET("/status") // 모니터링 상태 기록 정보 조회
}
