/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	_ "github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"telegrambot/internal/helpers"
	"telegrambot/internal/mysql/config"
	"telegrambot/internal/service"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "telegrambot",
	Short: "Mama's telegram bot",
	Long:  `Telegram bot for mama's clients'`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}
var Service *service.Service

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	viper.AddConfigPath(helpers.GetConfigDir())
	viper.SetConfigName("telegramBot")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}

	err = viper.UnmarshalKey("telegram", Service.Config)
	if err != nil {
		log.Fatalln(err)
	}

	var configDB *mysql.Config
	err = viper.UnmarshalKey("mysql", &configDB)
	if err != nil {
		log.Fatalln(err)
	}

	db, err := sqlx.Open("mysql", configDB.User+":"+configDB.Password+"@/"+configDB.DB+"?parseTime=true")
	if err != nil {
		log.Fatalln(err)
	}
	Service.Db = db

	bot, err := tgbotapi.NewBotAPI(Service.Config.Token)
	bot.Debug = Service.Config.Dev
	if err != nil {
		panic(err)
	}

	Service.Bot = bot
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "telegram", "", "telegram file (default is $HOME/.telegrambot.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
