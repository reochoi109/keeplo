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
	// cors

	// middleware
	r.Use(middleware.UseTraceID())

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
	registerLogHandler(api, handlerService)

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

	auth.POST("/signup", handlerService.SignupHandler)                                          // 회원가입
	auth.POST("/login", handlerService.LoginHandler)                                            // 로그인
	auth.GET("/me", middleware.AuthMiddleware(), handlerService.GetUserInfoHandler)             // 로그인 정보 조회
	auth.PUT("/me/nickname", middleware.AuthMiddleware(), handlerService.UpdateNicknameHandler) //
	auth.PUT("/me/password", middleware.AuthMiddleware(), handlerService.UpdatePasswordHandler) //
	auth.DELETE("/me/logout", middleware.AuthMiddleware(), handlerService.LogoutHandler)        // 로그아웃
	auth.DELETE("/me/resign", middleware.AuthMiddleware(), handlerService.ReSignHandler)        // 회원 탈퇴 요청
	auth.POST("/password", middleware.AuthMiddleware(), handlerService.CheckPassword)           // 비밀번호 검사
	auth.GET("/duplicate", handlerService.DuplicateEmail)                                       // 이메일 중복 검사

	// 보안 기능 (추후 이메일 연동 시 구현)
	// auth.POST("/forgot-password", handlerService.ForgotPasswordHandler) // 비밀번호 초기화 요청
}

func registerMonitorHandler(api *gin.RouterGroup, handlerService *handler.Handler) {
	monitor := api.Group("/monitor")

	monitor.POST("", handlerService.RegisterMonitorHandler)     // 모니터링 주소 추가
	monitor.GET("/list", handlerService.GetMonitorListHandler)  // 모니터 목록 조회
	monitor.GET("/:id", handlerService.GetMonitorHandler)       // 단일 모니터 조회
	monitor.PUT("/:id", handlerService.UpdateMonitorHandler)    // 모니터 수정
	monitor.DELETE("/:id", handlerService.RemoveMonitorHandler) // 모니터 삭제

	// 추가된 기능들
	monitor.PATCH("/:id/toggle", handlerService.ToggleMonitorHandler)      // 모니터 ON/OFF
	monitor.POST("/:id/trigger", handlerService.TriggerMonitorHandler)     // 수동 검사 요청
	monitor.GET("/protocols", handlerService.GetSupportedProtocolsHandler) // 지원 프로토콜 목록
}

func registerLogHandler(api *gin.RouterGroup, handlerService *handler.Handler) {
	logg := api.Group("/log", middleware.AuthMiddleware())

	// [1] 모니터 헬스 로그 전체 조회 (리스트, 필터 포함 가능)
	logg.GET("/health/:monitor_id", handlerService.GetMonitorHealthLogHandler)

	// [2] 실패 요약 (예: 에러 빈도, 마지막 실패 등)
	logg.GET("/health/:monitor_id/errors", handlerService.GetHealthErrorSummaryHandler)

	// [3] 응답 시간 시계열 데이터 (차트용)
	logg.GET("/health/:monitor_id/timeseries", handlerService.GetResponseTimeChartHandler)

	// [4] 차트
	logg.GET("/health/:monitor_id/timeseries", handlerService.GetResponseTimeChartHandler)

	// [5] 알림 전송 이력 (예: 이메일/슬랙 등)
	// logg.GET("/notifications/:monitor_id", handlerService.GetNotificationLogsHandler)
}
