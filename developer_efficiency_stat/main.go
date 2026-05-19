/*
.\developer_efficiency_stat.exe -input .\Higo版本研发详细计划-2026.xlsx -start 2026-04 -end 2026-05 -debug

// Windows 编译命令:
//
//	go build -o developer_efficiency_stat.exe .
//
// Windows 运行方式:
//
//	.\developer_efficiency_stat.exe -input "D:\data\需求统计.xlsx" -start "2026-04" -end "2026-06" -output "D:\data\result.csv"
//
// 开启调试输出(输出每个需求明细并按交付时长升序):
//
//	.\developer_efficiency_stat.exe -input "D:\data\需求统计.xlsx" -start "2026-04" -end "2026-04" -output "D:\data\result.csv" -debug -debug-output "D:\data\debug.csv"

*/

package main

import (
	"encoding/csv"
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
	timeLayoutYM         = "2006-01"
)

var releaseSheetRegexp = regexp.MustCompile(`^\d+\.\d+\.\d+_\d+\.\d+$`)

type demandRecord struct {
	Title       string
	Start       time.Time
	End         time.Time
	SourceSheet string
	StatMonth   string
}

type monthStat struct {
	Month          string
	DeliveryCount  int
	DeliveryP80Day float64
}

type demandAgg struct {
	Title    string
	MinStart time.Time
	MaxEnd   time.Time
}

type debugItem struct {
	Month       string
	SourceSheet string
	Title       string
	Start       time.Time
	End         time.Time
	Days        float64
}

type targetSheet struct {
	Name      string
	StatMonth string
}

func main() {
	inputPath, startMonth, endMonth, outputPath, debugMode, debugOutputPath, err := parseArgs(os.Args[1:])
	if err != nil {
		exitWithErr(err)
	}

	start, end, err := validateMonthRange(startMonth, endMonth)
	if err != nil {
		exitWithErr(err)
	}

	records, err := readDemandRecords(inputPath, start, end, debugMode)
	if err != nil {
		exitWithErr(err)
	}

	stats := calcStatsByMonth(records, start, end)

	if err = writeCSV(outputPath, stats); err != nil {
		exitWithErr(err)
	}

	if debugMode {
		debugItems := buildDebugItems(records, start, end)
		if err = writeDebugCSV(debugOutputPath, debugItems); err != nil {
			exitWithErr(err)
		}
	}
}

func parseArgs(args []string) (inputPath, startMonth, endMonth, outputPath string, debugMode bool, debugOutputPath string, err error) {
	if len(args) == 0 {
		return "", "", "", "", false, "", errors.New("参数不能为空，示例: -input <xlsx文件> -start <YYYY-MM> -end <YYYY-MM> [-output <csv文件>] [-debug] [-debug-output <csv文件>]")
	}

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-input":
			i++
			if i >= len(args) {
				return "", "", "", "", false, "", errors.New("-input 缺少参数")
			}
			inputPath = args[i]
		case "-start":
			i++
			if i >= len(args) {
				return "", "", "", "", false, "", errors.New("-start 缺少参数")
			}
			startMonth = args[i]
		case "-end":
			i++
			if i >= len(args) {
				return "", "", "", "", false, "", errors.New("-end 缺少参数")
			}
			endMonth = args[i]
		case "-output":
			i++
			if i >= len(args) {
				return "", "", "", "", false, "", errors.New("-output 缺少参数")
			}
			outputPath = args[i]
		case "-debug":
			debugMode = true
		case "-debug-output":
			i++
			if i >= len(args) {
				return "", "", "", "", false, "", errors.New("-debug-output 缺少参数")
			}
			debugOutputPath = args[i]
		default:
			return "", "", "", "", false, "", fmt.Errorf("未知参数: %s", args[i])
		}
	}

	if inputPath == "" || startMonth == "" || endMonth == "" {
		return "", "", "", "", false, "", errors.New("必须提供 -input -start -end 参数")
	}

	if outputPath == "" {
		outputPath = defaultOutputPath(inputPath, startMonth, endMonth)
	}
	if debugMode && debugOutputPath == "" {
		debugOutputPath = defaultDebugOutputPath(outputPath)
	}

	return inputPath, startMonth, endMonth, outputPath, debugMode, debugOutputPath, nil
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
	filename := fmt.Sprintf("%s_%s_to_%s.csv", base, startMonth, endMonth)
	return filepath.Join(filepath.Dir(inputPath), filename)
}

func defaultDebugOutputPath(outputPath string) string {
	base := strings.TrimSuffix(outputPath, filepath.Ext(outputPath))
	return base + "_debug.csv"
}

func readDemandRecords(xlsxPath string, start, end time.Time, debugMode bool) ([]demandRecord, error) {
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

	headerRowIdx, titleCol, startCol, endCol, err := detectHeaderRow(rows)
	if err != nil {
		return nil, err
	}

	mergeRefMap, err := buildMergeRefMap(file, sheet)
	if err != nil {
		return nil, err
	}

	demandMap := make(map[string]*demandAgg)
	for rowIdx := headerRowIdx + 2; rowIdx <= len(rows); rowIdx++ {
		titleCell, startCell, endCell := cellName(titleCol, rowIdx), cellName(startCol, rowIdx), cellName(endCol, rowIdx)

		title := strings.TrimSpace(readCellValueByMerge(file, sheet, titleCell, mergeRefMap))
		if title == "" {
			continue
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
				Title:    title,
				MinStart: startTime,
				MaxEnd:   endTime,
			}
			continue
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
			Start:       v.MinStart,
			End:         v.MaxEnd,
			SourceSheet: sheet,
			StatMonth:   statMonth,
		})
	}
	return out, nil
}

func detectHeaderRow(rows [][]string) (headerRowIdx, titleCol, startCol, endCol int, err error) {
	maxCheck := 2
	if len(rows) < maxCheck {
		maxCheck = len(rows)
	}
	for i := 0; i < maxCheck; i++ {
		titleCol, startCol, endCol, err = headerIndexes(rows[i])
		if err == nil {
			return i, titleCol, startCol, endCol, nil
		}
	}
	return -1, -1, -1, -1, errors.New("前两行未识别到表头(需求标题/研发开始时间/测试完成时间)")
}

func headerIndexes(header []string) (titleCol, startCol, endCol int, err error) {
	titleCol, startCol, endCol = -1, -1, -1
	for i, h := range header {
		normalized := normalizeHeader(h)
		switch {
		case isTitleHeader(normalized):
			titleCol = i + 1
		case isStartHeader(normalized):
			startCol = i + 1
		case isEndHeader(normalized):
			endCol = i + 1
		}
	}
	if titleCol == -1 || startCol == -1 || endCol == -1 {
		return -1, -1, -1, errors.New("表头缺少 需求标题/研发开始时间/测试完成时间")
	}
	return titleCol, startCol, endCol, nil
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
	for _, rec := range records {
		month := statMonthOf(rec)
		days := inclusiveWorkdayDiff(rec.Start, rec.End)
		if days < 0 {
			continue
		}
		byMonth[month] = append(byMonth[month], days)
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

func writeCSV(outputPath string, stats []monthStat) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建输出文件失败: %w", err)
	}
	defer file.Close()
	if _, err := file.WriteString("\uFEFF"); err != nil {
		return fmt.Errorf("写入 UTF-8 BOM 失败: %w", err)
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"统计月份", "需求交付数量", "需求交付时长P80"}
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, s := range stats {
		row := []string{
			s.Month,
			fmt.Sprintf("%d", s.DeliveryCount),
			fmt.Sprintf("%.2f", s.DeliveryP80Day),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return writer.Error()
}

func writeDebugCSV(outputPath string, items []debugItem) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建调试输出文件失败: %w", err)
	}
	defer file.Close()
	if _, err := file.WriteString("\uFEFF"); err != nil {
		return fmt.Errorf("写入 UTF-8 BOM 失败: %w", err)
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"统计月份", "来源Sheet", "需求标题", "需求开始时间", "需求结束时间", "交付时长"}
	if err := writer.Write(header); err != nil {
		return err
	}
	for _, item := range items {
		row := []string{
			item.Month,
			item.SourceSheet,
			item.Title,
			item.Start.Format(timeLayoutYMD),
			item.End.Format(timeLayoutYMD),
			fmt.Sprintf("%.2f", item.Days),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	return writer.Error()
}

func exitWithErr(err error) {
	fmt.Fprintf(os.Stderr, "错误: %v\n", err)
	os.Exit(1)
}
