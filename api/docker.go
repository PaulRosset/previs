package api

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func getBuildContext(imgDocker string) (io.Reader, error) {
	buildContextReader, errOnOpen := os.Open(imgDocker)
	if errOnOpen != nil {
		return nil, errOnOpen
	}
	return buildContextReader, nil
}

func createTarFileBeforeBuild(buildContext io.Reader, imgDocker string, pathDockerImage string) (io.Reader, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	contentDockerFileReader, err := getBuildContext(pathDockerImage + imgDocker)
	if err != nil {
		return nil, err
	}

	readDockerFile, err := ioutil.ReadAll(contentDockerFileReader)
	if err != nil {
		return nil, err
	}

	tarHeader := &tar.Header{
		Name: imgDocker,
		Size: int64(len(readDockerFile)),
	}
	err = tw.WriteHeader(tarHeader)
	if err != nil {
		return nil, err
	}
	_, err = tw.Write(readDockerFile)
	if err != nil {
		return nil, err
	}
	dockerFileTarReader := bytes.NewReader(buf.Bytes())
	return dockerFileTarReader, nil
}

func buildImage(cli *client.Client, imgDocker string, pathDockerImage string, ctx context.Context) error {
	buildContextReader, err := getBuildContext(imgDocker)
	if err != nil {
		return err
	}
	dockerFileTarReader, err := createTarFileBeforeBuild(buildContextReader, imgDocker, pathDockerImage)
	if err != nil {
		return err
	}
	builderResp, err := cli.ImageBuild(
		ctx,
		dockerFileTarReader,
		types.ImageBuildOptions{
			Context:    dockerFileTarReader,
			Dockerfile: imgDocker,
			Remove:     true,
			Tags:       []string{"previs"},
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

func startContainer(cli *client.Client, ctx context.Context) (string, error) {
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
		if status.StatusCode != 2 {
			fmt.Println("Your tests have passed")
			return respContainerCreater.ID, nil
		}
		fmt.Printf("Your tests have failed for the container %s\nLOGS:\n", respContainerCreater.ID)
		readerOutputLog, errOnLogs := cli.ContainerLogs(ctx, respContainerCreater.ID, types.ContainerLogsOptions{
			ShowStderr: true,
			ShowStdout: true,
		})
		if errOnLogs != nil {
			return "", errOnLogs
		}
		io.Copy(os.Stdout, readerOutputLog)
		err = cleaningImagesContainer(cli, ctx, respContainerCreater.ID)
		if err != nil {
			return "", err
		}
	}
	return respContainerCreater.ID, nil
}

func cleaningImagesContainer(cli *client.Client, ctx context.Context, containerID string) error {
	//cli.ImageRemove
	//cli.ContainerRemove
}

func Start(imgDocker string, pathDockerImage string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithVersion("1.38"))
	if err != nil {
		panic(err)
	}
	// err = buildImage(cli, imgDocker, pathDockerImage, ctx)
	// if err != nil {
	// 	panic(err)
	// }
	containerID, err := startContainer(cli, ctx)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
