# 빌드 스테이지
FROM golang:1.23.3-alpine3.20 AS builder

WORKDIR /app

# 의존성 파일 복사 및 다운로드
COPY go.mod go.sum ./
RUN go mod download

# 소스 코드 복사
COPY . .

# 애플리케이션 빌드
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# 실행 스테이지
FROM alpine:3.19

WORKDIR /app

# 타임존 설정
RUN apk --no-cache add tzdata
ENV TZ=Asia/Seoul

# 빌드 스테이지에서 빌드된 바이너리 복사
COPY --from=builder /app/server .

# 포트 노출
EXPOSE 8080

# 실행 명령
CMD ["./server"] 