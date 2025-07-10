package router

import (
	"context"
	"keeplo/internal/adapter/rest/handler"
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

	registerUserHandler(api)
	registerMonitorHandler(api)
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

	// https server
	// srv.ListenAndServeTLS(certFile string, keyFile string)

}

func registerUserHandler(api *gin.RouterGroup) {
	auth := api.Group("/auth")

	auth.POST("/signup", handler.SignupHandler)          // 회원가입
	auth.POST("/login", handler.LoginHandler)            // 로그인
	auth.GET("/me/:id", handler.GetUserInfoHandler)      // 로그인 정보 조회
	auth.PUT("me/:id", handler.UpdateUserInfoHandler)    // 사용자 정보 수정
	auth.DELETE("/me/:id/logout", handler.LogoutHandler) // 로그아웃
	auth.DELETE("/me/:id/resign", handler.ReSignHandler) // 회원 탈퇴 요청
	auth.GET("/duplicate", handler.DuplicateEmail)       // 이메일 중복 검사
}

func registerMonitorHandler(api *gin.RouterGroup) {
	monitor := api.Group("/monitor")

	monitor.GET("", handler.GetMonitorHandler)          // 모니터링 주소 정보 조회
	monitor.PUT("", handler.UpdateMonitorHandler)       // 모니터링 주소 정보 수정
	monitor.POST("", handler.RegisterMonitorHandler)    // 모니터링 주소 추가
	monitor.DELETE("", handler.RemoveMonitorHandler)    // 모니터링 주소 삭제
	monitor.GET("/list", handler.GetMonitorListHandler) // 모니터링 주소 목록 조회
}

func registerLogHandler(api *gin.RouterGroup) {
	logg := api.Group("/log")

	logg.GET("/health") // 모니터링 헬스 기록 정보 조회
	logg.GET("/status") // 모니터링 상태 기록 정보 조회
}
