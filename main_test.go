package main

import (

	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"github.com/stretchr/testify/assert"


)

func getPublicIPAddress() (string, error) {
	resp, err := http.Get("https://api.ipify.org?format=text") 

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	ipAddress := ""
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		ipAddress = strings.TrimSpace(string(bodyBytes))
	}
	return ipAddress, nil
}



func addIPToACRNetworkRule(resourceGroup, acrName string) error {
	ipAddress, err := getPublicIPAddress()
	if err != nil {
		return fmt.Errorf("error getting public IP address: %v", err)
	}

	cmd := exec.Command("az", "acr", "network-rule", "add", "--resource-group", resourceGroup, "--name", acrName, "--ip-address", ipAddress)
	err = cmd.Run()
	fmt.Println(err)
	if err != nil {
		return fmt.Errorf("error adding IP address to ACR network rule: %v", err)
	}

	fmt.Println("Successfully added IP address to ACR network rule.")
	return nil
}



func dockerLogin(acrURL, acrUsername, acrPassword string) error {
	cmd := exec.Command("docker", "login", acrURL, "-u", acrUsername, "-p", acrPassword)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error logging in to ACR: %v", err)
	}

	fmt.Println("Successfully logged in to ACR.")
	return nil
}


func dockerPushImage(imageNameLocal, acrURL, acrUsername, acrPassword string) (bool, error) {
	acrImageName := fmt.Sprintf("%s/%s", acrURL, imageNameLocal+".1")
	fmt.Println(acrImageName)

	// Tag the local Docker image with the ACR image name
	cmd := exec.Command("docker", "tag", imageNameLocal, acrImageName)
	err := cmd.Run()
	if err != nil {
		return false, fmt.Errorf("error tagging local Docker image with ACR image name: %v", err)
	}

	// Login to the ACR
	err = dockerLogin(acrURL, acrUsername, acrPassword)
	if err != nil {
		return false, err
	}

	// Push the tagged image to the ACR
	cmd = exec.Command("docker", "push", acrImageName)
	err = cmd.Run()
	if err != nil {
		return false, fmt.Errorf("error pushing image to ACR: %v", err)
	}

	fmt.Println("Image pushed to ACR successfully.")
	return true, nil
}



func TestTerraformOutputs(t *testing.T) {
	// Replace these variables with your own values
	resourceGroup := "<BLOG>"
	acrName := "<testacrmk2>"
	acrURL := "<testacrmk2.azurecr.io>"
	acrPassword := "< >"
	imageNameLocal := "hello-world:latest"
	expectedPushedResponse := true

	// Add IP to ACR network rule
	err := addIPToACRNetworkRule(resourceGroup, acrName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Push the image to ACR
	actualPushedResponse, err := dockerPushImage(imageNameLocal, acrURL, acrName, acrPassword)
	fmt.Println(actualPushedResponse)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	t.Run(fmt.Sprintf("Checking the images is pushed or not to the Private ACR: %s", acrName), func(t *testing.T) {

        assert.Equal(t, expectedPushedResponse, actualPushedResponse, "publicNetworkAccess mismatch")

    })

}
