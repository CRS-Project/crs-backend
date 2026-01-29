package mypdf

import (
	"bytes"
	"fmt"
	_ "image/png" // Penting untuk import gambar
	"strconv"

	"github.com/xuri/excelize/v2"
)

// GenerateExcel membuat file Excel berdasarkan data request
func GenerateExcel(req []GenerateRequestData) (*bytes.Buffer, string, error) {
	f := excelize.NewFile()

	// Hapus sheet default "Sheet1" nanti setelah kita buat sheet baru
	defaultSheet := "Sheet1"

	for i, r := range req {
		// Buat nama sheet unik, misal "CRS 1", "CRS 2"
		sheetName := fmt.Sprintf("CRS %d", i+1)
		index, err := f.NewSheet(sheetName)
		if err != nil {
			return nil, "", err
		}
		f.SetActiveSheet(index)

		// Set lebar kolom agar proporsional mirip PDF
		setColumnWidths(f, sheetName)

		// Siapkan Style
		styles, err := createExcelStyles(f)
		if err != nil {
			return nil, "", err
		}

		// Gambar Layout
		currentRow := 1
		if err := drawExcelHeader(f, sheetName, &currentRow, styles); err != nil {
			return nil, "", err
		}

		drawExcelPackageInfo(f, sheetName, r.PackageInfoData, &currentRow)

		// Spasi sebelum tabel discipline
		currentRow += 1

		drawExcelDisciplineSection(f, sheetName, r.DisciplineSectionData, &currentRow, styles)

		// Spasi sebelum tabel utama
		currentRow += 2

		drawExcelMainTable(f, sheetName, r.CommentRow, &currentRow, styles)
	}

	// Hapus sheet default jika tidak terpakai
	if len(req) > 0 {
		f.DeleteSheet(defaultSheet)
	}

	filename := "comment_resolution_sheet.xlsx"

	// Write ke buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, "", err
	}

	return &buf, filename, nil
}

// --- Helper Functions ---

// setColumnWidths mengatur lebar kolom (perkiraan konversi mm ke width excel)
func setColumnWidths(f *excelize.File, sheet string) {
	// Mapping lebar kolom berdasarkan visual PDF:
	// A: No, B: Page, C: SME Init, D: SME Comment
	// E: RefDocNo, F: RefDocTitle, G: DocStatus, H: Status, I: CloseOut

	f.SetColWidth(sheet, "A", "A", 6)  // No
	f.SetColWidth(sheet, "B", "B", 10) // Page
	f.SetColWidth(sheet, "C", "C", 12) // SME Initial
	f.SetColWidth(sheet, "D", "D", 40) // SME Comment (Lebar agar muat banyak teks)
	f.SetColWidth(sheet, "E", "E", 25) // Ref Doc No
	f.SetColWidth(sheet, "F", "F", 30) // Ref Doc Title
	f.SetColWidth(sheet, "G", "G", 15) // Doc Status
	f.SetColWidth(sheet, "H", "H", 10) // Status
	f.SetColWidth(sheet, "I", "I", 30) // SME Close Out
}

type excelStyles struct {
	Title        int
	HeaderGray   int
	HeaderYellow int
	BodyText     int
	LabelBold    int
	BorderBox    int
}

func createExcelStyles(f *excelize.File) (*excelStyles, error) {
	// Definisi Border Hitam Tipis
	border := []excelize.Border{
		{Type: "left", Color: "000000", Style: 1},
		{Type: "top", Color: "000000", Style: 1},
		{Type: "bottom", Color: "000000", Style: 1},
		{Type: "right", Color: "000000", Style: 1},
	}

	// Style Judul Utama
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 16, Family: "Arial"},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})

	// Style Header Abu-abu (Discipline & Table Header Kiri)
	headerGray, _ := f.NewStyle(&excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#E7E6E6"}, Pattern: 1},
		Font:      &excelize.Font{Bold: true, Size: 9, Family: "Arial"},
		Border:    border,
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "top", WrapText: true},
	})

	// Style Header Kuning (Table Header Kanan)
	headerYellow, _ := f.NewStyle(&excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#FFFF00"}, Pattern: 1},
		Font:      &excelize.Font{Bold: true, Size: 9, Family: "Arial"},
		Border:    border,
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "top", WrapText: true},
	})

	// Style Isi Tabel (Wrap Text)
	bodyText, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 9, Family: "Arial"},
		Border:    border,
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "top", WrapText: true},
	})

	// Style Label Bold (Tanpa Border)
	labelBold, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 9, Family: "Arial"},
	})

	// Style Kotak Kosong dengan Border (untuk Consolidator dll)
	borderBox, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 9, Family: "Arial"},
		Border:    border,
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "top", WrapText: true},
	})

	return &excelStyles{
		Title:        titleStyle,
		HeaderGray:   headerGray,
		HeaderYellow: headerYellow,
		BodyText:     bodyText,
		LabelBold:    labelBold,
		BorderBox:    borderBox,
	}, nil
}

func drawExcelHeader(f *excelize.File, sheet string, row *int, s *excelStyles) error {
	// Masukkan Logo
	// Catatan: Pastikan path gambar valid. Jika error, gambar tidak muncul tapi file tetap tergenerate.
	if err := f.AddPicture(sheet, fmt.Sprintf("A%d", *row), "./assets/image/Logo-CRS.png", &excelize.GraphicOptions{
		ScaleX: 0.2,
		ScaleY: 0.2,
	}); err != nil {
		fmt.Println("Warning: Logo not found or error loading:", err)
	}

	// Judul "COMMENT RESOLUTION SHEET"
	// Merge cell dari D sampai I agar center visualnya pas
	f.MergeCell(sheet, fmt.Sprintf("D%d", *row), fmt.Sprintf("I%d", *row+1))
	f.SetCellValue(sheet, fmt.Sprintf("D%d", *row), "COMMENT RESOLUTION SHEET")
	f.SetCellStyle(sheet, fmt.Sprintf("D%d", *row), fmt.Sprintf("I%d", *row+1), s.Title)

	*row += 4 // Turun baris setelah header
	return nil
}

func drawExcelPackageInfo(f *excelize.File, sheet string, data PackageInfoData, row *int) {
	startRow := *row

	// Kiri
	f.SetCellValue(sheet, fmt.Sprintf("A%d", startRow), "Package")
	f.SetCellValue(sheet, fmt.Sprintf("C%d", startRow), ": "+data.Package)

	f.SetCellValue(sheet, fmt.Sprintf("A%d", startRow+1), "FEED Contractor")
	f.SetCellValue(sheet, fmt.Sprintf("C%d", startRow+1), ": "+data.ContractorInitial)

	// Kanan
	f.SetCellValue(sheet, fmt.Sprintf("G%d", startRow), "Inc. Transmittal")
	f.SetCellValue(sheet, fmt.Sprintf("H%d", startRow), ":")

	f.SetCellValue(sheet, fmt.Sprintf("G%d", startRow+1), "Out. Transmittal")
	f.SetCellValue(sheet, fmt.Sprintf("H%d", startRow+1), ":")

	f.SetCellValue(sheet, fmt.Sprintf("G%d", startRow+2), "Out. Transmittal Date")
	f.SetCellValue(sheet, fmt.Sprintf("H%d", startRow+2), ":")

	*row += 4
}

func drawExcelDisciplineSection(f *excelize.File, sheet string, data DisciplineSectionData, row *int, s *excelStyles) {
	// Header Tabel Kecil
	f.SetCellValue(sheet, fmt.Sprintf("A%d", *row), "Discipline")
	f.SetCellStyle(sheet, fmt.Sprintf("A%d", *row), fmt.Sprintf("B%d", *row), s.HeaderGray)
	f.MergeCell(sheet, fmt.Sprintf("A%d", *row), fmt.Sprintf("B%d", *row)) // Merge A-B untuk Discipline

	f.SetCellValue(sheet, fmt.Sprintf("C%d", *row), "Consolidator")
	f.SetCellStyle(sheet, fmt.Sprintf("C%d", *row), fmt.Sprintf("D%d", *row), s.HeaderGray)
	f.MergeCell(sheet, fmt.Sprintf("C%d", *row), fmt.Sprintf("D%d", *row))

	*row++

	// Isi Tabel Kecil
	f.SetCellValue(sheet, fmt.Sprintf("A%d", *row), data.Discipline)
	f.SetCellStyle(sheet, fmt.Sprintf("A%d", *row), fmt.Sprintf("B%d", *row+1), s.BorderBox)
	f.MergeCell(sheet, fmt.Sprintf("A%d", *row), fmt.Sprintf("B%d", *row+1)) // Tinggi 2 baris

	f.SetCellValue(sheet, fmt.Sprintf("C%d", *row), data.Consolidator)
	f.SetCellStyle(sheet, fmt.Sprintf("C%d", *row), fmt.Sprintf("D%d", *row+1), s.BorderBox)
	f.MergeCell(sheet, fmt.Sprintf("C%d", *row), fmt.Sprintf("D%d", *row+1))

	*row += 3

	// Catatan kaki kecil
	f.SetCellValue(sheet, fmt.Sprintf("A%d", *row), "*Please manually sort page number in ascending order")
	f.SetCellStyle(sheet, fmt.Sprintf("A%d", *row), fmt.Sprintf("A%d", *row),
		mustNewStyle(f, &excelize.Style{Font: &excelize.Font{Italic: true, Size: 8, Family: "Arial"}}))
}

func drawExcelMainTable(f *excelize.File, sheet string, rows []CommentRow, row *int, s *excelStyles) {
	headers := []string{
		"No.", "Page *", "SME Initial", "SME\nComment",
		"Ref. Document No.", "Ref. Document Title",
		"Doc. Status", "Status", "SME Close Out\nComments",
	}

	colNames := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I"}

	// Draw Headers
	for i, header := range headers {
		cell := colNames[i] + strconv.Itoa(*row)
		f.SetCellValue(sheet, cell, header)

		// 4 Kolom pertama Abu-abu, sisanya Kuning (Sesuai logika PDF: i >= 4)
		if i >= 4 {
			f.SetCellStyle(sheet, cell, cell, s.HeaderYellow)
		} else {
			f.SetCellStyle(sheet, cell, cell, s.HeaderGray)
		}
	}
	*row++

	// Draw Data Rows
	for _, dataRow := range rows {
		vals := []string{
			dataRow.No,
			dataRow.Page,
			dataRow.SMEInitial,
			dataRow.SMEComment,
			dataRow.RefDocNo,
			dataRow.RefDocTitle,
			dataRow.DocStatus,
			dataRow.Status,
			dataRow.SMECloseComment,
		}

		// Cari tinggi baris maksimum berdasarkan panjang teks (manual estimation tidak sempurna di excel,
		// tapi wrap text akan menanganinya secara visual)

		for i, val := range vals {
			cell := colNames[i] + strconv.Itoa(*row)
			f.SetCellValue(sheet, cell, val)
			f.SetCellStyle(sheet, cell, cell, s.BodyText)
		}
		*row++
	}
}

// Helper kecil untuk error handling style inline
func mustNewStyle(f *excelize.File, s *excelize.Style) int {
	id, _ := f.NewStyle(s)
	return id
}
