package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	apperror "service/internal/pkg/error"
	"service/internal/pkg/form"
	reportservice "service/internal/report/service"
)

// ReportAPIHandler handles REST API endpoints for daily reports.
type ReportAPIHandler struct {
	queryService  *reportservice.ReportQueryService
	exportService *reportservice.ReportExportService
}

func NewReportAPIHandler(
	queryService *reportservice.ReportQueryService,
	exportService *reportservice.ReportExportService,
) *ReportAPIHandler {
	return &ReportAPIHandler{
		queryService:  queryService,
		exportService: exportService,
	}
}

// List handles GET /api/v1/reports
func (h *ReportAPIHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	f := form.ReportListForm{
		Sort:  q.Get("sort"),
		Page:  queryInt(q.Get("page"), 1),
		Limit: queryInt(q.Get("limit"), 20),
	}
	f.Defaults()

	data, meta, err := h.queryService.List(&f)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": data, "meta": meta})
}

// ListByUser handles GET /api/v1/reports/by-user
func (h *ReportAPIHandler) ListByUser(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	f := form.ReportByUserForm{
		FromDate: q.Get("fromDate"),
		ToDate:   q.Get("toDate"),
		SortBy:   q.Get("sortBy"),
		Sort:     q.Get("sort"),
		Page:     queryInt(q.Get("page"), 1),
		Limit:    queryInt(q.Get("limit"), 20),
	}
	f.Defaults()

	if f.FromDate == "" || f.ToDate == "" {
		writeJSONError(w, http.StatusBadRequest, "fromDate and toDate are required")
		return
	}

	data, meta, err := h.queryService.ListByUser(&f)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": data, "meta": meta})
}

// ListByDate handles GET /api/v1/reports/by-date
func (h *ReportAPIHandler) ListByDate(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	f := form.ReportByDateForm{
		FromDate: q.Get("fromDate"),
		ToDate:   q.Get("toDate"),
		Page:     queryInt(q.Get("page"), 1),
		Limit:    queryInt(q.Get("limit"), 20),
	}
	f.Defaults()

	if f.FromDate == "" || f.ToDate == "" {
		writeJSONError(w, http.StatusBadRequest, "fromDate and toDate are required")
		return
	}

	data, meta, err := h.queryService.ListByDate(&f)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": data, "meta": meta})
}

// GetByID handles GET /api/v1/reports/{id}
func (h *ReportAPIHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid report id")
		return
	}

	data, err := h.queryService.GetByID(uint(id64))
	if err != nil {
		if errors.Is(err, apperror.ErrReportNotFound) {
			writeJSONError(w, http.StatusNotFound, "report not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": data})
}

// Export handles GET /api/v1/reports/export and returns a .docx file download.
func (h *ReportAPIHandler) Export(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	f := form.ReportExportForm{
		FromDate: q.Get("fromDate"),
		ToDate:   q.Get("toDate"),
		Sort:     q.Get("sort"),
		GroupBy:  q.Get("groupBy"),
	}
	f.Defaults()

	docxBytes, err := h.exportService.Export(&f)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	w.Header().Set("Content-Disposition", `attachment; filename="daily-report-export.docx"`)
	w.WriteHeader(http.StatusOK)
	w.Write(docxBytes)
}

// --- helpers ---

func queryInt(s string, defaultVal int) int {
	if v, err := strconv.Atoi(s); err == nil {
		return v
	}
	return defaultVal
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeJSONError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
