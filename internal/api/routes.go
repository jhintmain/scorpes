package api

import "net/http"

func RegisterRoutes(r *Router) {
	r.Use(LoggerMiddleware)

	r.GET("/health", healthCheckHandler)

	// 모니터링 대상 목록 조회
	r.GET("/api/targets", healthCheckHandler)
	// 신규 대상 추가
	r.POST("/api/targets", healthCheckHandler)
	// 대상 삭제
	r.DELETE("/api/targets/:id", healthCheckHandler)
	// Uptime 상태 요약 조회
	r.GET("/api/status", healthCheckHandler)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    "OK",
	})
}
