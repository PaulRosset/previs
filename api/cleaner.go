package api

import "os"

// CleanUnusedDockerfile is cleaning Dockerfile file to keep everything clean
func CleanUnusedDockerfile(pathDockerImage string, imgDocker string) error {
	if _, err := os.Stat(pathDockerImage + "/" + imgDocker); err == nil {
		err = os.Remove(pathDockerImage + "/" + imgDocker)
		if err != nil {
			return err
		}
	}
	return nil
}
