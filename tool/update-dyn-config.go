package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Http Http `yaml:"http"`
}

type Http struct {
	Routers     map[string]Router     `yaml:"routers"`
	Services    map[string]Service    `yaml:"services"`
	Middlewares map[string]Middleware `yaml:"middlewares"`
}

type Router struct {
	Rule        string   `yaml:"rule"`
	Service     string   `yaml:"service"`
	EntryPoints []string `yaml:"entryPoints"`
	Middlewares []string `yaml:"middlewares"`
}

type Service struct {
	LoadBalancer LoadBalancer `yaml:"loadBalancer"`
}

type LoadBalancer struct {
	PassHostHeader bool     `yaml:"passHostHeader"`
	Servers        []Server `yaml:"servers"`
}

type Server struct {
	URL string `yaml:"url"`
}

type Middleware struct {
	Plugin struct {
		DemoWasm struct {
			Headers map[string]string `yaml:"headers"`
		} `yaml:"demowasm"`
	} `yaml:"plugin"`
}

func generateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	sb := strings.Builder{}
	for i := 0; i < n; i++ {
		sb.WriteByte(letters[rand.IntN(len(letters))])
	}
	return sb.String()
}

func main() {
	filename := flag.String("file", "dyn.yaml", "YAML file to read and write")
	loopCount := flag.Int("count", 1, "Number of times to modify the URL")
	sleepDuration := flag.Duration("sleep", 0, "Sleep duration in seconds between each modification")
	flag.Parse()

	for i := 0; i < *loopCount; i++ {
		// Read the YAML file
		data, err := os.ReadFile(*filename)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		// Unmarshal the YAML data
		var config Config
		err = yaml.Unmarshal(data, &config)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		servers := make([]Server, 3)
		for i := 0; i < len(servers); i++ {
			servers[i] = Server{URL: fmt.Sprintf("http://127.0.0.1:8000/%s", generateRandomString(32))}
		}
		// config.Http.Services["service-foo"].LoadBalancer.Servers = servers
		service := config.Http.Services["service-foo"]
		service.LoadBalancer.Servers = servers
		config.Http.Services["service-foo"] = service

		for i := 0; i < 100; i++ {
			svc := fmt.Sprintf("svc-%d", i)
			config.Http.Services[svc] = service

			config.Http.Routers[fmt.Sprintf("route-%s", svc)] = Router{
				Rule:        fmt.Sprintf("Host(`%s.example.com`)", svc),
				Service:     svc,
				EntryPoints: []string{"web"},
				Middlewares: []string{"demowasm"},
			}
		}

		// Marshal the YAML data back to bytes
		out, err := yaml.Marshal(&config)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		// Write the modified YAML back to the file
		err = os.WriteFile(*filename, out, 0644)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		fmt.Println("YAML file updated successfully.")
		time.Sleep(*sleepDuration)
	}
}
