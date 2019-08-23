package api

import (
	"context"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

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

// CleanProducedImages is cleaning the produced Image
func CleanProducedImages(ctx context.Context, cli *client.Client) error {
	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return err
	}
	for _, image := range images {
		if len(image.RepoTags) >= 1 && image.RepoTags[0] == "previs:latest" {
			_, err := cli.ImageRemove(ctx, image.ID, types.ImageRemoveOptions{
				Force:         true,
				PruneChildren: true,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// CleanProducedContainer is cleaning the produced container from the image
func CleanProducedContainer(ctx context.Context, cli *client.Client, idContainer string) error {
	err := cli.ContainerRemove(ctx, idContainer, types.ContainerRemoveOptions{
		Force: true,
	})
	if err != nil {
		return err
	}
	return nil
}

// CleanAll is the sum of the three Cleaning function
func CleanAll(ctx context.Context, cli *client.Client, idContainer string, imgDocker string, pathDockerImage string) error {
	err := CleanProducedImages(ctx, cli)
	if err != nil {
		return err
	}
	err = CleanProducedContainer(ctx, cli, idContainer)
	if err != nil {
		return err
	}
	err = CleanUnusedDockerfile(pathDockerImage, imgDocker)
	if err != nil {
		return err
	}
	return nil
}
