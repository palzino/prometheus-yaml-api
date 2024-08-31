package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alecthomas/jsonschema"
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
type IPAdress struct {
	IP string `json:"ip"`
}

// Global variable to store the current configuration
var currentConfig PrometheusConfig

const configFile = "./config/prometheus.yml"

func readConfigFromFile() {
	data, err := os.ReadFile(configFile)
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

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		log.Fatalf("Error writing config file: %v", err)
	}
}

func restartNormalPromHandler(c *gin.Context) {
	var ip IPAdress
	c.BindJSON(&ip)
	ipaddr := "http://" + ip.IP + "/-/reload"
	fmt.Println(ipaddr)
	// Create a new HTTP request
	req, err := http.NewRequest("POST", ipaddr, bytes.NewBuffer([]byte("")))
	if err != nil {
		fmt.Errorf("error creating request: %w", err)
		return
	}

	// Set the request headers if needed
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {

		fmt.Errorf("error sending request: %w", err)
		return

	}

	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Container restarted successfully"})
}

func getPromConf(c *gin.Context) {
	c.JSON(http.StatusOK, currentConfig)
}

func newTarget(c *gin.Context) {
	var newTarget ScrapeConfig
	if err := c.ShouldBindJSON(&newTarget); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentConfig.ScrapeConfigs = append(currentConfig.ScrapeConfigs, newTarget)
	writeConfigToFile(currentConfig)

	c.JSON(http.StatusOK, gin.H{"message": "New Target Added"})
}

func deleteTarget(c *gin.Context) {
	var target ScrapeConfig

	if err := c.ShouldBindJSON(&target); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, config := range currentConfig.ScrapeConfigs {
		if config.JobName == target.JobName {
			currentConfig.ScrapeConfigs = append(currentConfig.ScrapeConfigs[:i], currentConfig.ScrapeConfigs[i+1:]...)
			writeConfigToFile(currentConfig)
			c.JSON(http.StatusOK, gin.H{"message": "Scrape config deleted"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Target not found"})
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
	router.POST("/reload", restartNormalPromHandler)
	router.Run(":7042")
}
