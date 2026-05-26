/*
.\developer_efficiency_stat.exe -input .\Higo版本研发详细计划-2026.xlsx -start 202604 -end 202605

// Windows 编译命令:
//
//	go build -o developer_efficiency_stat.exe .
//
// Windows 运行方式:
//
//	.\developer_efficiency_stat.exe -input "D:\data\需求统计.xlsx" -start "202604" -end "202606" -output "D:\data\result.xlsx"

*/

package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

const (
	sheetDoneImmediately = "做完即发"
	timeLayoutYMDHMS     = "2006-01-02 15:04:05"
	timeLayoutYMDHM      = "2006-01-02 15:04"
	timeLayoutYMD        = "2006-01-02"
	timeLayoutYM         = "200601"
)

var releaseSheetRegexp = regexp.MustCompile(`^\d+\.\d+\.\d+_\d+\.\d+$`)

type demandRecord struct {
	Title       string
	Goal        string
	Requester   string
	Spec        string
	Link        string
	Start       time.Time
	End         time.Time
	SourceSheet string
	StatMonth   string
}

type monthStat struct {
	Month          string
	DeliveryCount  int
	SmallCount     int
	MediumCount    int
	LargeCount     int
	DeliveryP80Day float64
}

type demandAgg struct {
	Title     string
	Goal      string
	Requester string
	Spec      string
	Link      string
	MinStart  time.Time
	MaxEnd    time.Time
}

type debugItem struct {
	Month       string
	SourceSheet string
	Title       string
	Goal        string
	Requester   string
	Spec        string
	Link        string
	Start       time.Time
	End         time.Time
	Days        float64
}

type targetSheet struct {
	Name      string
	StatMonth string
}

func main() {
	inputPath, startMonth, endMonth, outputPath, err := parseArgs(os.Args[1:])
	if err != nil {
		exitWithErr(err)
	}

	start, end, err := validateMonthRange(startMonth, endMonth)
	if err != nil {
		exitWithErr(err)
	}

	records, err := readDemandRecords(inputPath, start, end)
	if err != nil {
		exitWithErr(err)
	}
	records = deduplicateByTitle(records)

	stats := calcStatsByMonth(records, start, end)
	debugItems := buildDebugItems(records, start, end)

	if err = writeXLSX(outputPath, stats, debugItems); err != nil {
		exitWithErr(err)
	}
}

func parseArgs(args []string) (inputPath, startMonth, endMonth, outputPath string, err error) {
	if len(args) == 0 {
		return "", "", "", "", errors.New("参数不能为空，示例: -input <xlsx文件> -start <YYYYMM> -end <YYYYMM> [-output <xlsx文件>]")
	}

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-input":
			i++
			if i >= len(args) {
				return "", "", "", "", errors.New("-input 缺少参数")
			}
			inputPath = args[i]
		case "-start":
			i++
			if i >= len(args) {
				return "", "", "", "", errors.New("-start 缺少参数")
			}
			startMonth = args[i]
		case "-end":
			i++
			if i >= len(args) {
				return "", "", "", "", errors.New("-end 缺少参数")
			}
			endMonth = args[i]
		case "-output":
			i++
			if i >= len(args) {
				return "", "", "", "", errors.New("-output 缺少参数")
			}
			outputPath = args[i]
		default:
			return "", "", "", "", fmt.Errorf("未知参数: %s", args[i])
		}
	}

	if inputPath == "" || startMonth == "" || endMonth == "" {
		return "", "", "", "", errors.New("必须提供 -input -start -end 参数")
	}

	if outputPath == "" {
		outputPath = defaultOutputPath(inputPath, startMonth, endMonth)
	}

	return inputPath, startMonth, endMonth, outputPath, nil
}

func validateMonthRange(startMonth, endMonth string) (time.Time, time.Time, error) {
	start, err := time.ParseInLocation(timeLayoutYM, startMonth, time.Local)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("开始月份格式错误: %w", err)
	}
	end, err := time.ParseInLocation(timeLayoutYM, endMonth, time.Local)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("结束月份格式错误: %w", err)
	}
	if end.Before(start) {
		return time.Time{}, time.Time{}, errors.New("结束月份不能早于开始月份")
	}
	return start, end, nil
}

func defaultOutputPath(inputPath, startMonth, endMonth string) string {
	base := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	filename := fmt.Sprintf("%s_%s_to_%s.xlsx", base, startMonth, endMonth)
	return filepath.Join(filepath.Dir(inputPath), filename)
}

func readDemandRecords(xlsxPath string, start, end time.Time) ([]demandRecord, error) {
	file, err := excelize.OpenFile(xlsxPath)
	if err != nil {
		return nil, fmt.Errorf("打开 xlsx 失败: %w", err)
	}
	defer file.Close()

	allSheetNames := file.GetSheetList()
	sheets := filterTargetSheets(allSheetNames, start, end)
	if len(sheets) == 0 {
		return nil, fmt.Errorf("未找到可统计的 sheet。请确认存在 \"%s\" 或符合 版本号_发布日期 的 sheet", sheetDoneImmediately)
	}

	var all []demandRecord
	for _, sheet := range sheets {
		records, sheetErr := parseOneSheet(file, sheet.Name, sheet.StatMonth)
		if sheetErr != nil {
			fmt.Fprintf(os.Stderr, "警告: 跳过 sheet[%s]，原因: %v\n", sheet.Name, sheetErr)
			continue
		}
		all = append(all, records...)
	}
	return all, nil
}

func filterTargetSheets(sheets []string, start, end time.Time) []targetSheet {
	var result []targetSheet
	for _, name := range sheets {
		if name == sheetDoneImmediately {
			result = append(result, targetSheet{Name: name})
			continue
		}
		if !releaseSheetRegexp.MatchString(name) {
			continue
		}

		parts := strings.Split(name, "_")
		if len(parts) != 2 {
			continue
		}

		releaseDate, err := parseReleaseDate(parts[1], start.Year())
		if err != nil {
			continue
		}

		month := time.Date(releaseDate.Year(), releaseDate.Month(), 1, 0, 0, 0, 0, time.Local)
		if !month.Before(start) && !month.After(end) {
			result = append(result, targetSheet{
				Name:      name,
				StatMonth: month.Format(timeLayoutYM),
			})
		}
	}

	return result
}

func parseReleaseDate(mmdd string, fallbackYear int) (time.Time, error) {
	parts := strings.Split(mmdd, ".")
	if len(parts) != 2 {
		return time.Time{}, errors.New("发布日期格式错误")
	}

	var month, day int
	if _, err := fmt.Sscanf(mmdd, "%d.%d", &month, &day); err != nil {
		return time.Time{}, fmt.Errorf("解析发布日期失败: %w", err)
	}
	return time.Date(fallbackYear, time.Month(month), day, 0, 0, 0, 0, time.Local), nil
}

func parseOneSheet(file *excelize.File, sheet, statMonth string) ([]demandRecord, error) {
	rows, err := file.GetRows(sheet)
	if err != nil {
		return nil, err
	}
	if len(rows) < 2 {
		return nil, nil
	}

	headerRowIdx, titleCol, goalCol, requesterCol, specCol, linkCol, startCol, endCol, err := detectHeaderRow(rows)
	if err != nil {
		return nil, err
	}

	mergeRefMap, err := buildMergeRefMap(file, sheet)
	if err != nil {
		return nil, err
	}

	demandMap := make(map[string]*demandAgg)
	for rowIdx := headerRowIdx + 2; rowIdx <= len(rows); rowIdx++ {
		titleCell := cellName(titleCol, rowIdx)
		goalCell := cellName(goalCol, rowIdx)
		requesterCell := cellName(requesterCol, rowIdx)
		specCell := cellName(specCol, rowIdx)
		startCell := cellName(startCol, rowIdx)
		endCell := cellName(endCol, rowIdx)

		title := strings.TrimSpace(readCellValueByMerge(file, sheet, titleCell, mergeRefMap))
		if title == "" {
			continue
		}
		goal := ""
		if goalCol > 0 {
			goal = strings.TrimSpace(readCellValueByMerge(file, sheet, goalCell, mergeRefMap))
		}
		requester := ""
		if requesterCol > 0 {
			requester = strings.TrimSpace(readCellValueByMerge(file, sheet, requesterCell, mergeRefMap))
		}
		spec := strings.TrimSpace(readCellValueByMerge(file, sheet, specCell, mergeRefMap))
		link := ""
		if linkCol > 0 {
			linkCell := cellName(linkCol, rowIdx)
			link = strings.TrimSpace(readCellLinkByMerge(file, sheet, linkCell, mergeRefMap))
		}

		startVal := strings.TrimSpace(getCellValue(file, sheet, startCell))
		endVal := strings.TrimSpace(getCellValue(file, sheet, endCell))
		if startVal == "" {
			continue
		}
		if endVal == "" {
			continue
		}

		startTime, errStart := parseAnyTime(startVal)
		endTime, errEnd := parseAnyTime(endVal)
		if errStart != nil || errEnd != nil {
			continue
		}
		if endTime.Before(startTime) {
			continue
		}

		key := sheet + "|" + title
		exist, ok := demandMap[key]
		if !ok {
			demandMap[key] = &demandAgg{
				Title:     title,
				Goal:      goal,
				Requester: requester,
				Spec:      spec,
				Link:      link,
				MinStart:  startTime,
				MaxEnd:    endTime,
			}
			continue
		}
		if exist.Goal == "" && goal != "" {
			exist.Goal = goal
		}
		if exist.Requester == "" && requester != "" {
			exist.Requester = requester
		}
		if exist.Spec == "" && spec != "" {
			exist.Spec = spec
		}
		if exist.Link == "" && link != "" {
			exist.Link = link
		}
		if startTime.Before(exist.MinStart) {
			exist.MinStart = startTime
		}
		if endTime.After(exist.MaxEnd) {
			exist.MaxEnd = endTime
		}
	}

	out := make([]demandRecord, 0, len(demandMap))
	for _, v := range demandMap {
		if statMonth != "" && monthKey(v.MaxEnd) != statMonth {
			continue
		}
		out = append(out, demandRecord{
			Title:       v.Title,
			Goal:        v.Goal,
			Requester:   v.Requester,
			Spec:        v.Spec,
			Link:        v.Link,
			Start:       v.MinStart,
			End:         v.MaxEnd,
			SourceSheet: sheet,
			StatMonth:   statMonth,
		})
	}
	return out, nil
}

func detectHeaderRow(rows [][]string) (headerRowIdx, titleCol, goalCol, requesterCol, specCol, linkCol, startCol, endCol int, err error) {
	maxCheck := 2
	if len(rows) < maxCheck {
		maxCheck = len(rows)
	}
	for i := 0; i < maxCheck; i++ {
		titleCol, goalCol, requesterCol, specCol, linkCol, startCol, endCol, err = headerIndexes(rows[i])
		if err == nil {
			return i, titleCol, goalCol, requesterCol, specCol, linkCol, startCol, endCol, nil
		}
	}
	return -1, -1, -1, -1, -1, -1, -1, -1, errors.New("前两行未识别到表头(需求标题/需求规格/研发开始时间/测试完成时间)")
}

func headerIndexes(header []string) (titleCol, goalCol, requesterCol, specCol, linkCol, startCol, endCol int, err error) {
	titleCol, goalCol, requesterCol, specCol, linkCol, startCol, endCol = -1, -1, -1, -1, -1, -1, -1
	for i, h := range header {
		normalized := normalizeHeader(h)
		switch {
		case isTitleHeader(normalized):
			titleCol = i + 1
		case isGoalHeader(normalized):
			goalCol = i + 1
		case isRequesterHeader(normalized):
			requesterCol = i + 1
		case isSpecHeader(normalized):
			specCol = i + 1
		case isLinkHeader(normalized):
			linkCol = i + 1
		case isStartHeader(normalized):
			startCol = i + 1
		case isEndHeader(normalized):
			endCol = i + 1
		}
	}
	if titleCol == -1 || specCol == -1 || startCol == -1 || endCol == -1 {
		return -1, -1, -1, -1, -1, -1, -1, errors.New("表头缺少 需求标题/需求规格/研发开始时间/测试完成时间")
	}
	return titleCol, goalCol, requesterCol, specCol, linkCol, startCol, endCol, nil
}

func normalizeHeader(s string) string {
	s = strings.TrimSpace(s)
	replacer := strings.NewReplacer(" ", "", "\t", "", "\n", "", "\r", "", "　", "")
	return replacer.Replace(s)
}

func isTitleHeader(h string) bool {
	return h == "需求标题" || h == "需求名称" || (strings.Contains(h, "需求") && (strings.Contains(h, "标题") || strings.Contains(h, "名称")))
}

func isStartHeader(h string) bool {
	return h == "研发开始时间" || h == "开发开始时间" ||
		(strings.Contains(h, "研发") && strings.Contains(h, "开始")) ||
		(strings.Contains(h, "开发") && strings.Contains(h, "开始"))
}

func isSpecHeader(h string) bool {
	return h == "需求规格" || h == "需求规模" || h == "规格" ||
		(strings.Contains(h, "需求") && strings.Contains(h, "规格")) ||
		(strings.Contains(h, "需求") && strings.Contains(h, "规模"))
}

func isGoalHeader(h string) bool {
	return h == "需求目标" || h == "目标" || h == "目标描述" || h == "需求目的" ||
		(strings.Contains(h, "需求") && strings.Contains(h, "目标")) ||
		(strings.Contains(h, "目标") && strings.Contains(h, "描述"))
}

func isRequesterHeader(h string) bool {
	return h == "需求人" || h == "提出人" || h == "提需求人" || h == "需求方" ||
		(strings.Contains(h, "需求") && strings.Contains(h, "人")) ||
		(strings.Contains(h, "提出") && strings.Contains(h, "人"))
}

func isLinkHeader(h string) bool {
	return h == "需求链接" || h == "链接" || h == "地址" ||
		(strings.Contains(h, "需求") && strings.Contains(h, "链接"))
}

func isEndHeader(h string) bool {
	return h == "测试完成时间" || h == "测试结束时间" ||
		(strings.Contains(h, "测试") && strings.Contains(h, "完成")) ||
		(strings.Contains(h, "测试") && strings.Contains(h, "结束"))
}

func buildMergeRefMap(file *excelize.File, sheet string) (map[string]string, error) {
	mergeCells, err := file.GetMergeCells(sheet)
	if err != nil {
		return nil, err
	}

	refMap := make(map[string]string)
	for _, mc := range mergeCells {
		startAxis, endAxis := mc.GetStartAxis(), mc.GetEndAxis()
		startCol, startRow, err1 := excelize.CellNameToCoordinates(startAxis)
		endCol, endRow, err2 := excelize.CellNameToCoordinates(endAxis)
		if err1 != nil || err2 != nil {
			continue
		}

		for col := startCol; col <= endCol; col++ {
			for row := startRow; row <= endRow; row++ {
				name, err := excelize.CoordinatesToCellName(col, row)
				if err != nil {
					continue
				}
				refMap[name] = startAxis
			}
		}
	}
	return refMap, nil
}

func readCellValueByMerge(file *excelize.File, sheet, cell string, mergeRefMap map[string]string) string {
	if ref, ok := mergeRefMap[cell]; ok {
		return getCellValue(file, sheet, ref)
	}
	return getCellValue(file, sheet, cell)
}

func readCellLinkByMerge(file *excelize.File, sheet, cell string, mergeRefMap map[string]string) string {
	targetCell := cell
	if ref, ok := mergeRefMap[cell]; ok {
		targetCell = ref
	}
	if hasLink, url, err := file.GetCellHyperLink(sheet, targetCell); err == nil && hasLink && strings.TrimSpace(url) != "" {
		return strings.TrimSpace(url)
	}
	return getCellValue(file, sheet, targetCell)
}

func cellName(col, row int) string {
	name, err := excelize.CoordinatesToCellName(col, row)
	if err != nil {
		return ""
	}
	return name
}

func getCellValue(file *excelize.File, sheet, cell string) string {
	v, err := file.GetCellValue(sheet, cell)
	if err != nil {
		return ""
	}
	return v
}

func parseAnyTime(s string) (time.Time, error) {
	layouts := []string{
		timeLayoutYMDHMS,
		timeLayoutYMDHM,
		timeLayoutYMD,
		"2006/1/2 15:04:05",
		"2006/1/2 15:04",
		"2006/1/2",
		"2006-1-2 15:04:05",
		"2006-1-2 15:04",
		"2006-1-2",
		"2006.1.2 15:04:05",
		"2006.1.2 15:04",
		"2006.1.2",
		"01-02-06 15:04:05",
		"01-02-06 15:04",
		"01-02-06",
		"01-02-2006 15:04:05",
		"01-02-2006 15:04",
		"01-02-2006",
	}

	s = strings.TrimSpace(s)
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			return t, nil
		}
	}

	// Excel 可能是纯数字日期序列
	if n, nErr := parseExcelSerial(s); nErr == nil {
		return excelize.ExcelDateToTime(n, false)
	}

	return time.Time{}, fmt.Errorf("无法识别的时间格式: %s", s)
}

func parseExcelSerial(s string) (float64, error) {
	s = strings.TrimSpace(s)
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func calcStatsByMonth(records []demandRecord, start, end time.Time) []monthStat {
	byMonth := make(map[string][]float64)
	sizeCountByMonth := make(map[string]map[string]int)
	for _, rec := range records {
		month := statMonthOf(rec)
		days := inclusiveWorkdayDiff(rec.Start, rec.End)
		if days < 0 {
			continue
		}
		byMonth[month] = append(byMonth[month], days)
		if _, ok := sizeCountByMonth[month]; !ok {
			sizeCountByMonth[month] = map[string]int{
				"small":  0,
				"medium": 0,
				"large":  0,
			}
		}
		switch normalizeSpec(rec.Spec) {
		case "small":
			sizeCountByMonth[month]["small"]++
		case "medium":
			sizeCountByMonth[month]["medium"]++
		case "large":
			sizeCountByMonth[month]["large"]++
		}
	}

	months := monthRange(start, end)
	stats := make([]monthStat, 0, len(months))
	for _, m := range months {
		durations := byMonth[m]
		p80 := 0.0
		if len(durations) > 0 {
			p80 = percentile(durations, 0.8)
		}
		stats = append(stats, monthStat{
			Month:          m,
			DeliveryCount:  len(durations),
			SmallCount:     sizeCountByMonth[m]["small"],
			MediumCount:    sizeCountByMonth[m]["medium"],
			LargeCount:     sizeCountByMonth[m]["large"],
			DeliveryP80Day: p80,
		})
	}
	return stats
}

func buildDebugItems(records []demandRecord, start, end time.Time) []debugItem {
	items := make([]debugItem, 0, len(records))
	for _, rec := range records {
		month := statMonthOf(rec)
		if month < start.Format(timeLayoutYM) || month > end.Format(timeLayoutYM) {
			continue
		}
		days := inclusiveWorkdayDiff(rec.Start, rec.End)
		if days < 0 {
			continue
		}
		items = append(items, debugItem{
			Month:       month,
			SourceSheet: rec.SourceSheet,
			Title:       rec.Title,
			Goal:        rec.Goal,
			Requester:   rec.Requester,
			Spec:        rec.Spec,
			Link:        rec.Link,
			Start:       rec.Start,
			End:         rec.End,
			Days:        days,
		})
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Month != items[j].Month {
			return items[i].Month < items[j].Month
		}
		if items[i].Days != items[j].Days {
			return items[i].Days < items[j].Days
		}
		if !items[i].Start.Equal(items[j].Start) {
			return items[i].Start.Before(items[j].Start)
		}
		return items[i].Title < items[j].Title
	})
	return items
}

func monthKey(t time.Time) string {
	return t.Format(timeLayoutYM)
}

func statMonthOf(rec demandRecord) string {
	if rec.StatMonth != "" {
		return rec.StatMonth
	}
	return monthKey(rec.End)
}

func monthRange(start, end time.Time) []string {
	var out []string
	cur := time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, time.Local)
	last := time.Date(end.Year(), end.Month(), 1, 0, 0, 0, 0, time.Local)
	for !cur.After(last) {
		out = append(out, cur.Format(timeLayoutYM))
		cur = cur.AddDate(0, 1, 0)
	}
	return out
}

func inclusiveWorkdayDiff(start, end time.Time) float64 {
	startDate := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.Local)
	endDate := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.Local)
	if endDate.Before(startDate) {
		return -1
	}

	weekdayCount := 0.0
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
			continue
		}
		weekdayCount++
	}
	return weekdayCount
}

func percentile(values []float64, p float64) float64 {
	if len(values) == 0 {
		return 0
	}
	cp := make([]float64, len(values))
	copy(cp, values)
	sort.Float64s(cp)

	index := int(math.Ceil(p*float64(len(cp)))) - 1
	if index < 0 {
		index = 0
	}
	if index >= len(cp) {
		index = len(cp) - 1
	}
	return cp[index]
}

func normalizeSpec(spec string) string {
	s := normalizeHeader(spec)
	switch s {
	case "小", "small", "s":
		return "small"
	case "中", "medium", "m":
		return "medium"
	case "大", "large", "l":
		return "large"
	default:
		return ""
	}
}

func specRank(spec string) int {
	switch normalizeSpec(spec) {
	case "small":
		return 1
	case "medium":
		return 2
	case "large":
		return 3
	default:
		return 0
	}
}

func deduplicateByTitle(records []demandRecord) []demandRecord {
	bestByTitle := make(map[string]demandRecord)
	order := make([]string, 0, len(records))
	for _, rec := range records {
		title := strings.TrimSpace(rec.Title)
		if title == "" {
			continue
		}
		rec.Title = title
		exist, ok := bestByTitle[title]
		if !ok {
			bestByTitle[title] = rec
			order = append(order, title)
			continue
		}
		if shouldReplaceDemand(exist, rec) {
			bestByTitle[title] = rec
		}
	}

	out := make([]demandRecord, 0, len(bestByTitle))
	for _, title := range order {
		if rec, ok := bestByTitle[title]; ok {
			out = append(out, rec)
		}
	}
	return out
}

func shouldReplaceDemand(current, candidate demandRecord) bool {
	curRank := specRank(current.Spec)
	candRank := specRank(candidate.Spec)
	if candRank != curRank {
		return candRank > curRank
	}

	curDays := inclusiveWorkdayDiff(current.Start, current.End)
	candDays := inclusiveWorkdayDiff(candidate.Start, candidate.End)
	if curDays < 0 {
		return candDays >= 0
	}
	if candDays < 0 {
		return false
	}
	if candDays != curDays {
		return candDays < curDays
	}
	return false
}

func writeXLSX(outputPath string, stats []monthStat, items []debugItem) error {
	f := excelize.NewFile()
	summarySheet := "汇总"
	detailSheet := "明细"
	f.SetSheetName("Sheet1", summarySheet)
	if _, err := f.NewSheet(detailSheet); err != nil {
		return fmt.Errorf("创建输出工作表失败: %w", err)
	}

	summaryHeader := []string{"统计月份", "交付数量", "小", "中", "大", "交付时长P80"}
	for col, h := range summaryHeader {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		if err := f.SetCellValue(summarySheet, cell, h); err != nil {
			return fmt.Errorf("写入汇总表头失败: %w", err)
		}
	}
	for i, s := range stats {
		row := i + 2
		values := []interface{}{s.Month, s.DeliveryCount, s.SmallCount, s.MediumCount, s.LargeCount, s.DeliveryP80Day}
		for col, v := range values {
			cell, _ := excelize.CoordinatesToCellName(col+1, row)
			if err := f.SetCellValue(summarySheet, cell, v); err != nil {
				return fmt.Errorf("写入汇总数据失败: %w", err)
			}
		}
	}

	detailHeader := []string{"统计月份", "来源Sheet", "需求标题", "需求目标", "需求人", "需求规格", "需求链接", "需求开始时间", "需求结束时间", "交付时长"}
	for col, h := range detailHeader {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		if err := f.SetCellValue(detailSheet, cell, h); err != nil {
			return fmt.Errorf("写入明细表头失败: %w", err)
		}
	}
	for i, item := range items {
		row := i + 2
		values := []interface{}{
			item.Month,
			item.SourceSheet,
			item.Title,
			item.Goal,
			item.Requester,
			item.Spec,
			item.Link,
			item.Start.Format(timeLayoutYMD),
			item.End.Format(timeLayoutYMD),
			item.Days,
		}
		for col, v := range values {
			cell, _ := excelize.CoordinatesToCellName(col+1, row)
			if err := f.SetCellValue(detailSheet, cell, v); err != nil {
				return fmt.Errorf("写入明细数据失败: %w", err)
			}
		}
		if strings.TrimSpace(item.Link) != "" {
			linkCell, _ := excelize.CoordinatesToCellName(7, row)
			if err := f.SetCellHyperLink(detailSheet, linkCell, strings.TrimSpace(item.Link), "External"); err != nil {
				return fmt.Errorf("写入需求链接失败: %w", err)
			}
		}
	}

	if idx, err := f.GetSheetIndex(summarySheet); err == nil {
		f.SetActiveSheet(idx)
	}
	if err := f.SaveAs(outputPath); err != nil {
		return fmt.Errorf("保存输出文件失败: %w", err)
	}
	return nil
}

func exitWithErr(err error) {
	fmt.Fprintf(os.Stderr, "错误: %v\n", err)
	os.Exit(1)
}
