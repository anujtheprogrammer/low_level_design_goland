package main

type Document interface {
	Render(data []byte) []byte
}

type PDFDocument struct{}

func (PDFDocument) Render(d []byte) []byte {
	return d
}

type ExcelDocument struct{}

func (ExcelDocument) Render(d []byte) []byte {
	return d
}

type Exporter interface {
	CreateDocument() Document
	Export(data []byte) []byte
}

type PDFExporter struct{}

func (PDFExporter) CreateDocument() Document {
	return PDFDocument{}
}
func (p PDFExporter) Export(d []byte) []byte {
	return p.CreateDocument().Render(d)
}

type ExcelExporter struct{}

func (ExcelExporter) CreateDocument() Document {
	return ExcelDocument{}
}
func (e ExcelExporter) Export(d []byte) []byte {
	return e.CreateDocument().Render(d)
}

// chosen once at account setup - every later call just use it
func RunExport(exp Exporter, data []byte) []byte {
	return exp.Export(data)
}
