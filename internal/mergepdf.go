package mergepdf

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

func mergeOddEvenPdfs(oddPdfPath string, evenPdfPath string, outputPdfPath string) error {

	conf := model.NewDefaultConfiguration()

	err := mergeCreateZipFile(oddPdfPath, evenPdfPath, outputPdfPath, conf)
	if err != nil {
		return err
	}

	return nil
}

func mergeCreateZip(rs1, rs2 io.ReadSeeker, w io.Writer, conf *model.Configuration) error {

	if rs1 == nil {
		return errors.New("pdfcpu: MergeCreateZip: missing rs1")
	}
	if rs2 == nil {
		return errors.New("pdfcpu: MergeCreateZip: missing rs2")
	}
	if w == nil {
		return errors.New("pdfcpu: MergeCreateZip: missing w")
	}

	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.MERGECREATEZIP
	conf.ValidationMode = model.ValidationRelaxed

	ctxDest, err := api.ReadAndValidate(rs1, conf)
	if err != nil {
		return err
	}
	if ctxDest.XRefTable.Version() == model.V20 {
		return pdfcpu.ErrUnsupportedVersion
	}
	ctxDest.EnsureVersionForWriting()

	if _, err = pdfcpu.RemoveBookmarks(ctxDest); err != nil {
		return err
	}

	ctxSrc, err := api.ReadAndValidate(rs2, conf)
	if err != nil {
		return err
	}
	if ctxSrc.XRefTable.Version() == model.V20 {
		return pdfcpu.ErrUnsupportedVersion
	}

	ctxSrc2, err := reversePages(ctxSrc)
	if err != nil {
		return err
	}

	if err := pdfcpu.MergeXRefTables("", ctxSrc2, ctxDest, true, false); err != nil {
		return err
	}

	if err := api.OptimizeContext(ctxDest); err != nil {
		return err
	}

	return api.WriteContext(ctxDest, w)
}

// MergeCreateZipFile zips inFile1 and inFile2 into outFile.
func mergeCreateZipFile(inFile1, inFile2, outFile string, conf *model.Configuration) (err error) {
	f1, err := os.Open(inFile1)
	if err != nil {
		return err
	}

	f2, err := os.Open(inFile2)
	if err != nil {
		return err
	}

	f, err := os.Create(outFile)
	if err != nil {
		return err
	}

	defer func() {
		cerr := f.Close()
		if err == nil {
			err = cerr
		}
	}()

	err = mergeCreateZip(f1, f2, f, conf)
	return err
}

func reversePages(ctxSrc *model.Context) (ctxDest *model.Context, err error) {

	cnt := ctxSrc.PageCount
	conf := model.NewDefaultConfiguration()
	ctxDestTmp, err := pdfcpu.CreateContextWithXRefTable(conf, types.PaperSize["A4"])
	if err != nil {
		return nil, err
	}

	pages := make([]int, cnt)
	i := 1
	j := cnt - 1
	for j >= 0 {
		pages[j] = i
		i++
		j--
	}

	err = pdfcpu.AddPages(ctxSrc, ctxDestTmp, pages, false)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	writerBuf := bufio.NewWriter(&buf)
	err = api.WriteContext(ctxDestTmp, writerBuf)
	if err != nil {
		return nil, err
	}

	readerBuf := bytes.NewReader(buf.Bytes())
	ctxDest, err = api.ReadAndValidate(readerBuf, conf)
	if err != nil {
		return nil, err
	}

	return ctxDest, nil
}
