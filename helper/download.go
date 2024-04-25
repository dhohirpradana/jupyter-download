package helper

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gopkg.in/validator.v2"
	"jupyter-download/entity"
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

	randomUUID := uuid.New()
	randomUUIDString := randomUUID.String()
	username := dlDir.Username
	podName := "jupyter-" + username + "-0"
	dir := dlDir.Dir
	lastDir := filepath.Base(dir)
	temp := randomUUIDString + "/" + lastDir
	args := []string{"-n", "sapujagad2", "--kubeconfig", "kubeconfig", "cp", podName + ":" + dir, temp}
	err = Exec("kubectl", args)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	args1 := []string{"-jr", temp + ".zip", temp}
	err = Exec("zip", args1)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	defer func() {
		os.RemoveAll(randomUUIDString)
		os.Remove(temp + ".zip")
	}()

	c.Set(fiber.HeaderContentType, "application/zip")
	return c.SendFile(temp + ".zip")
}
