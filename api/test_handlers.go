package api

import (
	"errors"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// TestRequest 테스트 요청 구조체
type TestRequest struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// TestResponse 테스트 응답 구조체
type TestResponse struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	ResponseTime int64  `json:"response_time_ms"`
	Data         any    `json:"data,omitempty"`
}

// RegisterTestHandlers 테스트 API 핸들러 등록
func RegisterTestHandlers(r *gin.Engine) {
	// 테스트 API 그룹
	testGroup := r.Group("/test")
	{
		// 기본 테스트 API
		testGroup.POST("/echo", EchoHandler)

		// 랜덤 지연 API
		testGroup.POST("/delay", DelayHandler)

		// 랜덤 에러 API
		testGroup.POST("/error", ErrorHandler)

		// 모든 상황 랜덤 API (지연, 에러, 정상)
		testGroup.POST("/random", RandomHandler)
	}
}

// EchoHandler 요청을 그대로 응답하는 핸들러
func EchoHandler(c *gin.Context) {
	var req TestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, TestResponse{
			Success: false,
			Message: "잘못된 요청 형식입니다: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, TestResponse{
		Success:      true,
		Message:      "에코 응답입니다",
		ResponseTime: 0,
		Data:         req,
	})
}

// DelayHandler 랜덤 지연 후 응답하는 핸들러
func DelayHandler(c *gin.Context) {
	var req TestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, TestResponse{
			Success: false,
			Message: "잘못된 요청 형식입니다: " + err.Error(),
		})
		return
	}

	// 0~2초 사이 랜덤 지연
	delayMs := rand.Intn(2000)
	time.Sleep(time.Duration(delayMs) * time.Millisecond)

	c.JSON(http.StatusOK, TestResponse{
		Success:      true,
		Message:      "지연 응답입니다",
		ResponseTime: int64(delayMs),
		Data:         req,
	})
}

// ErrorHandler 랜덤 확률로 에러를 반환하는 핸들러
func ErrorHandler(c *gin.Context) {
	var req TestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, TestResponse{
			Success: false,
			Message: "잘못된 요청 형식입니다: " + err.Error(),
		})
		return
	}

	// 50% 확률로 에러 발생
	if rand.Float32() < 0.5 {
		err := errors.New("랜덤 에러가 발생했습니다")
		c.Set("error", err) // 로거에서 사용할 에러 정보 설정
		c.JSON(http.StatusInternalServerError, TestResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, TestResponse{
		Success: true,
		Message: "정상 응답입니다",
		Data:    req,
	})
}

// RandomHandler 랜덤 지연, 에러, 정상 응답을 모두 포함하는 핸들러
func RandomHandler(c *gin.Context) {
	var req TestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, TestResponse{
			Success: false,
			Message: "잘못된 요청 형식입니다: " + err.Error(),
		})
		return
	}

	// 랜덤 지연 (0~3초)
	delayMs := rand.Intn(3000)
	time.Sleep(time.Duration(delayMs) * time.Millisecond)

	// 랜덤 상태 결정 (30% 확률로 에러)
	if rand.Float32() < 0.3 {
		statusCodes := []int{
			http.StatusBadRequest,
			http.StatusUnauthorized,
			http.StatusForbidden,
			http.StatusNotFound,
			http.StatusInternalServerError,
			http.StatusServiceUnavailable,
		}

		// 랜덤 에러 상태 코드 선택
		statusCode := statusCodes[rand.Intn(len(statusCodes))]
		errorMsg := "랜덤 에러가 발생했습니다 (상태 코드: " + http.StatusText(statusCode) + ")"

		err := errors.New(errorMsg)
		c.Set("error", err)
		c.JSON(statusCode, TestResponse{
			Success:      false,
			Message:      errorMsg,
			ResponseTime: int64(delayMs),
		})
		return
	}

	c.JSON(http.StatusOK, TestResponse{
		Success:      true,
		Message:      "랜덤 테스트 정상 응답입니다",
		ResponseTime: int64(delayMs),
		Data:         req,
	})
}
