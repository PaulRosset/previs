package api

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"

	tm "github.com/buger/goterm"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"golang.org/x/net/context"
)

type building struct {
	Status   string `json:"status"`
	Stream   string `json:"stream"`
	Progress string `json:"progress"`
}

func showLoaderWhenBuildImage(Body io.ReadCloser) {
	tm.Clear()
	decode := json.NewDecoder(Body)
	for decode.More() {
		tm.MoveCursor(1, 1)
		var Build building
		err := decode.Decode(&Build)
		if err != nil {
			break
		}
		if Build.Status != "" {
			tm.Printf("Building Image\nStatus: %s\nProgress: +%v", Build.Status, Build.Progress)
		} else if len(Build.Stream) > 1 && Build.Stream[:4] == "Step" {
			tm.Printf("%+v", Build.Stream)
		}
		tm.Flush()
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
	showLoaderWhenBuildImage(builderResp.Body)
	return nil
}

func array_contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func startContainer(ctx context.Context, cli *client.Client, imgDocker string, pathDockerImage string, envVar []string, services []string) (string, error) {
	hostconfig := &container.HostConfig{}
	if array_contains(services, "docker") {
		hostconfig = &container.HostConfig{
			Binds: []string{"/var/run/docker.sock:/var/run/docker.sock"},
		}
	}
	containerName := randSeq(10) + "-previs"
	respContainerCreater, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "previs",
		Tty:   true,
		Env:   envVar,
	}, hostconfig, nil, containerName)
	if err != nil {
		return fmt.Sprintf("ContainerCreate failed with name '%s'", containerName), err
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
		readerOutputLog, errOnLogs := cli.ContainerLogs(ctx, respContainerCreater.ID, types.ContainerLogsOptions{
			ShowStderr: true,
			ShowStdout: true,
		})
		if errOnLogs != nil {
			return "", errOnLogs
		}
		defer readerOutputLog.Close()
		io.Copy(os.Stdout, readerOutputLog)
		err = CleanAll(ctx, cli, respContainerCreater.ID, imgDocker, pathDockerImage)
		if err != nil {
			return "", err
		}
		if status.StatusCode == 0 {
			fmt.Printf("\nYour build has passed\n\n")
			return respContainerCreater.ID, nil
		}
		fmt.Printf("\nYour build has failed for the container %s\n\nLOGS:\n", respContainerCreater.ID)
	}
	return respContainerCreater.ID, nil
}

// Start the pipeline of test: Build,Launch,Clean
func Start(imgDocker string, pathDockerImage string, envVar []string, services []string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithVersion("1.38"))
	if err != nil {
		return err
	}
	err = buildImage(ctx, cli, imgDocker, pathDockerImage)
	if err != nil {
		errOnCleanImages := CleanProducedImages(ctx, cli)
		if errOnCleanImages != nil {
			return errOnCleanImages
		}
		return err
	}
	idContainer, err := startContainer(ctx, cli, imgDocker, pathDockerImage, envVar, services)
	if err != nil {
		errOnCleanCOntainer := CleanProducedContainer(ctx, cli, idContainer)
		if errOnCleanCOntainer != nil {
			return errOnCleanCOntainer
		}
		return err
	}
	return nil
}
