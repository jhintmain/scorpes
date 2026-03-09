-- name: InsertCheckResult :one
INSERT INTO check_results (target_id, status_code, response_time_ms, is_healthy, error_message, response_body, response_headers)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetUptimeSummary :many
SELECT
    t.id,
    t.name,
    t.url,
    COUNT(cr.id) AS total_checks,
    COUNT(cr.id) FILTER (WHERE cr.is_healthy = true) AS healthy_checks,
    CASE
        WHEN COUNT(cr.id) > 0
        THEN ROUND((COUNT(cr.id) FILTER (WHERE cr.is_healthy = true))::numeric / COUNT(cr.id) * 100, 2)
        ELSE 0
    END AS uptime_percentage,
    MAX(cr.checked_at) AS last_checked_at
FROM targets t
LEFT JOIN check_results cr ON cr.target_id = t.id AND cr.checked_at >= now() - make_interval(hours => @hours::int)
WHERE t.deleted_at IS NULL AND t.is_active = true
GROUP BY t.id, t.name, t.url
ORDER BY t.name;
