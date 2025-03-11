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
