package image

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/dto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestSyncGulie2KImageFormStripsResolutionFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	require.NoError(t, writer.WriteField("prompt", "cat"))
	require.NoError(t, writer.WriteField("size", "3840x2160"))
	require.NoError(t, writer.WriteField("quality", "high"))
	require.NoError(t, writer.WriteField("image_size", "4K"))
	require.NoError(t, writer.Close())

	req := httptest.NewRequest(http.MethodPost, "/v1/images/edits", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	require.NoError(t, req.ParseMultipartForm(1<<20))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	syncGulie2KImageForm(c, &dto.ImageRequest{Size: "16:9"})

	require.Equal(t, []string{"16:9"}, c.Request.MultipartForm.Value["size"])
	require.Nil(t, c.Request.MultipartForm.Value["quality"])
	require.Nil(t, c.Request.MultipartForm.Value["image_size"])
	require.Equal(t, "16:9", c.Request.PostForm.Get("size"))
	require.Empty(t, c.Request.PostForm.Get("quality"))
}
