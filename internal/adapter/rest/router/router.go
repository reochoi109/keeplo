package router

import (
	"context"
	"keeplo/internal/adapter/rest/handler"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}

	// https server
	// srv.ListenAndServeTLS(certFile string, keyFile string)

}

func registerUserHandler(api *gin.RouterGroup) {
	auth := api.Group("/auth")

	auth.POST("/signup", handler.SingupHandler)   // 회원가입
	auth.POST("/login", handler.LoginHandler)     // 로그인
	auth.GET("/me", handler.GetUserHandler)       // 로그인 정보 조회
	auth.DELETE("/logout", handler.LogoutHandler) // 로그아웃
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
