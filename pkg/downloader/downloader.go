package downloader

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// PaperMCAPI represents the PaperMC API endpoints
const (
	PaperMCBaseURL = "https://api.papermc.io/v2/projects"
	PaperMCProject = "paper"
)

// PaperMCVersionsResponse represents the response from the PaperMC API for versions
type PaperMCVersionsResponse struct {
	ProjectID   string   `json:"project_id"`
	ProjectName string   `json:"project_name"`
	Versions    []string `json:"versions"`
}

// PaperMCBuildsResponse represents the response from the PaperMC API for builds
type PaperMCBuildsResponse struct {
	ProjectID   string `json:"project_id"`
	ProjectName string `json:"project_name"`
	Version     string `json:"version"`
	Builds      []int  `json:"builds"`
}

// PaperMCBuildResponse represents the response from the PaperMC API for a specific build
type PaperMCBuildResponse struct {
	ProjectID   string `json:"project_id"`
	ProjectName string `json:"project_name"`
	Version     string `json:"version"`
	Build       int    `json:"build"`
	Downloads   struct {
		Application struct {
			Name   string `json:"name"`
			Sha256 string `json:"sha256"`
		} `json:"application"`
	} `json:"downloads"`
}

// DownloadPaperMC downloads the PaperMC server jar
func DownloadPaperMC(serverPath, version string) (string, error) {
	// If version is "latest", get the latest version
	if version == "latest" {
		var err error
		version, err = getLatestPaperMCVersion()
		if err != nil {
			return "", fmt.Errorf("failed to get latest PaperMC version: %w", err)
		}
		fmt.Printf("Using latest PaperMC version: %s\n", version)
	}

	// Get the latest build for the version
	latestBuild, err := getLatestPaperMCBuild(version)
	if err != nil {
		return "", fmt.Errorf("failed to get latest PaperMC build: %w", err)
	}

	// Get the build details
	buildDetails, err := getPaperMCBuildDetails(version, latestBuild)
	if err != nil {
		return "", fmt.Errorf("failed to get PaperMC build details: %w", err)
	}

	// Download the server jar
	jarName := buildDetails.Downloads.Application.Name
	downloadURL := fmt.Sprintf("%s/%s/versions/%s/builds/%d/downloads/%s",
		PaperMCBaseURL, PaperMCProject, version, latestBuild, jarName)

	jarPath := filepath.Join(serverPath, jarName)
	if err := downloadFile(downloadURL, jarPath); err != nil {
		return "", fmt.Errorf("failed to download PaperMC server jar: %w", err)
	}

	return jarPath, nil
}

// DownloadVanilla downloads the vanilla Minecraft server jar
func DownloadVanilla(serverPath, version string) (string, error) {
	// For vanilla, we need to handle the version differently
	// For now, we'll use a hardcoded URL for the latest version
	// In a real implementation, we would need to fetch the version manifest from Mojang
	
	// This is a placeholder URL for the latest vanilla server
	vanillaURL := "https://piston-data.mojang.com/v1/objects/8dd1a28015f51b1803213892b50b7b4fc76e594d/server.jar"
	jarName := "minecraft_server." + version + ".jar"
	jarPath := filepath.Join(serverPath, jarName)
	
	if err := downloadFile(vanillaURL, jarPath); err != nil {
		return "", fmt.Errorf("failed to download vanilla server jar: %w", err)
	}
	
	return jarPath, nil
}

// DownloadFabric downloads the Fabric server jar
func DownloadFabric(serverPath, mcVersion, loaderVersion string) (string, error) {
	// Default installer version
	installerVersion := "1.0.1"
	
	// If no loader version is provided, use the default
	if loaderVersion == "" {
		loaderVersion = "0.16.10"
	}
	
	// Construct the URL for the Fabric server jar
	fabricServerURL := fmt.Sprintf("https://meta.fabricmc.net/v2/versions/loader/%s/%s/%s/server/jar", 
		mcVersion, loaderVersion, installerVersion)
	
	// Expected jar filename
	jarName := fmt.Sprintf("fabric-server-mc.%s-loader.%s-launcher.%s.jar", 
		mcVersion, loaderVersion, installerVersion)
	jarPath := filepath.Join(serverPath, jarName)
	
	fmt.Printf("Downloading Fabric server jar for Minecraft %s with loader %s...\n", 
		mcVersion, loaderVersion)
	
	// Download the Fabric server jar
	if err := downloadFile(fabricServerURL, jarPath); err != nil {
		return "", fmt.Errorf("failed to download Fabric server jar: %w", err)
	}
	
	fmt.Println("Fabric server jar downloaded successfully")
	fmt.Println("Note: Most mods will also require you to install Fabric API into the mods folder")
	
	// Create mods directory
	modsDir := filepath.Join(serverPath, "mods")
	if err := os.MkdirAll(modsDir, 0755); err != nil {
		fmt.Printf("Warning: Failed to create mods directory: %v\n", err)
	}
	
	return jarPath, nil
}

// getLatestPaperMCVersion gets the latest PaperMC version
func getLatestPaperMCVersion() (string, error) {
	url := fmt.Sprintf("%s/%s", PaperMCBaseURL, PaperMCProject)
	
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get PaperMC versions: %s", resp.Status)
	}
	
	var versionsResp PaperMCVersionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&versionsResp); err != nil {
		return "", err
	}
	
	if len(versionsResp.Versions) == 0 {
		return "", fmt.Errorf("no PaperMC versions found")
	}
	
	// Return the latest version (last in the list)
	return versionsResp.Versions[len(versionsResp.Versions)-1], nil
}

// getLatestPaperMCBuild gets the latest build for a PaperMC version
func getLatestPaperMCBuild(version string) (int, error) {
	url := fmt.Sprintf("%s/%s/versions/%s", PaperMCBaseURL, PaperMCProject, version)
	
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to get PaperMC builds: %s", resp.Status)
	}
	
	var buildsResp PaperMCBuildsResponse
	if err := json.NewDecoder(resp.Body).Decode(&buildsResp); err != nil {
		return 0, err
	}
	
	if len(buildsResp.Builds) == 0 {
		return 0, fmt.Errorf("no PaperMC builds found for version %s", version)
	}
	
	// Return the latest build (last in the list)
	return buildsResp.Builds[len(buildsResp.Builds)-1], nil
}

// getPaperMCBuildDetails gets the details for a specific PaperMC build
func getPaperMCBuildDetails(version string, build int) (*PaperMCBuildResponse, error) {
	url := fmt.Sprintf("%s/%s/versions/%s/builds/%d", PaperMCBaseURL, PaperMCProject, version, build)
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get PaperMC build details: %s", resp.Status)
	}
	
	var buildResp PaperMCBuildResponse
	if err := json.NewDecoder(resp.Body).Decode(&buildResp); err != nil {
		return nil, err
	}
	
	return &buildResp, nil
}

// downloadFile downloads a file from a URL to a local path
func downloadFile(url, filePath string) error {
	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()
	
	// Create a client with a timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	// Get the data
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: %s", resp.Status)
	}
	
	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	
	return nil
}
