package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

var debugMode bool

func main() {
	configFilenamePtr := flag.String("config-filename", "config-example.json", "the configuration filename")
	listenAddressPtr := flag.String("listen-address", "", "the address on which to listen")
	listenPortPtr := flag.Int("listen-port", 0, "the port on which to listen")
	debugModePtr := flag.Bool("debug", false, "enable debug mode")
	flag.Parse()

	configFilename := *configFilenamePtr
	listenAddress := *listenAddressPtr
	listenPort := *listenPortPtr
	debugMode := *debugModePtr

	if !debugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	if listenPort == 0 {
		// Grab the listen port from the environment if it exists.
		listenPortStr, gotPort := os.LookupEnv("NOMAD_PORT_http")
		if gotPort {
			var err error
			listenPort, err = strconv.Atoi(listenPortStr)
			if err != nil {
				log.Fatal("Failed to convert port from environment to integer!")
			}
		}
	}

	config, err := LoadConfig(configFilename)
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.Use(gin.Recovery())

	r.POST("/:endpoint/*key", func(c *gin.Context) {
		endpoint := c.Param("endpoint")
		key := c.Param("key")
		reader := bufio.NewReader(c.Request.Body)

		err := config.processRequest(endpoint, key, reader)
		if err == nil {
			c.Status(http.StatusNoContent)
			return
		}

		errStr := fmt.Sprintf("%s", err)
		if len(errStr) == 3 {
			// Assume it's a status code.
			statusCode, err := strconv.Atoi(errStr)
			if err == nil {
				c.Status(statusCode)
				return
			}
			errStr = fmt.Sprintf("%s", err)
		}

		c.Data(500, gin.MIMEPlain, []byte(errStr))
	})

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", listenAddress, listenPort))
	if err != nil {
		panic(err)
	}

	log.Printf("API server listening on %s:%d\n", listenAddress, listener.Addr().(*net.TCPAddr).Port)
	panic(http.Serve(listener, r))
}
