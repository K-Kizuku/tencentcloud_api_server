package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/K-Kizuku/tencentcloud_api_server/lib/config"
	"github.com/K-Kizuku/tencentcloud_api_server/lib/sign"
	"github.com/rs/cors"
	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/regions"
	trtc "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/trtc/v20190722"
)

func main() {
	config.LoadEnv()

	credential := common.NewCredential(
		config.SecretID,
		config.SecretKey,
	)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /userSig", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		req := UserSigReq{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		userSig, err := sign.GenUserSig(int(config.SdkAppID), config.SdkAppSecret, req.ID, 3600)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"userSig": userSig})
	})
	mux.HandleFunc("GET /mixStreaming", func(w http.ResponseWriter, r *http.Request) {

		cpf := profile.NewClientProfile()
		cpf.HttpProfile.ReqMethod = "POST"
		cpf.HttpProfile.ReqTimeout = 30
		cpf.HttpProfile.Endpoint = "trtc.tencentcloudapi.com"
		cpf.Language = "en-US"

		client, err := trtc.NewClient(credential, regions.Singapore, cpf)

		if err != nil {
			fmt.Print(err)
		}

		req := trtc.NewStartPublishCdnStreamRequest()
		req.SdkAppId = common.Uint64Ptr(config.SdkAppID)
		req.RoomId = common.StringPtr("hoge")
		req.RoomIdType = common.Uint64Ptr(1)
		req.WithTranscoding = common.Uint64Ptr(1)
		req.AudioParams = &trtc.McuAudioParams{
			AudioEncode: &trtc.AudioEncode{
				BitRate:    common.Uint64Ptr(8),
				Channel:    common.Uint64Ptr(2),
				SampleRate: common.Uint64Ptr(44100),
			},
		}
		req.VideoParams = &trtc.McuVideoParams{
			VideoEncode: &trtc.VideoEncode{
				BitRate: common.Uint64Ptr(7500),
				Fps:     common.Uint64Ptr(60),
				Width:   common.Uint64Ptr(1920),
				Height:  common.Uint64Ptr(1080),
				Gop:     common.Uint64Ptr(1),
			},
		}
		req.VideoParams.LayoutParams = &trtc.McuLayoutParams{
			MixLayoutMode:          common.Uint64Ptr(3),
			PureAudioHoldPlaceMode: common.Uint64Ptr(0),
		}

		d := &trtc.McuPublishCdnParam{
			PublishCdnUrl: common.StringPtr(config.RTMP_URL),
			IsTencentCdn:  common.Uint64Ptr(1),
		}
		var dd []*trtc.McuPublishCdnParam
		req.PublishCdnParams = append(dd, d)

		req.AgentParams = &trtc.AgentParams{
			UserId:      common.StringPtr(config.AgentName),
			UserSig:     common.StringPtr(config.AgentSign),
			MaxIdleTime: common.Uint64Ptr(30),
		}

		response, err := client.StartPublishCdnStream(req)
		if _, ok := err.(*errors.TencentCloudSDKError); ok {
			fmt.Printf("An API error has returned: %s", err)
			return
		}
		if err != nil {
			panic(err)
		}

		fmt.Printf("%s", response.ToJsonString())

	})

	c := cors.AllowAll()
	handler := c.Handler(mux)
	log.Println("listen and serve ... on port 8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}

type UserSigReq struct {
	ID string `json:"id"`
}
