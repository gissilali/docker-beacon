package main

import (
	"docker-watch/pkg/discord"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type DockerContainer struct {
	Names  string `json:"Names"`
	Status string `json:"Status"`
	ID     string `json:"ID"`
}

func main() {
	bot, err := discord.OpenSession()

	if err != nil {
		fmt.Printf("error creating docker session: %v\n", err)
		return
	}

	output, err := runDockerPs()

	if err != nil {
		fmt.Printf("error executing docker ps: %s\n", err)
	}

	stringOutput := string(output)

	dockerContainers := getDockerContainersFromOutput(strings.Split(stringOutput, "\n"))

	if len(dockerContainers) == 0 {
		fmt.Println("All good now!")
	} else {
		sendNotification(bot, dockerContainers)
	}

	_, err = writeToFile("docker.json", "["+stringOutput+"]")

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

func writeToFile(filename string, data string) (any, error) {
	file, err := os.Create(filename)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	_, err = file.WriteString(data)

	if err != nil {
		return nil, err
	}

	return nil, nil
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
