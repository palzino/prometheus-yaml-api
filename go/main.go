package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/alecthomas/jsonschema"
	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

// Define structs to match the Prometheus YAML structure
type PrometheusConfig struct {
	Global        GlobalConfig   `yaml:"global"`
	ScrapeConfigs []ScrapeConfig `yaml:"scrape_configs"`
}

type GlobalConfig struct {
	ScrapeInterval string            `yaml:"scrape_interval"`
	ExternalLabels map[string]string `yaml:"external_labels"`
}

type ScrapeConfig struct {
	JobName        string           `yaml:"job_name" json:"jobName"`
	Scheme         string           `yaml:"scheme,omitempty" json:"scheme,omitempty"`
	MetricsPath    string           `yaml:"metrics_path,omitempty" json:"metricsPath,omitempty"`
	ScrapeInterval string           `yaml:"scrape_interval,omitempty" json:"scrapeInterval,omitempty"`
	ScrapeTimeout  string           `yaml:"scrape_timeout,omitempty" json:"scrapeTimeout,omitempty"`
	StaticConfigs  []StaticConfig   `yaml:"static_configs,omitempty" json:"staticConfigs,omitempty"`
	BasicAuth      *BasicAuthConfig `yaml:"basic_auth,omitempty" json:"basicAuth,omitempty"`
}

type StaticConfig struct {
	Targets []string `yaml:"targets" json:"targets"`
}

type BasicAuthConfig struct {
	Username string `yaml:"username,omitempty" json:"username,omitempty"`
	Password string `yaml:"password,omitempty" json:"password,omitempty"`
}

// Global variable to store the current configuration
var currentConfig PrometheusConfig

const configFile = "prometheus.yml"

func readConfigFromFile() {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	if err := yaml.Unmarshal(data, &currentConfig); err != nil {
		log.Fatalf("Error unmarshalling config file: %v", err)
	}
}

func writeConfigToFile(config PrometheusConfig) {
	data, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatalf("Error marshalling config: %v", err)
	}

	if err := ioutil.WriteFile(configFile, data, 0644); err != nil {
		log.Fatalf("Error writing config file: %v", err)
	}
}
func restartProm() (string, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		for _, name := range container.Names {
			if name == "/prometheus" {
				fmt.Println(name)
				fmt.Println(container.ID)
				timeout := 10
				options := containertypes.StopOptions{Timeout: &timeout}
				err := cli.ContainerRestart(ctx, container.ID, options)
				if err != nil {
					fmt.Println("Error restarting container:", err)
					return "", errors.New("could not restart container")
				} else {
					fmt.Println("Container restarted successfully")
					return "", nil
				}
			}
		}
	}
	return "", err
}
func restartprom(c *gin.Context) {
	message, err := restartProm()
	if err != nil {
		log.Fatal(message, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Container restarted successfully"})
}

func getPromConf(c *gin.Context) {
	c.JSON(http.StatusOK, currentConfig)
}

func newTarget(c *gin.Context) {
	var newtarget ScrapeConfig
	if err := c.ShouldBindJSON(&newtarget); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentConfig.ScrapeConfigs = append(currentConfig.ScrapeConfigs, newtarget)
	writeConfigToFile(currentConfig)
	message, err := restartProm()
	if err != nil {
		log.Fatal(message, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"messgae": "New Target Added"})
	}

}
func deleteTarget(c *gin.Context) {
	var target ScrapeConfig

	if err := c.ShouldBindJSON(&target); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	indexToDelete := -1
	for i, config := range currentConfig.ScrapeConfigs {
		if config.JobName == target.JobName {
			indexToDelete = i
			currentConfig.ScrapeConfigs = append(currentConfig.ScrapeConfigs[:indexToDelete], currentConfig.ScrapeConfigs[indexToDelete+1:]...)
			writeConfigToFile(currentConfig)
			c.JSON(http.StatusAccepted, gin.H{"message": "scrape config deleted"})

		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "not valid target"})
		}
	}

}
func getSchema(c *gin.Context) {
	schema := jsonschema.Reflect(&PrometheusConfig{})

	// Print out the JSON schema
	schemaJSON, _ := json.MarshalIndent(schema, "", "  ")
	stringSchema := string(schemaJSON)
	c.JSON(http.StatusOK, stringSchema)
}
func main() {
	// Read configuration from file on startup
	readConfigFromFile()

	router := gin.Default()
	router.Use(cors.Default())

	// Alternatively, for more customized settings:
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                             // Allows all origins
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"}, // You can adjust methods
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},      // You can adjust headers
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.LoadHTMLGlob("src/html/*")
	router.Static("/static", "./src/static")

	// Setup routes
	router.GET("/", func(c *gin.Context) {
		c.HTML(
			// Set the HTTP status to 200 (OK)
			http.StatusOK,
			// Use the index.html template
			"index.html",
			// Pass the data that the page uses (in this case, 'title')
			gin.H{
				"title": "Home Page",
			},
		)
	})
	router.GET("/promconf", getPromConf)
	router.GET("/schema", getSchema)
	router.POST("/newtarget", newTarget)
	router.POST("/deletetarget", deleteTarget)
	router.GET("/promrestart", restartprom)
	router.Run(":7042")
}
