// Copyright © 2016 River Yang <comicme_yanghe@icloud.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"os"
	"encoding/json"
	"github.com/docker/docker/daemon/network"
	"github.com/spf13/cobra"
	"strconv"
	"github.com/gobuild/log"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ruleng",
	RunE: func(cmd *cobra.Command, args []string) error {
		hostip, err := cmd.Flags().GetString("hostip")
		if err != nil {
			return err
		}

		hostIP := os.Getenv("HOSTIP")
		if hostIP == "" {
			hostIP = hostip
		}

		log.Printf("Host IP: %v", hostIP)

		port, err := cmd.Flags().GetUint("docker-port")
		if err != nil {
			return err
		}

		dockerPort := os.Getenv("DOCKER_PORT")
		if dockerPort == "" {
			dockerPort = strconv.Itoa(int(port))
		}

		log.Printf("Docker port: %v", dockerPort)

		name, err := cmd.Flags().GetString("hostname")
		if err != nil {
			return err
		}

		hostname := os.Getenv("HOSTNAME")
		if name != "" {
			hostname = name
		}

		log.Printf("Hostname: %v", hostname)

		profile, err := cmd.Flags().GetString("profile")
		if err != nil {
			return err
		}

		log.Printf("Profile path: %v", profile)

		resp, err := http.Get("http://" + hostIP + ":" + dockerPort + "/containers/" + hostname + "/json")
		if err != nil {
			return err
		}

		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		maps := make(map[string]interface{})
		err = json.Unmarshal(bytes, &maps)
		if err != nil {
			return err
		}

		networkSettings, err := json.Marshal(maps["NetworkSettings"])
		if err != nil {
			return err
		}

		settings := network.Settings{}
		err = json.Unmarshal(networkSettings, &settings)

		var profiles string
		for port, bindings := range settings.Ports {
			for _, binding := range bindings {
				log.Println(binding.HostIP + ":" + binding.HostPort + " -> " + port.Port())
				if binding.HostPort != "" && port.Port() != "" {
					p := "export PORT_" + port.Port() + "=" + binding.HostPort + "\n"
					log.Println(p)
					profiles += p
				}
			}
		}

		if profiles != "" {
			os.Remove(profile)
			ioutil.WriteFile(profile, []byte(profiles), os.ModePerm)
		}

		return nil
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.Flags().String("hostip", "127.0.0.1", "host ip")
	RootCmd.Flags().Uint("docker-port", 2375, "Docker daemon port")
	RootCmd.Flags().String("hostname", "", "Docker container hostname")
	RootCmd.Flags().String("profile", "/tmp/profile", "Profile template path")
}

