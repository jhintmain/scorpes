package api

import "net/http"

func RegisterRoutes(r *Router, h *TargetHandler) {
	r.Use(LoggerMiddleware)

	r.GET("/health", healthCheckHandler)

	// Target API
	r.GET("/api/targets", h.ListTargets)
	r.POST("/api/targets", h.CreateTarget)
	r.PUT("/api/targets/{id}", h.UpdateTarget)
	r.DELETE("/api/targets/{id}", h.DeleteTarget)

	// Uptime 상태 요약 조회
	r.GET("/api/status", healthCheckHandler)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    "OK",
	})
}
