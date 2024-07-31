package main

import (
	"docker-watch/pkg/discord"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type DockerContainer struct {
	Names  string `json:"Names"`
	Status string `json:"Status"`
	ID     string `json:"ID"`
}

func main() {
	bot, err := openDiscordSession()
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	output, err := runDockerPs()

	if err != nil {
		fmt.Printf("error executing docker ps: %s\n", err)
	}

	dockerContainers := parseDockerContainers(string(output))

	if len(dockerContainers) == 0 {
		fmt.Println("All good now!")
	} else {
		sendNotification(bot, dockerContainers)
	}
}

func parseDockerContainers(outputString string) []DockerContainer {
	return getDockerContainersFromOutput(strings.Split(outputString, "\n"))
}

func openDiscordSession() (*discord.DockerBot, error) {
	bot, err := discord.OpenSession()

	if err != nil {
		return nil, fmt.Errorf("error opening Discord session: %v", err)
	}

	return bot, nil
}

func sendNotification(bot *discord.DockerBot, containers []DockerContainer) {
	// Map the slice of structs to a slice of strings
	names := make([]string, len(containers))
	for i, container := range containers {
		names[i] = container.Names
	}

	var commaSeparatedNames string
	if len(names) == 1 {
		commaSeparatedNames = names[0]
	} else if len(names) > 1 {
		commaSeparatedNames = strings.Join(names[:len(names)-1], "`, `") + "` & `" + names[len(names)-1]
	}

	fmt.Println(commaSeparatedNames)
	message, err := bot.SendMessage(fmt.Sprintf("%d dockers dead %s", len(containers), "`"+commaSeparatedNames) + "`")
	if err != nil {
		return
	}

	fmt.Printf("Sent a message:%v\n", message.Content)
}

func runDockerPs() ([]byte, error) {
	outputTemplate := `{"ID": "{{ .ID }}","Image": "{{ .Image }}","Names": "{{ .Names }}", "Ports": "{{ .Ports }}", "Status": "{{ .Status }}"}`
	cmd := exec.Command("docker", "ps", "-a", "--format", outputTemplate, "--filter", "status=exited")
	output, err := cmd.CombinedOutput()
	return output, err
}

func getDockerContainersFromOutput(input []string) []DockerContainer {
	var result []DockerContainer
	for _, str := range input {
		if str != "" {
			var container DockerContainer
			err := json.Unmarshal([]byte(str), &container)
			if err != nil {
				fmt.Println("Error:-->", err)
			}
			result = append(result, container)
		}
	}
	return result
}
