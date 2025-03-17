package main

import (
	"log"
	"net/http"
	"os"

	"github.com/baekjunho/devops-prac-server/api"
	"github.com/baekjunho/devops-prac-server/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	// Gin 라우터 생성
	r := gin.New() // Default() 대신 New()를 사용하여 기본 미들웨어 제외

	// 로깅 미들웨어 등록
	r.Use(logger.LoggerMiddleware())
	r.Use(gin.Recovery()) // panic 복구 미들웨어

	// 헬스 체크 API
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	// 루트 경로 핸들러
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "DevOps 인프라 테스트용 서버가 실행 중입니다!",
		})
	})

	// 테스트 API 핸들러 등록
	api.RegisterTestHandlers(r)

	// 포트 설정 (환경 변수에서 가져오거나 기본값 사용)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("서버 시작 실패: %v", err)
	}
}
