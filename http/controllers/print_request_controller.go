package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/render"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"

	"safa-went/database/models"
	"safa-went/http/requests"
	"safa-went/http/resources"
	"safa-went/internal/responses"
)

type PrintRequest struct {
	DB *gorm.DB
}

func (c *PrintRequest) findRaw(id uint, item *models.PrintRequest) error {
	return c.DB.Preload("Requester").Preload("Approver").First(item, id).Error
}

// GetAllPrintRequest godoc
// @Summary List all PrintRequest records
// @Description Get a paginated list of print requests with optional filters
// @Tags PrintRequest
// @Produce json
// @Security BearerAuth
// @Param page               query int    false "Page number (default 1)"
// @Param per_page           query int    false "Items per page (default 15)"
// @Param requester_names    query string false "Comma-separated requester names"
// @Param approver_names     query string false "Comma-separated approver names"
// @Param color_copies_min   query int    false "Minimum color copies"
// @Param color_copies_max   query int    false "Maximum color copies"
// @Param bw_copies_min      query int    false "Minimum BW copies"
// @Param bw_copies_max      query int    false "Maximum BW copies"
// @Param requested_at_from  query string false "Start date (RFC3339)"
// @Param requested_at_to    query string false "End date (RFC3339)"
// @Param sort               query string false "Sort column"
// @Success 200 {object} resources.PrintRequestCollection
// @Failure 500 {object} responses.ErrorBody
// @Router /print-requests [get]
func (c *PrintRequest) GetAllPrintRequest(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	perPage, _ := strconv.ParseInt(r.URL.Query().Get("per_page"), 10, 64)

	requesterNames := splitCSV(r.URL.Query().Get("requester_names"))
	approverNames := splitCSV(r.URL.Query().Get("approver_names"))

	colorCopiesMin, _ := strconv.Atoi(r.URL.Query().Get("color_copies_min"))
	colorCopiesMax, _ := strconv.Atoi(r.URL.Query().Get("color_copies_max"))
	bwCopiesMin, _ := strconv.Atoi(r.URL.Query().Get("bw_copies_min"))
	bwCopiesMax, _ := strconv.Atoi(r.URL.Query().Get("bw_copies_max"))

	var requestedAtFrom, requestedAtTo *time.Time
	if v := r.URL.Query().Get("requested_at_from"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			requestedAtFrom = &t
		}
	}
	if v := r.URL.Query().Get("requested_at_to"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			requestedAtTo = &t
		}
	}

	sortBy := r.URL.Query().Get("sort")

	q := resources.NewPrintRequestQuery(c.DB)
	collection, err := q.Filter(
		requesterNames, approverNames,
		colorCopiesMin, colorCopiesMax,
		bwCopiesMin, bwCopiesMax,
		requestedAtFrom, requestedAtTo,
		sortBy, page, perPage,
	)
	if err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	render.JSON(w, r, collection)
}

// GetPrintRequestByID godoc
// @Summary Get a PrintRequest by ID
// @Tags PrintRequest
// @Produce json
// @Security BearerAuth
// @Param id path int true "PrintRequest ID"
// @Success 200 {object} resources.PrintRequestResource
// @Failure 400 {object} responses.ErrorBody
// @Failure 404 {object} responses.ErrorBody
// @Router /print-requests/{id} [get]
func (c *PrintRequest) GetPrintRequestByID(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid id")
		return
	}
	item, err := resources.NewPrintRequestQuery(c.DB).Find(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			responses.JSONError(w, r, http.StatusNotFound, "not found")
			return
		}
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	render.JSON(w, r, item)
}

// CreatePrintRequest godoc
// @Summary Create a PrintRequest
// @Tags PrintRequest
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body requests.PrintRequestPayload true "PrintRequest payload"
// @Success 201 {object} resources.PrintRequestResource
// @Failure 400 {object} responses.ErrorBody
// @Failure 500 {object} responses.ErrorBody
// @Router /print-requests [post]
func (c *PrintRequest) CreatePrintRequest(w http.ResponseWriter, r *http.Request) {
	var payload requests.PrintRequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid request body")
		return
	}
	item := models.PrintRequest{
		RequestedAt: payload.RequestedAt,
		ColorCopies: payload.ColorCopies,
		BWCopies:    payload.BWCopies,
		Description: payload.Description,
		RequesterID: payload.RequesterID,
		ApproverID:  payload.ApproverID,
	}
	if err := c.DB.Create(&item).Error; err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	created, err := resources.NewPrintRequestQuery(c.DB).Find(item.ID)
	if err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, created)
}

// UpdatePrintRequest godoc
// @Summary Update a PrintRequest
// @Tags PrintRequest
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "PrintRequest ID"
// @Param payload body requests.PrintRequestUpdatePayload true "PrintRequest update payload"
// @Success 200 {object} resources.PrintRequestResource
// @Failure 400 {object} responses.ErrorBody
// @Failure 404 {object} responses.ErrorBody
// @Router /print-requests/{id} [put]
func (c *PrintRequest) UpdatePrintRequest(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid id")
		return
	}
	var item models.PrintRequest
	if err := c.findRaw(id, &item); err != nil {
		if err == gorm.ErrRecordNotFound {
			responses.JSONError(w, r, http.StatusNotFound, "not found")
			return
		}
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	var payload requests.PrintRequestUpdatePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid request body")
		return
	}
	updates := map[string]interface{}{}
	if payload.RequestedAt != nil {
		updates["requested_at"] = *payload.RequestedAt
	}
	if payload.ColorCopies != nil {
		updates["color_copies"] = *payload.ColorCopies
	}
	if payload.BWCopies != nil {
		updates["bw_copies"] = *payload.BWCopies
	}
	if payload.Description != nil {
		updates["description"] = *payload.Description
	}
	if payload.RequesterID != nil {
		updates["requester_id"] = *payload.RequesterID
	}
	if payload.ApproverID != nil {
		updates["approver_id"] = *payload.ApproverID
	}
	if len(updates) == 0 {
		responses.JSONError(w, r, http.StatusBadRequest, "no fields to update")
		return
	}
	if err := c.DB.Model(&item).Updates(updates).Error; err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	updated, err := resources.NewPrintRequestQuery(c.DB).Find(id)
	if err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	render.JSON(w, r, updated)
}

// DeletePrintRequest godoc
// @Summary Delete a PrintRequest
// @Tags PrintRequest
// @Security BearerAuth
// @Param id path int true "PrintRequest ID"
// @Success 204
// @Failure 400 {object} responses.ErrorBody
// @Failure 404 {object} responses.ErrorBody
// @Router /print-requests/{id} [delete]
func (c *PrintRequest) DeletePrintRequest(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid id")
		return
	}
	var item models.PrintRequest
	if err := c.findRaw(id, &item); err != nil {
		if err == gorm.ErrRecordNotFound {
			responses.JSONError(w, r, http.StatusNotFound, "not found")
			return
		}
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	if err := c.DB.Delete(&item).Error; err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ExportByRequester generates an Excel report of print requests grouped by requester for a date range.
func (c *PrintRequest) ExportByRequester(w http.ResponseWriter, r *http.Request) {
	startStr := r.URL.Query().Get("start_date")
	endStr := r.URL.Query().Get("end_date")
	if startStr == "" || endStr == "" {
		responses.JSONError(w, r, http.StatusBadRequest, "start_date and end_date are required")
		return
	}
	startDate, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid start_date format, use YYYY-MM-DD")
		return
	}
	endDateRaw, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid end_date format, use YYYY-MM-DD")
		return
	}
	endDate := endDateRaw.Add(24*time.Hour - time.Second)

	var rows []struct {
		RequesterName  string
		ColorCopies    int
		BWCopies       int
		TotalAllCopies int
	}
	if err := c.DB.Table("print_requests").
		Select("requesters.name as requester_name, SUM(print_requests.color_copies) as color_copies, SUM(print_requests.bw_copies) as bw_copies, SUM(print_requests.color_copies + print_requests.bw_copies) as total_all_copies").
		Joins("JOIN requesters ON requesters.id = print_requests.requester_id").
		Where("print_requests.deleted_at IS NULL AND print_requests.requested_at BETWEEN ? AND ?", startDate, endDate).
		Group("requesters.id, requesters.name").
		Order("total_all_copies DESC").
		Scan(&rows).Error; err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	f := excelize.NewFile()
	defer f.Close()
	sheet := "Rapor"
	f.SetSheetName("Sheet1", sheet)

	// Row 1: Title merged A1:D1, bold 14pt centered
	title := fmt.Sprintf("Fotokopi Raporu (%s - %s)", startDate.Format("02.01.2006"), endDateRaw.Format("02.01.2006"))
	f.SetCellValue(sheet, "A1", title)
	f.MergeCell(sheet, "A1", "D1")
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 14},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	f.SetCellStyle(sheet, "A1", "D1", titleStyle)
	f.SetRowHeight(sheet, 1, 20)

	// Row 2: Column headers, bold
	headers := []string{"Talep Eden Kişi", "Renkli Kopya", "Siyah-Beyaz Kopya", "Toplam Kopya"}
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 2)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}
	f.SetRowHeight(sheet, 2, 20)

	// Rows 3+: Data
	for i, row := range rows {
		rn := i + 3
		f.SetCellValue(sheet, fmt.Sprintf("A%d", rn), row.RequesterName)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", rn), row.ColorCopies)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", rn), row.BWCopies)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", rn), row.TotalAllCopies)
	}

	// Totals row
	lastData := len(rows) + 2
	totalsRow := lastData + 1
	boldStyle, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
	f.SetCellValue(sheet, fmt.Sprintf("A%d", totalsRow), "Genel Toplam")
	f.SetCellFormula(sheet, fmt.Sprintf("B%d", totalsRow), fmt.Sprintf("SUM(B3:B%d)", lastData))
	f.SetCellFormula(sheet, fmt.Sprintf("C%d", totalsRow), fmt.Sprintf("SUM(C3:C%d)", lastData))
	f.SetCellFormula(sheet, fmt.Sprintf("D%d", totalsRow), fmt.Sprintf("SUM(D3:D%d)", lastData))
	f.SetCellStyle(sheet, fmt.Sprintf("A%d", totalsRow), fmt.Sprintf("D%d", totalsRow), boldStyle)

	fileName := fmt.Sprintf("fotokopi_raporu_%s_%s.xlsx", startStr, endStr)
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	if writeErr := f.Write(w); writeErr != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, writeErr.Error())
	}
}

// ExportComparison generates an Excel comparison report for two date ranges.
func (c *PrintRequest) ExportComparison(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	p1StartStr := q.Get("first_start_date")
	p1EndStr := q.Get("first_end_date")
	p2StartStr := q.Get("second_start_date")
	p2EndStr := q.Get("second_end_date")
	if p1StartStr == "" || p1EndStr == "" || p2StartStr == "" || p2EndStr == "" {
		responses.JSONError(w, r, http.StatusBadRequest, "first_start_date, first_end_date, second_start_date, second_end_date are required")
		return
	}

	parseDate := func(s string) (time.Time, error) { return time.Parse("2006-01-02", s) }
	p1Start, err := parseDate(p1StartStr)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid first_start_date")
		return
	}
	p1EndRaw, err := parseDate(p1EndStr)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid first_end_date")
		return
	}
	p2Start, err := parseDate(p2StartStr)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid second_start_date")
		return
	}
	p2EndRaw, err := parseDate(p2EndStr)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid second_end_date")
		return
	}
	p1End := p1EndRaw.Add(24*time.Hour - time.Second)
	p2End := p2EndRaw.Add(24*time.Hour - time.Second)

	type PeriodRow struct {
		RequesterID    uint
		RequesterName  string
		ColorCopies    int
		BWCopies       int
		TotalAllCopies int
	}
	fetchPeriod := func(start, end time.Time) ([]PeriodRow, error) {
		var rows []PeriodRow
		err := c.DB.Table("print_requests").
			Select("requesters.id as requester_id, requesters.name as requester_name, SUM(print_requests.color_copies) as color_copies, SUM(print_requests.bw_copies) as bw_copies, SUM(print_requests.color_copies + print_requests.bw_copies) as total_all_copies").
			Joins("JOIN requesters ON requesters.id = print_requests.requester_id").
			Where("print_requests.deleted_at IS NULL AND print_requests.requested_at BETWEEN ? AND ?", start, end).
			Group("requesters.id, requesters.name").
			Scan(&rows).Error
		return rows, err
	}

	p1Rows, err := fetchPeriod(p1Start, p1End)
	if err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	p2Rows, err := fetchPeriod(p2Start, p2End)
	if err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// Build map requesterID -> row for each period
	p1Map := make(map[uint]PeriodRow, len(p1Rows))
	for _, row := range p1Rows {
		p1Map[row.RequesterID] = row
	}
	p2Map := make(map[uint]PeriodRow, len(p2Rows))
	for _, row := range p2Rows {
		p2Map[row.RequesterID] = row
	}

	// Collect all requester IDs in union
	type CombinedRow struct {
		RequesterName                string
		P1Color, P1BW, P1Total       int
		P2Color, P2BW, P2Total       int
		DiffColor, DiffBW, DiffTotal int
	}
	seen := make(map[uint]bool)
	var combined []CombinedRow
	for _, row := range p1Rows {
		if seen[row.RequesterID] {
			continue
		}
		seen[row.RequesterID] = true
		p2 := p2Map[row.RequesterID]
		combined = append(combined, CombinedRow{
			RequesterName: row.RequesterName,
			P1Color:       row.ColorCopies, P1BW: row.BWCopies, P1Total: row.TotalAllCopies,
			P2Color: p2.ColorCopies, P2BW: p2.BWCopies, P2Total: p2.TotalAllCopies,
			DiffColor: p2.ColorCopies - row.ColorCopies,
			DiffBW:    p2.BWCopies - row.BWCopies,
			DiffTotal: p2.TotalAllCopies - row.TotalAllCopies,
		})
	}
	for _, row := range p2Rows {
		if seen[row.RequesterID] {
			continue
		}
		seen[row.RequesterID] = true
		combined = append(combined, CombinedRow{
			RequesterName: row.RequesterName,
			P2Color:       row.ColorCopies, P2BW: row.BWCopies, P2Total: row.TotalAllCopies,
			DiffColor: row.ColorCopies, DiffBW: row.BWCopies, DiffTotal: row.TotalAllCopies,
		})
	}
	// Sort by total_difference DESC
	for i := 0; i < len(combined)-1; i++ {
		for j := i + 1; j < len(combined); j++ {
			if combined[j].DiffTotal > combined[i].DiffTotal {
				combined[i], combined[j] = combined[j], combined[i]
			}
		}
	}

	f := excelize.NewFile()
	defer f.Close()
	sheet := "Karşılaştırma"
	f.SetSheetName("Sheet1", sheet)

	p1Label := fmt.Sprintf("%s - %s", p1Start.Format("02.01.2006"), p1EndRaw.Format("02.01.2006"))
	p2Label := fmt.Sprintf("%s - %s", p2Start.Format("02.01.2006"), p2EndRaw.Format("02.01.2006"))

	// Row 1: Title merged A1:J1, bold 14pt centered wrap, row height 40
	titleText := fmt.Sprintf("Fotokopi Karşılaştırma Raporu\n%s vs %s", p1Label, p2Label)
	f.SetCellValue(sheet, "A1", titleText)
	f.MergeCell(sheet, "A1", "J1")
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 14},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
	})
	f.SetCellStyle(sheet, "A1", "J1", titleStyle)
	f.SetRowHeight(sheet, 1, 40)

	// Row 2: empty

	// Row 3: Headers, bold, centered, wrap text, row height 40
	headers := []string{
		"Talep Eden Kişi",
		fmt.Sprintf("1. Dönem Renkli\n(%s)", p1Label),
		fmt.Sprintf("1. Dönem S-B\n(%s)", p1Label),
		fmt.Sprintf("1. Dönem Toplam\n(%s)", p1Label),
		fmt.Sprintf("2. Dönem Renkli\n(%s)", p2Label),
		fmt.Sprintf("2. Dönem S-B\n(%s)", p2Label),
		fmt.Sprintf("2. Dönem Toplam\n(%s)", p2Label),
		"Renkli Fark",
		"S-B Fark",
		"Toplam Fark",
	}
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
	})
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 3)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}
	f.SetRowHeight(sheet, 3, 40)

	// Rows 4+: Data
	for i, row := range combined {
		rn := i + 4
		f.SetCellValue(sheet, fmt.Sprintf("A%d", rn), row.RequesterName)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", rn), row.P1Color)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", rn), row.P1BW)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", rn), row.P1Total)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", rn), row.P2Color)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", rn), row.P2BW)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", rn), row.P2Total)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", rn), row.DiffColor)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", rn), row.DiffBW)
		f.SetCellValue(sheet, fmt.Sprintf("J%d", rn), row.DiffTotal)
	}

	// Totals row
	lastData := len(combined) + 3
	totalsRow := lastData + 1
	boldStyle, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
	f.SetCellValue(sheet, fmt.Sprintf("A%d", totalsRow), "Genel Toplam")
	for _, col := range []string{"B", "C", "D", "E", "F", "G", "H", "I", "J"} {
		f.SetCellFormula(sheet, fmt.Sprintf("%s%d", col, totalsRow), fmt.Sprintf("SUM(%s4:%s%d)", col, col, lastData))
	}
	f.SetCellStyle(sheet, fmt.Sprintf("A%d", totalsRow), fmt.Sprintf("J%d", totalsRow), boldStyle)

	// Conditional formatting for difference columns H-J (green if >0, red if <0)
	greenStyle, _ := f.NewConditionalStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"C6EFCE"}, Pattern: 1},
	})
	redStyle, _ := f.NewConditionalStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"FFC7CE"}, Pattern: 1},
	})
	if lastData >= 4 {
		diffRange := fmt.Sprintf("H4:J%d", lastData)
		f.SetConditionalFormat(sheet, diffRange, []excelize.ConditionalFormatOptions{
			{Type: "cell", Criteria: ">", Value: "0", Format: &greenStyle},
			{Type: "cell", Criteria: "<", Value: "0", Format: &redStyle},
		})
	}

	// Column widths
	f.SetColWidth(sheet, "A", "A", 20)
	f.SetColWidth(sheet, "B", "J", 12)

	fileName := fmt.Sprintf("karsilastirma_%s_%s_vs_%s_%s.xlsx", p1StartStr, p1EndStr, p2StartStr, p2EndStr)
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	if writeErr := f.Write(w); writeErr != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, writeErr.Error())
	}
}

// ExportAll generates an Excel export of all print requests.
func (c *PrintRequest) ExportAll(w http.ResponseWriter, r *http.Request) {
	var items []models.PrintRequest
	if err := c.DB.Preload("Requester").Preload("Approver").
		Where("deleted_at IS NULL").
		Order("requested_at DESC").
		Find(&items).Error; err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	f := excelize.NewFile()
	defer f.Close()
	sheet := "Tüm Talepler"
	f.SetSheetName("Sheet1", sheet)

	// Row 1: Title merged A1:I1
	titleText := "Tüm Fotokopi Talepleri Raporu"
	f.SetCellValue(sheet, "A1", titleText)
	f.MergeCell(sheet, "A1", "I1")
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 14},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	f.SetCellStyle(sheet, "A1", "I1", titleStyle)
	f.SetRowHeight(sheet, 1, 20)

	// Row 2: Column headers, bold
	headers := []string{"ID", "Talep Eden Kişi", "Onaylayan Kişi", "Renkli Kopya", "Siyah-Beyaz Kopya", "Toplam Kopya", "Açıklama", "Talep Tarihi", "Durum"}
	headerStyle, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 2)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}
	f.SetRowHeight(sheet, 2, 20)

	// Rows 3+: Data
	for i, item := range items {
		rn := i + 3
		requesterName := item.Requester.Name
		if requesterName == "" {
			requesterName = "Bilinmeyen Talep Eden"
		}
		approverName := item.Approver.Name
		if approverName == "" {
			approverName = "Henüz Onaylanmamış"
		}
		desc := "Açıklama yok"
		if item.Description != nil && *item.Description != "" {
			desc = *item.Description
		}
		status := "Onaylandı"
		if item.ApproverID == 0 {
			status = "Beklemede"
		}
		f.SetCellValue(sheet, fmt.Sprintf("A%d", rn), item.ID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", rn), requesterName)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", rn), approverName)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", rn), item.ColorCopies)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", rn), item.BWCopies)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", rn), item.ColorCopies+item.BWCopies)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", rn), desc)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", rn), item.RequestedAt.Format("02.01.2006 15:04"))
		f.SetCellValue(sheet, fmt.Sprintf("I%d", rn), status)
	}

	// Totals row (D, E, F columns summed)
	lastData := len(items) + 2
	totalsRow := lastData + 1
	boldStyle, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
	f.SetCellValue(sheet, fmt.Sprintf("A%d", totalsRow), "Genel Toplam")
	f.SetCellFormula(sheet, fmt.Sprintf("D%d", totalsRow), fmt.Sprintf("SUM(D3:D%d)", lastData))
	f.SetCellFormula(sheet, fmt.Sprintf("E%d", totalsRow), fmt.Sprintf("SUM(E3:E%d)", lastData))
	f.SetCellFormula(sheet, fmt.Sprintf("F%d", totalsRow), fmt.Sprintf("SUM(F3:F%d)", lastData))
	f.SetCellStyle(sheet, fmt.Sprintf("A%d", totalsRow), fmt.Sprintf("I%d", totalsRow), boldStyle)

	// Auto filter on row 2
	f.AutoFilter(sheet, fmt.Sprintf("A2:I%d", lastData), []excelize.AutoFilterOptions{})

	// Column widths
	f.SetColWidth(sheet, "A", "A", 8)
	f.SetColWidth(sheet, "B", "B", 20)
	f.SetColWidth(sheet, "C", "C", 20)
	f.SetColWidth(sheet, "D", "D", 12)
	f.SetColWidth(sheet, "E", "E", 15)
	f.SetColWidth(sheet, "F", "F", 12)
	f.SetColWidth(sheet, "G", "G", 30)
	f.SetColWidth(sheet, "H", "H", 15)
	f.SetColWidth(sheet, "I", "I", 12)

	fileName := fmt.Sprintf("tum_fotokopi_talepleri_%s.xlsx", time.Now().Format("2006-01-02"))
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	if writeErr := f.Write(w); writeErr != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, writeErr.Error())
	}
}

// splitCSV splits a comma-separated string into a slice of non-empty strings.
func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := make([]string, 0)
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			part := s[start:i]
			if part != "" {
				parts = append(parts, part)
			}
			start = i + 1
		}
	}
	return parts
}
