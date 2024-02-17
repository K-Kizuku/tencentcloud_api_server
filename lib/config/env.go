package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	SdkAppID     uint64
	SdkAppSecret string
	SecretID     string
	SecretKey    string
	AgentName    string
	AgentSign    string
	RTMP_URL     string
)

func LoadEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Printf("読み込み出来ませんでした: %v", err)
	}

	SdkAppID, err = strconv.ParseUint(os.Getenv("TENCENTCLOUD_APP_API_ID"), 10, 64)
	if err != nil {
		panic(err)
	}
	SecretID = os.Getenv("TENCENTCLOUD_SECRET_ID")
	SecretKey = os.Getenv("TENCENTCLOUD_SECRET_KEY")
	SdkAppSecret = os.Getenv("TENCENTCLOUD_API_SECRET_KEY")
	AgentName = os.Getenv("AGENT_NAME")
	AgentSign = os.Getenv("AGENT_SIGNATURE")
	RTMP_URL = os.Getenv("RTMP_PUSH_URL")
}
