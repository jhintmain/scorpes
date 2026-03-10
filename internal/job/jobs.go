package job

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hooneun/scorpes/internal/util"
)

var (
	client = &http.Client{
		Timeout: 3 * time.Second,
	}
)

// HealthCheck [job] : health check
func HealthCheck() {
	resp, err := client.Get(util.ApiURL("/health"))
	if err != nil {
		fmt.Println("healthcheck failed:", err)
		return
	}

	defer resp.Body.Close()

	fmt.Println("health status:", resp.Status)
}
