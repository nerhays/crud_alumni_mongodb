package service

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"path/filepath"
	"testing"
	"time"

	"crud_alumni/app/model"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- Mock repo ---
type mockFileRepo struct {
	created       *model.File
	files         []model.File
	getByIDResult *model.File
	deleteErr     error
	createErr     error
}

func (m *mockFileRepo) Create(file *model.File) error {
	if m.createErr != nil {
		return m.createErr
	}
	// simulate DB assigning ID
	file.ID = primitive.NewObjectID()
	m.created = file
	return nil
}

func (m *mockFileRepo) GetAll() ([]model.File, error) {
	return m.files, nil
}

func (m *mockFileRepo) GetByUserID(userID string) ([]model.File, error) {
	return m.files, nil
}

func (m *mockFileRepo) GetByID(id string) (*model.File, error) {
	return m.getByIDResult, nil
}

func (m *mockFileRepo) DeleteByID(id string) error {
	return m.deleteErr
}

// --- helpers for multipart ---
func makeMultipart(bodyFieldName, filename, contentType string, content []byte) (string, *bytes.Buffer, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="`+bodyFieldName+`"; filename="`+filename+`"`)
	h.Set("Content-Type", contentType)
	part, err := w.CreatePart(h)
	if err != nil {
		return "", nil, err
	}
	if _, err := part.Write(content); err != nil {
		return "", nil, err
	}
	w.Close()
	return w.FormDataContentType(), &buf, nil
}

// Need textproto import in this file
// add at top: "net/textproto"
//
// (go test will fail if not added; ensure you add the import)

func TestUploadFile_SuccessFoto(t *testing.T) {
	// prepare mock repo
	mock := &mockFileRepo{}
	svc := NewFileService(mock)

	// setup Fiber app with middleware to set locals
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", primitive.NewObjectID().Hex())
		c.Locals("role", "user")
		return c.Next()
	})
	// route uses svc.UploadFile and extracts category from path param
	app.Post("/file/:category", func(c *fiber.Ctx) error {
		return svc.UploadFile(c, c.Params("category"))
	})

	// create small PNG-like bytes (not a real PNG but content-type matters)
	content := []byte{0x89, 0x50, 0x4E, 0x47}
	ct, body, err := makeMultipart("file", "photo.png", "image/png", content)
	if err != nil {
		t.Fatalf("failed to create multipart: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/file/foto", body)
	req.Header.Set("Content-Type", ct)

	// perform request
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 201 {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 201 got %d, body: %s", resp.StatusCode, string(b))
	}

	// ensure repo.Create was called and metadata looks ok
	if mock.created == nil {
		t.Fatalf("expected repo.Create to be called")
	}
	if mock.created.Category != "foto" {
		t.Errorf("expected category foto, got %s", mock.created.Category)
	}
	if mock.created.FileType != "image/png" {
		t.Errorf("expected filetype image/png, got %s", mock.created.FileType)
	}
}

func TestGetAllFiles_Admin(t *testing.T) {
	mock := &mockFileRepo{
		files: []model.File{
			{
				ID:           primitive.NewObjectID(),
				FileName:     "a.txt",
				OriginalName: "a.txt",
				Category:     "sertifikat",
				UploadedAt:   time.Now(),
			},
		},
	}
	svc := NewFileService(mock)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("role", "admin")
		c.Locals("user_id", primitive.NewObjectID().Hex())
		return c.Next()
	})
	app.Get("/file", func(c *fiber.Ctx) error {
		return svc.GetAllFiles(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/file", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var payload struct {
		Count int           `json:"count"`
		Data  []model.File  `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if payload.Count != 1 {
		t.Errorf("expected count 1, got %d", payload.Count)
	}
}

func TestGetFileByID_Unauthorized(t *testing.T) {
	otherUID := primitive.NewObjectID()
	mock := &mockFileRepo{
		getByIDResult: &model.File{
			ID:     primitive.NewObjectID(),
			UserID: otherUID,
		},
	}
	svc := NewFileService(mock)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		// set current user different from file owner
		c.Locals("role", "user")
		c.Locals("user_id", primitive.NewObjectID().Hex())
		return c.Next()
	})
	app.Get("/file/:id", func(c *fiber.Ctx) error {
		return svc.GetFileByID(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/file/"+primitive.NewObjectID().Hex(), nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 403 {
		t.Fatalf("expected 403, got %d", resp.StatusCode)
	}
}

func TestDeleteFile_Success(t *testing.T) {
	uid := primitive.NewObjectID()
	// create temporary file to be deleted
	tmpdir := t.TempDir()
	tmpfile := filepath.Join(tmpdir, "to_delete.txt")
	if err := os.WriteFile(tmpfile, []byte("ok"), 0644); err != nil {
		t.Fatalf("write tmp file: %v", err)
	}

	mock := &mockFileRepo{
		getByIDResult: &model.File{
			ID:       primitive.NewObjectID(),
			UserID:   uid,
			FilePath: tmpfile,
		},
	}
	svc := NewFileService(mock)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("role", "user")
		c.Locals("user_id", uid.Hex()) // same owner
		return c.Next()
	})
	app.Delete("/file/:id", func(c *fiber.Ctx) error {
		return svc.DeleteFile(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/file/"+primitive.NewObjectID().Hex(), nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 200 got %d, body: %s", resp.StatusCode, string(b))
	}
	// ensure file removed
	if _, err := os.Stat(tmpfile); !os.IsNotExist(err) {
		t.Fatalf("expected file to be removed, stat err: %v", err)
	}
}
