-- 모니터링 대상 테이블
-- 헬스체크를 수행할 URL 및 설정 정보를 관리한다.
-- soft delete 방식으로 삭제하여 check_results 이력을 보존한다.
CREATE TABLE targets (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),  -- API 노출 시 sequential ID 추측 방지
    name             TEXT NOT NULL,                               -- 대상 이름 (표시용)
    url              TEXT NOT NULL,                               -- 헬스체크 호출 URL
    method           TEXT NOT NULL DEFAULT 'GET',                 -- HTTP 메서드 (GET, HEAD, POST 등)
    interval_seconds INTEGER NOT NULL DEFAULT 60,                 -- 체크 주기 (초)
    timeout_seconds  INTEGER NOT NULL DEFAULT 10,                 -- 요청 타임아웃 (초)
    is_active        BOOLEAN NOT NULL DEFAULT TRUE,               -- 모니터링 활성 여부 (삭제 없이 일시 중지 가능)
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at       TIMESTAMPTZ                                  -- soft delete 시각 (NULL이면 미삭제)
);

-- soft delete 필터링용
CREATE INDEX idx_targets_deleted_at ON targets (deleted_at);
-- 활성 대상만 빠르게 조회 (스케줄러에서 사용)
CREATE INDEX idx_targets_is_active ON targets (is_active) WHERE deleted_at IS NULL;
