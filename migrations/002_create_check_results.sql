-- 헬스체크 결과 테이블
-- 각 모니터링 대상에 대한 체크 결과를 append-only로 기록한다.
CREATE TABLE check_results (
    id               BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,  -- 대량 삽입에 UUID보다 효율적
    target_id        UUID NOT NULL REFERENCES targets(id),             -- 대상 FK (soft delete 대상도 참조 유지)
    status_code      INTEGER,                                          -- HTTP 상태 코드 (타임아웃/DNS 실패 시 NULL)
    response_time_ms INTEGER NOT NULL,                                 -- 응답 소요 시간 (ms)
    is_healthy       BOOLEAN NOT NULL,                                 -- 체크 시점 건강 판정 결과
    error_message    TEXT,                                              -- 실패 사유 (성공 시 NULL)
    response_body    TEXT,                                              -- 응답 본문 (디버깅 용도, 선택적)
    response_headers JSONB,                                             -- 응답 헤더 (키-값 구조)
    checked_at       TIMESTAMPTZ NOT NULL DEFAULT now()                -- 체크 수행 시각
);

-- 특정 대상의 최근 결과 조회용 (GET /api/status)
CREATE INDEX idx_check_results_target_checked ON check_results (target_id, checked_at DESC);
-- 전체 결과 시간순 조회용
CREATE INDEX idx_check_results_checked_at ON check_results (checked_at DESC);
