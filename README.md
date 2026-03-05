# scorpes

Go 1.26 기반 HTTP API 서버.

## 요구사항

- Go 1.26+
- Docker & Docker Compose (컨테이너 실행 시)
- [air](https://github.com/air-verse/air) (핫리로드 사용 시)

## 시작하기

### 로컬 실행

```bash
# 직접 실행
make run

# 핫리로드 (air 설치 필요)
go install github.com/air-verse/air@latest
make dev
```

### Docker

```bash
# 개발 환경 (air 핫리로드)
make docker-dev

# 프로덕션 빌드
make docker-prod
```

## API

| Method | Path | 설명 |
|--------|------|------|
| GET | `/health` | 헬스 체크 |

```bash
curl http://localhost:8080/health
# {"message":"Health check successful"}
```

## 빌드

```bash
# 바이너리 빌드 (./tmp/main)
make build

# 테스트
make test

# 빌드 아티팩트 삭제
make clean
```
