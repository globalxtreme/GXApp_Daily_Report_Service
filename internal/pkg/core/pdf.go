package core

import (
	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
	"log"
)

type PDF struct {
	PDFG *wkhtml.PDFGenerator
}

func (x *PDF) NewGenerator(layout string, data interface{}) {
	pdfg, err := wkhtml.NewPDFGenerator()
	if err != nil {
		log.Panicf("New Generator invalid: %s", err)
	}

	x.PDFG = pdfg

	buffer := xtremepkg.PDFHTMLTemplate(layout, data)
	page := wkhtml.NewPageReader(&buffer)
	page.DisableExternalLinks.Set(true)

	page.HeaderHTML.Set("internal/pkg/layout/pdf/component/header.html")
	page.FooterHTML.Set("internal/pkg/layout/pdf/component/footer.html")

	x.PDFG.AddPage(page)

	x.PDFG.MarginTop.Set(25)
	x.PDFG.MarginBottom.Set(29)
	x.PDFG.MarginLeft.Set(0)
	x.PDFG.MarginRight.Set(0)
}

func (x *PDF) Save(path string, filename string) error {
	path = xtremepkg.SetStorageAppDir(path)
	xtremepkg.CheckAndCreateDirectory(path)

	err := x.PDFG.Create()
	if err != nil {
		return err
	}

	err = x.PDFG.WriteFile(path + filename)
	if err != nil {
		return err
	}

	return nil
}
