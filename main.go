package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Gin 라우터 생성
	r := gin.Default()

	// 헬스 체크 API
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	// 루트 경로 핸들러
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "DevOps 인프라 테스트용 서버가 실행 중입니다.",
		})
	})

	// 포트 설정 (환경 변수에서 가져오거나 기본값 사용)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 서버 시작
	log.Printf("서버가 :%s 포트에서 시작됩니다...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("서버 시작 실패: %v", err)
	}
}
