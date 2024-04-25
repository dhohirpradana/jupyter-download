package helper

import (
	"github.com/gofiber/fiber/v2"
	"gopkg.in/validator.v2"
	"jupyter-folder-download/entity"
	"os"
	"path/filepath"
)

type DownloadHandler struct {
}

func InitDownload() DownloadHandler {
	return DownloadHandler{}
}

func (h DownloadHandler) FolderDownload(c *fiber.Ctx) (err error) {
	var dlDir entity.DlDir

	if err := c.BodyParser(&dlDir); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	if err := validator.Validate(dlDir); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	username := dlDir.Username
	podName := "jupyter-" + username + "-0"
	dir := dlDir.Dir
	lastDir := filepath.Base(dir)
	args := []string{"-n", "sapujagad2", "--kubeconfig", "kubeconfig", "cp", podName + ":" + dir, lastDir}
	err = Exec("kubectl", args)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	args1 := []string{"-r", lastDir + ".zip", lastDir}
	err = Exec("zip", args1)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	defer func() {
		os.RemoveAll(lastDir)
		os.Remove(lastDir + ".zip")
	}()

	c.Set(fiber.HeaderContentType, "application/zip")
	return c.SendFile(lastDir + ".zip")
}
