# DevOps 인프라 테스트용 서버

이 프로젝트는 DevOps 인프라 테스트를 위한 간단한 Golang 서버입니다. Gin 프레임워크를 사용하여 구현되었으며, 헬스 체크 API를 제공합니다.

## 기능

- 헬스 체크 API (`/health`)
- 루트 경로 메시지 (`/`)

## 로컬에서 실행하기

### 요구사항

- Go 1.16 이상
- Docker (선택 사항)

### 직접 실행

```bash
# 의존성 설치
go mod download

# 서버 실행
go run main.go
```

기본적으로 서버는 8080 포트에서 실행됩니다. 환경 변수 `PORT`를 설정하여 포트를 변경할 수 있습니다.

### Docker로 실행

```bash
# 이미지 빌드
docker build -t devops-prac-server .

# 컨테이너 실행
docker run -p 8080:8080 devops-prac-server
```

## API 엔드포인트

- `GET /health`: 서버 상태 확인

  - 응답: `{"status": "healthy"}`

- `GET /`: 루트 경로
  - 응답: `{"message": "DevOps 인프라 테스트용 서버가 실행 중입니다."}`

## Kubernetes에 배포하기

이 서버는 Kubernetes 환경에서 실행하기 위해 설계되었습니다. 다음은 기본적인 배포 예시입니다:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: devops-prac-server
spec:
  replicas: 3
  selector:
    matchLabels:
      app: devops-prac-server
  template:
    metadata:
      labels:
        app: devops-prac-server
    spec:
      containers:
      - name: devops-prac-server
        image: [DOCKERHUB_USERNAME]/devops-prac-server:latest
        ports:
        - containerPort: 8080
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
```

## CI/CD

이 프로젝트는 GitHub Actions를 사용하여 Docker 이미지를 자동으로 빌드하고 DockerHub에 푸시합니다. 워크플로우는 `.github/workflows/docker-publish.yml`에 정의되어 있습니다.

GitHub Actions를 사용하기 위해서는 다음 시크릿을 설정해야 합니다:

- `DOCKERHUB_TOKEN`: DockerHub 인증을 위한 액세스 토큰

## 로깅 시스템 테스트

이 프로젝트는 로깅 시스템을 테스트하기 위한 API를 제공합니다. 다양한 상황에서 로그가 어떻게 생성되는지 테스트할 수 있습니다.

### 로그 구조

로그는 다음과 같은 구조로 생성됩니다:

```json
{
  "timestamp": "2023-06-01T12:34:56Z",
  "level": "INFO",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "service_name": "devops-prac-server",
  "url": "/test/echo",
  "response_time_mil": 123,
  "http_method": "POST",
  "http_status_code": 200,
  "http_request_body": { "message": "테스트" },
  "http_response_body": { "success": true, "message": "에코 응답입니다" },
  "error_message": "",
  "stack_trace": ""
}
```

### 테스트 API 엔드포인트

다음 API 엔드포인트를 사용하여 다양한 로그 상황을 테스트할 수 있습니다:

1. **에코 API**: 요청을 그대로 응답합니다.

   ```
   POST /test/echo
   ```

2. **지연 API**: 0~2초 사이의 랜덤한 지연 후 응답합니다.

   ```
   POST /test/delay
   ```

3. **에러 API**: 50% 확률로 에러를 반환합니다.

   ```
   POST /test/error
   ```

4. **랜덤 API**: 랜덤 지연, 에러, 정상 응답을 모두 포함합니다.
   ```
   POST /test/random
   ```

### 테스트 방법

#### 개별 API 테스트

curl을 사용하여 개별 API를 테스트할 수 있습니다:

```bash
# 에코 API 테스트
curl -X POST http://localhost:8080/test/echo -H "Content-Type: application/json" -d '{"message": "테스트 메시지"}'

# 지연 API 테스트
curl -X POST http://localhost:8080/test/delay -H "Content-Type: application/json" -d '{"message": "지연 테스트"}'

# 에러 API 테스트
curl -X POST http://localhost:8080/test/error -H "Content-Type: application/json" -d '{"message": "에러 테스트"}'

# 랜덤 API 테스트
curl -X POST http://localhost:8080/test/random -H "Content-Type: application/json" -d '{"message": "랜덤 테스트"}'
```

#### 부하 테스트

부하 테스트를 위해 siege 도구를 사용할 수 있습니다. 프로젝트에 포함된 스크립트를 사용하여 간편하게 테스트할 수 있습니다:

```bash
# siege 설치 (macOS)
brew install siege

# 부하 테스트 실행 (기본값: 동시성 10, 시간 60초)
./test/load_test.sh

# 사용자 정의 동시성 및 시간으로 테스트
./test/load_test.sh 20 120  # 동시성 20, 시간 120초
```

### 로그 확인

서버 실행 시 표준 출력으로 로그가 출력됩니다. 이 로그를 파일로 저장하거나 로그 수집 시스템으로 전송할 수 있습니다:

```bash
# 로그를 파일로 저장
go run main.go > server.log

# 로그를 파일로 저장하고 터미널에도 출력
go run main.go | tee server.log
```
