/*
Copyright © 2020 BizFly Cloud

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bizflycloud/gobizfly"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"github.com/spf13/cobra"
)

var (
	cfgFile  string
	email    string
	password string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bizfly",
	Short: "BizFly Cloud Command Line",
	Long:  `BizFly Cloud Command Line`,
	PreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("Pre run")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bizfly.yaml)")

	rootCmd.PersistentFlags().StringVar(&email, "email", os.Getenv("BIZFLY_CLOUD_EMAIL"), "Your BizFly Cloud Email")
	rootCmd.MarkFlagRequired("email")

	rootCmd.PersistentFlags().StringVar(&password, "password", os.Getenv("BIZFLY_CLOUD_PASSWORD"), "Your BizFly CLoud Password")
	rootCmd.MarkFlagRequired("password")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".bizfly" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".bizfly")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func getApiClient(cmd *cobra.Command) (*gobizfly.Client, context.Context) {
	email, err := cmd.Flags().GetString("email")
	if email == "" {
		log.Fatal("Email is required")
	}
	password,  err := cmd.Flags().GetString("password")
	if password == "" {
		log.Fatal("Password is required")
	}

	client, err := gobizfly.NewClient(gobizfly.WithTenantName(email))

	if err != nil {
		log.Fatal(err)
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()
	//TODO Get token from cache
	tok, err := client.Token.Create(ctx, &gobizfly.TokenCreateRequest{AuthMethod: "password", Username: email, Password: password})
	if err != nil {
		log.Fatal(err)
	}

	client.SetKeystoneToken(tok.KeystoneToken)

	return client, ctx
}
