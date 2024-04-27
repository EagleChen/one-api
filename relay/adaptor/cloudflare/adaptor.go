package cloudflare

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/songquanpeng/one-api/relay/adaptor/openai"
	"github.com/songquanpeng/one-api/relay/meta"
	"github.com/songquanpeng/one-api/relay/relaymode"

	"github.com/gin-gonic/gin"
	"github.com/songquanpeng/one-api/relay/adaptor"
	"github.com/songquanpeng/one-api/relay/model"
)

type Adaptor struct {
}

func (a *Adaptor) Init(meta *meta.Meta) {

}

func (a *Adaptor) GetRequestURL(meta *meta.Meta) (string, error) {
	// cloudflare worker url需要使用 accountID, 用户必须填入包含 accountID 的代理地址
	// 形如：https://api.cloudflare.com/client/v4/accounts/%s/ai/run
	fullRequestURL := fmt.Sprintf("%s/%s", meta.BaseURL, meta.ActualModelName)
	return fullRequestURL, nil
}

func (a *Adaptor) SetupRequestHeader(c *gin.Context, req *http.Request, meta *meta.Meta) error {
	adaptor.SetupCommonRequestHeader(c, req, meta)
	req.Header.Set("Authorization", "Bearer "+meta.APIKey)
	return nil
}

func (a *Adaptor) ConvertRequest(c *gin.Context, relayMode int, request *model.GeneralOpenAIRequest) (any, error) {
	if request == nil {
		return nil, errors.New("request is nil")
	}
	// 目前仅支持 text generation， 后续可以支持 asr, text to image 等
	switch relayMode {
	case relaymode.Embeddings:
		return nil, errors.New("not supported yet")
	default:
		return request, nil
	}
}

func (a *Adaptor) ConvertImageRequest(request *model.ImageRequest) (any, error) {
	if request == nil {
		return nil, errors.New("request is nil")
	}
	return request, nil
}

func (a *Adaptor) DoRequest(c *gin.Context, meta *meta.Meta, requestBody io.Reader) (*http.Response, error) {
	return adaptor.DoRequestHelper(a, c, meta, requestBody)
}

func (a *Adaptor) DoResponse(c *gin.Context, resp *http.Response, meta *meta.Meta) (usage *model.Usage, err *model.ErrorWithStatusCode) {
	if meta.IsStream {
		streamHandler(c, resp, meta)
	} else {
		errMsg := "not implemented"
		switch meta.Mode {
		case relaymode.Embeddings:
			err = openai.ErrorWrapper(errors.New(errMsg), errMsg, http.StatusInternalServerError)
		default:
			err, usage = textGenerationHandler(c, resp)
		}
	}
	return
}

func (a *Adaptor) GetModelList() []string {
	return ModelList
}

func (a *Adaptor) GetChannelName() string {
	return "cloudflare-workers-ai"
}
