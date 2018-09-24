package api

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

// IsDockerInstall is verifying that docker deamon is running on the machine.
func IsDockerInstall() {
	cmd := exec.Command("docker", "--version")
	_, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "You have to install docker for using Previs\nFor a quick install:\n> curl -fsSL get.docker.com -o get-docker.sh\n> sh get-docker.sh\n")
		os.Exit(2)
	}
}

func buildImage(ctx context.Context, cli *client.Client, imgDocker string, pathDockerImage string) error {
	dockerFileTarReader, err := archive.TarWithOptions(pathDockerImage, &archive.TarOptions{})
	if err != nil {
		return err
	}
	builderResp, err := cli.ImageBuild(
		ctx,
		dockerFileTarReader,
		types.ImageBuildOptions{
			Context:     dockerFileTarReader,
			Dockerfile:  imgDocker,
			Remove:      true,
			ForceRemove: true,
			Tags:        []string{"previs"},
		})
	if err != nil {
		return err
	}

	defer builderResp.Body.Close()
	_, err = io.Copy(os.Stdout, builderResp.Body)
	if err != nil {
		return err
	}
	return nil
}

func startContainer(ctx context.Context, cli *client.Client, imgDocker string, pathDockerImage string) (string, error) {
	respContainerCreater, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "previs",
		Tty:   true,
	}, nil, nil, "previs")
	if err != nil {
		return "", err
	}
	err = cli.ContainerStart(ctx, respContainerCreater.ID, types.ContainerStartOptions{})
	if err != nil {
		return "", err
	}
	statusChan, errChan := cli.ContainerWait(ctx, respContainerCreater.ID, container.WaitConditionNextExit)
	select {
	case errOnWaiting := <-errChan:
		return "", errOnWaiting
	case status := <-statusChan:
		if status.StatusCode == 0 {
			fmt.Println("\nYour tests have passed\n")
			return respContainerCreater.ID, nil
		}
		fmt.Printf("\nYour tests have failed for the container %s\n\nLOGS:\n", respContainerCreater.ID)
		readerOutputLog, errOnLogs := cli.ContainerLogs(ctx, respContainerCreater.ID, types.ContainerLogsOptions{
			ShowStderr: true,
			ShowStdout: true,
		})
		if errOnLogs != nil {
			return "", errOnLogs
		}
		defer readerOutputLog.Close()
		io.Copy(os.Stdout, readerOutputLog)
		err = cleaningImagesContainer(ctx, cli, respContainerCreater.ID, imgDocker, pathDockerImage)
		if err != nil {
			return "", err
		}
	}
	return respContainerCreater.ID, nil
}

func cleaningImagesContainer(ctx context.Context, cli *client.Client, idContainer string, imgDocker string, pathDockerImage string) error {
	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return err
	}
	for _, image := range images {
		if image.RepoTags[0] == "previs:latest" {
			_, err := cli.ImageRemove(ctx, image.ID, types.ImageRemoveOptions{
				Force:         true,
				PruneChildren: true,
			})
			if err != nil {
				return err
			}
		}
	}
	err = cli.ContainerRemove(ctx, idContainer, types.ContainerRemoveOptions{
		Force: true,
	})
	if err != nil {
		return err
	}
	if _, err := os.Stat(pathDockerImage + "/" + imgDocker); err == nil {
		err = os.Remove(pathDockerImage + "/" + imgDocker)
		if err != nil {
			return err
		}
	}
	return nil
}

// Start the pipeline of test: Build,Launch,Clean
func Start(imgDocker string, pathDockerImage string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithVersion("1.38"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error encountered: %+v\n", err)
		os.Exit(2)
	}
	err = buildImage(ctx, cli, imgDocker, pathDockerImage)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error encountered: %+v\n", err)
		os.Exit(2)
	}
	_, err = startContainer(ctx, cli, imgDocker, pathDockerImage)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error encountered: %+v\n", err)
		os.Exit(2)
	}
}
