#!/bin/bash

# 로그 테스트를 위한 부하 테스트 스크립트
# 사용법: ./load_test.sh [동시성] [시간(초)]

# 기본값 설정
CONCURRENCY=${1:-10}
DURATION=${2:-60}
BASE_URL="http://localhost:8080"

# 필요한 도구 확인
if ! command -v siege &> /dev/null; then
    echo "siege가 설치되어 있지 않습니다. 설치하려면 다음 명령어를 실행하세요:"
    echo "  brew install siege (macOS)"
    echo "  apt-get install siege (Ubuntu/Debian)"
    echo "  yum install siege (CentOS/RHEL)"
    exit 1
fi

# 테스트 URL 파일 생성
cat > urls.txt << EOF
${BASE_URL}/health
${BASE_URL}/
${BASE_URL}/test/echo POST {"message": "에코 테스트", "data": {"timestamp": "$(date +%s)"}}
${BASE_URL}/test/delay POST {"message": "지연 테스트", "data": {"timestamp": "$(date +%s)"}}
${BASE_URL}/test/error POST {"message": "에러 테스트", "data": {"timestamp": "$(date +%s)"}}
${BASE_URL}/test/random POST {"message": "랜덤 테스트", "data": {"timestamp": "$(date +%s)"}}
EOF

echo "부하 테스트를 시작합니다..."
echo "동시성: $CONCURRENCY, 시간: ${DURATION}초"
echo "테스트 URL: $BASE_URL"

# siege 실행
siege -c $CONCURRENCY -t ${DURATION}S -f urls.txt -H "Content-Type: application/json"

# 임시 파일 정리
rm urls.txt

echo "부하 테스트가 완료되었습니다." 