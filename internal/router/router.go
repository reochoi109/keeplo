package router

import (
	"keeplo/internal/handler"

	"github.com/gin-gonic/gin"
)

func Run() {
	r := gin.Default()

	api := r.Group("/api/v1")
	registerUser(api)
	registerSite(api)
	registerPing(api)

	r.Run(":8000")
}

func registerUser(api *gin.RouterGroup) {
	user := api.Group("/user")

	user.POST("/login")      // 로그인
	user.POST("/logout")     // 로그아웃
	user.POST("/validation") // 중복 확인
	user.POST("")            // 회원가입
	user.GET("")             // 내 정보 확인
	user.PUT("")             // 내 정보 수정
	user.DELETE("")          // 탈퇴
}

func registerSite(api *gin.RouterGroup) {
	site := api.Group("/site")

	site.GET("")        // 전체 목록 조회
	site.POST("")       // 추가
	site.GET("/:id")    // 사이트 단건 조회
	site.PUT("/:id")    // 수정
	site.DELETE("/:id") // 삭제

	site.GET("/:id/status") //  특정 사이트 상태 히스토리
	site.GET("/:id/stat")   //  평균 응답시간, 실패 비율 등
	site.POST("/:id/check") //  수동 즉시 체크
}

func registerPing(api *gin.RouterGroup) {
	ping := api.Group("/ping")
	ping.GET("", handler.PingHandler)
}
