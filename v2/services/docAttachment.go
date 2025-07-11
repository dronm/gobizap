package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dronm/ds/pgds"
	"github.com/dronm/session"

	"github.com/dronm/gobizap/v2/models"
)

const (
	CACHE_DIR = "CACHE"
)

type DocAttachmentService struct {
	DB      *pgds.PgProvider
	Session session.Session
}

func NewDocAttachmentService(db *pgds.PgProvider, sess session.Session) *DocAttachmentService {
	return &DocAttachmentService{DB: db, Session: sess}
}

// file ID is unique inside ref
func (s *DocAttachmentService) GetAttachmentCacheFileName(baseDir string, refDataType string, refID int, fileID string) string {
	return filepath.Join(baseDir, CACHE_DIR, GetMd5(fmt.Sprintf("att_%s%d_%s", refDataType, refID, fileID)))
}

func (s *DocAttachmentService) GetPreviewCacheFileName(baseDir string, refDataType string, refID int, fileID string) string {
	return filepath.Join(baseDir, CACHE_DIR, GetMd5(fmt.Sprintf("prev_%s%d_%s", refDataType, refID, fileID)))
}

func (s *DocAttachmentService) runCMD(progName, commands, previewName string, toPDF bool) error {
	cmd_args := strings.Split(commands, " ")
	cmd := exec.Command(progName, cmd_args...)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error converting doc to image: %v, params: %s %s", err, progName, commands)
	}

	if toPDF {
		var thbn_n string
		if FileExists(previewName + "-1.jpg") {
			thbn_n = previewName + "-1.jpg"
		} else if FileExists(previewName + "-01.jpg") {
			thbn_n = previewName + "-01.jpg"
		} else if FileExists(previewName + "-001.jpg") {
			thbn_n = previewName + "-001.jpg"
		}
		// thbn_n -->> previewName
		os.Rename(thbn_n, previewName)
	}

	return nil
}

// realName for mime type!!!
// attName - attachment name
// pName - preview name
// realName
func (s *DocAttachmentService) GenThumbnail(attName, pName, realName string) error {
	var fExt string
	f_parts := strings.Split(realName, ".")
	if len(f_parts) > 0 {
		fExt = strings.ToLower(f_parts[len(f_parts)-1])
	}

	pdftoppm_fmt := "-l 1 -scale-to 300 -jpeg %s %s" //-q no comment or errors

	var cmd_name string
	var cmd_s string
	var pdf bool
	if fExt == "doc" || fExt == "docx" || fExt == "xls" || fExt == "xlsx" || fExt == "odt" || fExt == "ods" {
		// openoffice first to pdf

		//export HOME=CACHE && /usr/lib/libreoffice/program/./soffice --headless --convert-to pdf --outdir CACHE CACHE %s
		if err := s.runCMD("soffice",
			fmt.Sprintf("--headless --convert-to pdf:writer_pdf_Export --outdir CACHE %s", attName),
			pName, true); err != nil {
			return err
		}
		// got full pdf attName.pdf
		// pdf to image
		if err := s.runCMD("pdftoppm", fmt.Sprintf(pdftoppm_fmt, attName+".pdf", pName), pName, true); err != nil {
			return err
		}
		os.Remove(attName + ".pdf") // remove temp full pdf file
		return nil

	} else if fExt == "pdf" {
		pdf = true
		cmd_name = "pdftoppm"
		cmd_s = fmt.Sprintf(pdftoppm_fmt, attName, pName)

	} else {
		cmd_name = "convert"
		cmd_s = fmt.Sprintf("-define jpeg:size=500x180 %s -auto-orient -thumbnail 250x100 -unsharp 0x.5 %s", attName, pName)
	}

	return s.runCMD(cmd_name, cmd_s, pName, pdf)
}

func (s *DocAttachmentService) GenAttachmentThumbnail(baseDir string, refDataType string, refID int, fileInfo *models.DocAttachmentContentInfo, attBuf io.Reader) ([]byte, error) {
	att_n := s.GetAttachmentCacheFileName(baseDir, refDataType, refID, fileInfo.ID)
	file_att, err := os.OpenFile(att_n, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return []byte{}, fmt.Errorf("GenAttachmentThumbnail os.OpenFile() failed: %v", err)
	}
	defer file_att.Close()
	_, err = io.Copy(file_att, attBuf)
	if err != nil {
		return []byte{}, fmt.Errorf("GenAttachmentThumbnail io.Copy() failed: %v", err)
	}

	preview_fn := s.GetPreviewCacheFileName(baseDir, refDataType, refID, fileInfo.ID)
	if err := s.GenThumbnail(att_n, preview_fn, fileInfo.Name); err != nil {
		return []byte{}, fmt.Errorf("GenAttachmentThumbnail GenThumbnail() failed: %v", err)
	}
	defer os.Remove(preview_fn)

	var preview_bt []byte
	preview_bt, err = os.ReadFile(preview_fn)
	if err != nil {
		return []byte{}, fmt.Errorf("GenAttachmentThumbnail os.ReadFile() failed: %v", err)
	}
	return preview_bt, nil
}

func (s *DocAttachmentService) AddFileThumbnailToDb(baseDir string, file io.Reader, fileInfo *models.DocAttachmentContentInfo, ref *models.Ref) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return []byte{}, fmt.Errorf("AddFileThumbnailToDb() failed: %v", err)
	}

	fileCont := buf.Bytes()
	fileInfo.Size = int64(buf.Len())
	// thumbnail
	preview_bt, err := s.GenAttachmentThumbnail(baseDir, *ref.DataType, ref.Keys.ID, fileInfo, buf)
	if err != nil {
		return []byte{}, fmt.Errorf("AddFileThumbnailToDb  s.GenAttachmentThumbnail() failed: %v", err)
	}

	if err := s.StoreAttachment(ref, fileInfo, fileCont, preview_bt); err != nil {
		return []byte{}, err
	}

	return preview_bt, nil
}

func (s *DocAttachmentService) StoreAttachment(ref *models.Ref, fileInfo *models.DocAttachmentContentInfo, fileData []byte, previewData []byte) error {
	poolConn, connID, err := s.DB.GetPrimary()
	if err != nil {
		return fmt.Errorf("GetPrimary() failed: %v", err)
	}
	defer s.DB.Release(poolConn, connID)
	conn := poolConn.Conn()

	fileInfo.Size = int64(len(fileData))
	if _, err := conn.Exec(context.Background(), `BEGIN`); err != nil {
		return fmt.Errorf("StoreAttachment conn.Exec() begin failed: %v", err)
	}

	if _, err := conn.Exec(context.Background(),
		`DELETE FROM attachments
		WHERE ref->>'dataType' = $1 AND (ref->'keys'->>'id')::int = $2 AND content_info->>'id' = $3`,
		ref.DataType, ref.Keys.ID, fileInfo.ID,
	); err != nil {

		conn.Exec(context.Background(), `ROLLBACK`)

		return fmt.Errorf("StoreAttachment conn.Exec() delete failed: %v", err)
	}

	if _, err := conn.Exec(context.Background(),
		`INSERT INTO attachments
		(ref, content_info, content_data, content_preview)
		VALUES ($1, $2, $3, $4)`,
		ref,
		fileInfo,
		fileData,
		previewData,
	); err != nil {

		conn.Exec(context.Background(), `ROLLBACK`)

		return fmt.Errorf("StoreAttachment conn.Exec() insert failed: %v", err)
	}

	if _, err := conn.Exec(context.Background(), `COMMIT`); err != nil {
		return fmt.Errorf("StoreAttachment conn.Exec() commit failed: %v", err)
	}
	return nil
}

func (s *DocAttachmentService) AddFile(ctx context.Context, file multipart.File, docAtt models.DocAttachment) (*models.DocAttachment, error) {
	var err error
	docAtt.ContentPreview, err = s.AddFileThumbnailToDb(".", file, &docAtt.ContentInfo, &docAtt.Ref)
	if err != nil {
		return nil, fmt.Errorf("AddFile s.AddFileThumbnailToDb() failed: %v", err)
	}

	return &docAtt, nil
}

func (s *DocAttachmentService) ClearCache(baseDir string, ref models.Ref, contentID string) error {
	att_n := s.GetAttachmentCacheFileName(baseDir, *ref.DataType, ref.Keys.ID, contentID)
	if FileExists(att_n) {
		if err := os.Remove(att_n); err != nil {
			return fmt.Errorf("ClearCache os.Remove() on att_n failed: %v", err)
		}
	}
	preview_fn := s.GetPreviewCacheFileName(baseDir, *ref.DataType, ref.Keys.ID, contentID)
	if FileExists(preview_fn) {
		if err := os.Remove(preview_fn); err != nil {
			return fmt.Errorf("ClearCache os.Remove() on preview_fn failed: %v", err)
		}
	}
	return nil
}

func (s *DocAttachmentService) DelFile(ctx context.Context, ref models.Ref, contentId string) error {
	poolConn, connID, err := s.DB.GetPrimary()
	if err != nil {
		return fmt.Errorf("GetPrimary() failed: %v", err)
	}
	defer s.DB.Release(poolConn, connID)
	conn := poolConn.Conn()

	if _, err := conn.Exec(context.Background(),
		`DELETE FROM attachments
		WHERE ref->>'dataType' = $1
			AND (ref->'keys'->>'id')::int = $2
			AND content_info->>'id' = $3`,
		ref.DataType,
		ref.Keys.ID,
		contentId,
	); err != nil {
		return fmt.Errorf("conn.Exec() delete failed: %v", err)
	}

	return s.ClearCache(".", ref, contentId)
}

func (s *DocAttachmentService) GetFile(ctx context.Context, ref models.Ref, contentId string) (
	cacheFileName string, attachmentName string, retErr error,
) {
	poolConn, connID, err := s.DB.GetSecondary("")
	if err != nil {
		retErr = fmt.Errorf("GetSecondary() failed: %v", err)
		return
	}
	defer s.DB.Release(poolConn, connID)
	conn := poolConn.Conn()

	var attId int64
	if err := conn.QueryRow(context.Background(),
		`SELECT
			id,
			coalesce(content_info->>'name', '')
		FROM attachments
		WHERE ref->>'dataType' = $1
			AND (ref->'keys'->>'id')::int = $2
			AND content_info->>'id' = $3`,
		ref.DataType,
		ref.Keys.ID,
		contentId,
	).Scan(&attId, &attachmentName); err != nil {
		retErr = fmt.Errorf("conn.QueryRow() select failed: %v", err)
		return
	}

	cacheFileName = s.GetAttachmentCacheFileName(".", *ref.DataType, ref.Keys.ID, contentId)
	if !FileExists(cacheFileName) {
		// no cache, read from db && save
		var fileContent []byte
		if err := conn.QueryRow(context.Background(),
			`SELECT
				content_data
			FROM attachments
			WHERE id = $1`,
			attId,
		).Scan(&fileContent); err != nil {
			retErr = fmt.Errorf("conn.QueryRow() select content_data failed: %v", err)
			return
		}
		cacheFile, err := os.Create(cacheFileName)
		if err != nil {
			retErr = fmt.Errorf("os.Create() failed: %v", err)
			return
		}
		defer cacheFile.Close()
		if _, err := cacheFile.Write(fileContent); err != nil {
			retErr = fmt.Errorf("cacheFile.Write() failed: %v", err)
			return
		}
	}

	return
}
