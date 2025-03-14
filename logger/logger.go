package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LogLevel 로그 레벨 정의
type LogLevel string

const (
	INFO  LogLevel = "INFO"
	ERROR LogLevel = "ERROR"
	WARN  LogLevel = "WARN"
	DEBUG LogLevel = "DEBUG"
)

// LogEntry 로그 구조체 정의
type LogEntry struct {
	Timestamp        string      `json:"timestamp"`
	Level            LogLevel    `json:"level"`
	RequestID        string      `json:"request_id"`
	ServiceName      string      `json:"service_name"`
	URL              string      `json:"url"`
	ResponseTimeMil  int64       `json:"response_time_mil"`
	HTTPMethod       string      `json:"http_method"`
	HTTPStatusCode   int         `json:"http_status_code"`
	HTTPRequestBody  interface{} `json:"http_request_body,omitempty"`
	HTTPResponseBody interface{} `json:"http_response_body,omitempty"`
	ErrorMessage     string      `json:"error_message,omitempty"`
	StackTrace       string      `json:"stack_trace,omitempty"`
}

// 서비스 이름 환경 변수
var serviceName = os.Getenv("SERVICE_NAME")

func init() {
	// 서비스 이름이 설정되지 않았으면 기본값 사용
	if serviceName == "" {
		serviceName = "devops-prac-server"
	}
}

// LoggerMiddleware Gin 미들웨어 - 요청/응답 로깅
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 요청 시작 시간
		startTime := time.Now()

		// 요청 ID 생성
		requestID := uuid.New().String()
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		// 요청 바디 캡처
		var requestBodyBytes []byte
		if c.Request.Body != nil && c.Request.ContentLength > 0 {
			requestBodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBodyBytes))
		}

		// 응답 바디 캡처를 위한 writer
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// 다음 핸들러 실행
		c.Next()

		// 응답 시간 계산
		duration := time.Since(startTime).Milliseconds()

		// 로그 레벨 결정
		level := INFO
		if c.Writer.Status() >= 400 {
			level = ERROR
		}

		// 요청 바디 파싱
		var requestBody interface{}
		if len(requestBodyBytes) > 0 {
			json.Unmarshal(requestBodyBytes, &requestBody)
		}

		// 응답 바디 파싱
		var responseBody interface{}
		if len(blw.body.Bytes()) > 0 {
			json.Unmarshal(blw.body.Bytes(), &responseBody)
		}

		// 에러 정보 수집
		var errorMessage, stackTrace string
		if err, exists := c.Get("error"); exists {
			if err, ok := err.(error); ok {
				errorMessage = err.Error()
			}
		}

		// 로그 엔트리 생성
		logEntry := LogEntry{
			Timestamp:        time.Now().Format(time.RFC3339),
			Level:            level,
			RequestID:        requestID,
			ServiceName:      serviceName,
			URL:              c.Request.URL.String(),
			ResponseTimeMil:  duration,
			HTTPMethod:       c.Request.Method,
			HTTPStatusCode:   c.Writer.Status(),
			HTTPRequestBody:  requestBody,
			HTTPResponseBody: responseBody,
			ErrorMessage:     errorMessage,
			StackTrace:       stackTrace,
		}

		// JSON으로 변환하여 출력
		logJSON, _ := json.Marshal(logEntry)
		fmt.Println(string(logJSON))
	}
}

// bodyLogWriter 응답 바디를 캡처하기 위한 구조체
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 응답 바디 캡처
func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
