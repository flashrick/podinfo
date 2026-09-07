package http

import (
	"net/http"
	"os"

	"runtime"
	"strconv"

	"github.com/stefanprodan/podinfo/pkg/version"
)

func previewCommitSHA() string {
	if sha := os.Getenv("PREVIEW_COMMIT_SHA"); sha != "" {
		return sha
	}
	if version.REVISION != "" {
		return version.REVISION
	}
	return "unknown"
}

// Info godoc
// @Summary Runtime information
// @Description returns the runtime information
// @Tags HTTP API
// @Accept json
// @Produce json
// @Success 200 {object} http.RuntimeResponse
// @Router /api/info [get]
func (s *Server) infoHandler(w http.ResponseWriter, r *http.Request) {
	_, span := s.tracer.Start(r.Context(), "infoHandler")
	defer span.End()

	data := RuntimeResponse{
		Hostname:     s.config.Hostname,
		Version:      version.VERSION,
		Revision:     version.REVISION,
		App:          "podinfo",
		CommitSHA:    previewCommitSHA(),
		Logo:         s.config.UILogo,
		Color:        s.config.UIColor,
		Message:      s.config.UIMessage,
		GOOS:         runtime.GOOS,
		GOARCH:       runtime.GOARCH,
		Runtime:      runtime.Version(),
		NumGoroutine: strconv.FormatInt(int64(runtime.NumGoroutine()), 10),
		NumCPU:       strconv.FormatInt(int64(runtime.NumCPU()), 10),
	}

	s.JSONResponse(w, r, data)
}

type RuntimeResponse struct {
	Hostname     string `json:"hostname"`
	Version      string `json:"version"`
	Revision     string `json:"revision"`
	App          string `json:"app"`
	CommitSHA    string `json:"commit_sha"`
	Color        string `json:"color"`
	Logo         string `json:"logo"`
	Message      string `json:"message"`
	GOOS         string `json:"goos"`
	GOARCH       string `json:"goarch"`
	Runtime      string `json:"runtime"`
	NumGoroutine string `json:"num_goroutine"`
	NumCPU       string `json:"num_cpu"`
}
