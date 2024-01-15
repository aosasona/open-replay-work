package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/viper"
)

const (
	MaxChunkItems    = 5
	RegistryAddress  = "localhost:5000"
	CommandsFileName = "commands.sh"
)

type Service struct {
	Image         string `mapstructure:"image"`
	ContainerName string `mapstructure:"container_name"`
}

type ComposeContent struct {
	Services map[string]Service `mapstructure:"services"`
}

func main() {
	var content ComposeContent

	viper.SetConfigName("docker-compose")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Failed to read config file: %s", err.Error())
	}

	err = viper.Unmarshal(&content)
	if err != nil {
		log.Fatalf("Failed to read config file: %s", err.Error())
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage: ./bin generate|run")
		os.Exit(0)
	}

	cmd := os.Args[1]
	switch cmd {
	case "generate":
		generateFile(content, true)
	case "clone":
		panic("not implemented!")
	default:
		log.Fatal("Command not recognised")
	}
}

func generateFile(content ComposeContent, joinAll bool) {
	commands, count := makeCommands(content, joinAll)

	commandsFile := must("open or create "+CommandsFileName+" file",
		func() (*os.File, error) {
			return os.OpenFile(CommandsFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		})

	defer commandsFile.Close()

	_, _ = commandsFile.Write([]byte(commands))

	slog.Info("Finished generating redistribution commands", slog.Int("count", count))
}

func makeCommands(content ComposeContent, joinAll bool) (string, int) {
	var (
		ReplacementMap = map[string]string{
			"nginx:latest":                 "nginx:latest",
			"debian:stable-slim":           "debian:stable-slim",
			"postgresql:14.5.0":            "bitnami/postgresql:14.5.0",
			"caddy:latest":                 "caddy:latest",
			"minio:2020.10.9-debian-10-r6": "bitnami/minio:2020.10.9-debian-10-r6",
			"minio:2023.2.10-debian-11-r1": "bitnami/minio:2023.2.10-debian-11-r1",
		}

		count = 0

		commands string
	)

	for _, value := range content.Services {
		// get the image name and version
		imageParts := strings.Split(value.Image, "/")
		imageName := imageParts[len(imageParts)-1]

		// Some images are missing from the AWS gallery, we need to manually replace them with the docker hub alternative
		remoteImageName := value.Image
		if alt, ok := ReplacementMap[imageName]; ok {
			remoteImageName = alt
		}

		// if the image name is something like bitnami/postgresql, keep that instead of just postgresql
		if len(imageParts) > 2 {
			lastIdx := len(imageParts) - 1
			imageName = imageParts[lastIdx-1] + "/" + imageParts[lastIdx]
		}

		localImageName := fmt.Sprintf("%s/%s", RegistryAddress, imageName)

		cmd := fmt.Sprintf(`docker pull %s && \
docker tag %s %s && \
docker push %s`, remoteImageName, remoteImageName, localImageName, localImageName)

		commands += cmd
		if joinAll {
			commands += " && \\"
		} else {
			commands += "\n\n"
		}

		count += 1
	}

	return strings.Trim(commands, "&& \\"), count
}

func must[T any](text string, fn func() (T, error)) T {
	value, err := fn()
	if err != nil {
		log.Fatalf("Failed to %s: %s", text, err.Error())
	}

	return value
}
