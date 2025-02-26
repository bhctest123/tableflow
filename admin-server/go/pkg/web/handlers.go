package web

import (
	"github.com/gin-gonic/gin"
	"github.com/tus/tusd/pkg/filestore"
	"github.com/tus/tusd/pkg/handler"
	"net/http"
	"os"
	"strings"
	"tableflow/go/pkg/file"
	"tableflow/go/pkg/model"
	"tableflow/go/pkg/tf"
	"tableflow/go/pkg/types"
	"tableflow/go/pkg/util"
)

func APIKeyAuthMiddleware(isAuthorized func(c *gin.Context, apiKey string) bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) == 0 {
			tf.Log.Infow("Missing authorization header in request", "host", c.Request.Host, "referer", c.Request.Referer(), "uri", c.Request.RequestURI, "user_agent", c.Request.UserAgent())
			c.AbortWithStatusJSON(http.StatusUnauthorized, types.Res{Err: "unauthorized"})
			return
		}
		authHeader = strings.ToLower(authHeader)
		authHeaderParts := strings.SplitN(authHeader, " ", 2)
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "bearer" || len(authHeaderParts[1]) == 0 {
			tf.Log.Infow("Malformed authorization header in request", "host", c.Request.Host, "referer", c.Request.Referer(), "uri", c.Request.RequestURI, "user_agent", c.Request.UserAgent())
			c.AbortWithStatusJSON(http.StatusUnauthorized, types.Res{Err: "unauthorized"})
			return
		}
		if !isAuthorized(c, authHeaderParts[1]) {
			tf.Log.Infow("Unable to authorize API key", "host", c.Request.Host, "referer", c.Request.Referer(), "uri", c.Request.RequestURI, "user_agent", c.Request.UserAgent())
			c.AbortWithStatusJSON(http.StatusUnauthorized, types.Res{Err: "unauthorized"})
			return
		}
	}
}

// tusFileHandler TODO: Break this out into its own service eventually
func tusFileHandler(uploadAdditionalStorageHandler func(*model.Upload, *os.File) error,
	uploadLimitCheck func(*model.Upload, *os.File) (int, error),
	uploadChunkHandler func(upload *model.Upload, chunk [][]string, isLastChunk bool),
	getAllowedValidateTypes func(string) map[string]bool) *handler.UnroutedHandler {

	store := filestore.FileStore{
		Path: file.TempUploadsDirectory,
	}
	composer := handler.NewStoreComposer()
	store.UseIn(composer)
	fileHandler, err := handler.NewUnroutedHandler(handler.Config{
		BasePath:                "/file-import/v1/files",
		StoreComposer:           composer,
		NotifyCompleteUploads:   true,
		DisableDownload:         true,
		RespectForwardedHeaders: true,
	})
	if err != nil {
		tf.Log.Fatalw("Unable to create tus file upload handler", "error", err)
		return nil
	}
	go func() {
		for {
			event := <-fileHandler.CompleteUploads
			tf.Log.Infow("File upload to disk completed", "tus_id", event.Upload.ID)

			util.SafeGo(func() {
				// TODO: Implement a recover function that updates the upload error
				file.UploadCompleteHandler(event, uploadAdditionalStorageHandler, uploadLimitCheck, uploadChunkHandler, getAllowedValidateTypes)
			}, "tus_id", event.Upload.ID)
		}
	}()
	return fileHandler
}
