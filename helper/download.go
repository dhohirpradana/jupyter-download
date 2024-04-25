package helper

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gopkg.in/validator.v2"
	"jupyter-download/entity"
	"os"
	"path/filepath"
	"strings"
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
		_ = removeFiles(randomUUIDString, temp)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	defer func() {
		_ = removeFiles(randomUUIDString, temp)
	}()

	c.Set(fiber.HeaderContentType, "application/zip")
	return c.SendFile(temp + ".zip")
}

func (h DownloadHandler) FilesDownload(c *fiber.Ctx) (err error) {
	var dlFiles entity.DlFiles

	if err := c.BodyParser(&dlFiles); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	if err := validator.Validate(dlFiles); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	randomUUID := uuid.New()
	randomUUIDString := randomUUID.String()
	username := dlFiles.Username
	podName := "jupyter-" + username + "-0"
	temp := randomUUIDString

	for _, file := range dlFiles.Files {
		fmt.Println("file:", file)
		trimFile := strings.TrimPrefix(file, "/home/jupyter")
		args := []string{"-n", "sapujagad2", "--kubeconfig", "kubeconfig", "cp", podName + ":" + file, temp + "/" + trimFile}
		err = Exec("kubectl", args)
		if err != nil {
			_ = removeFiles(randomUUIDString, temp)
			return fiber.NewError(fiber.StatusUnprocessableEntity, file+", "+err.Error())
		}

		// Check file
		destFilePath := temp + "/" + trimFile
		if _, err := os.Stat(destFilePath); err == nil {
			fmt.Println(file + ": OK")
		} else if os.IsNotExist(err) {
			fmt.Println(file + ": No such file or directory")
			_ = removeFiles(randomUUIDString, temp)
			return fiber.NewError(fiber.StatusUnprocessableEntity, file+": No such file or directory")
		} else {
			fmt.Println(file + ": " + err.Error())
			_ = removeFiles(randomUUIDString, temp)
			return fiber.NewError(fiber.StatusUnprocessableEntity, file+": "+err.Error())
		}
	}

	args1 := []string{"-jr", temp + ".zip", temp}
	err = Exec("zip", args1)
	if err != nil {
		_ = removeFiles(randomUUIDString, temp)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	defer func() {
		_ = removeFiles(randomUUIDString, temp)
	}()

	c.Set(fiber.HeaderContentType, "application/zip")
	return c.SendFile(temp + ".zip")
}

func removeFiles(dir, file string) error {
	err := os.RemoveAll(dir)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = os.Remove(file + ".zip")
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}
