package wirescli

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	Rest *resty.Client
}

// New returns a new instance of the nube common apis
func New(url string, port int) *Client {
	rest := &Client{
		Rest: resty.New(),
	}
	rest.Rest.SetBaseURL(fmt.Sprintf("http://%s:%d", url, port))
	return rest
}

type Path struct {
	Path string
}

var Paths = struct {
	Auth   Path
	Upload Path
	Export Path
}{
	Auth:   Path{Path: "/api/auth/login"},
	Upload: Path{Path: "/api/editor/c/0/import"},
	Export: Path{Path: "/api/editor/export/all"},
}

type Response struct {
	StatusCode int         `json:"code"`
	Message    interface{} `json:"message"`
	resty      *resty.Response
}

func (response Response) buildResponse(resp *resty.Response, err error) *Response {
	response.StatusCode = resp.StatusCode()
	response.resty = resp
	if resp.IsError() {
		response.Message = resp.Error()
	}
	if resp.StatusCode() == 0 {
		response.Message = "server is unreachable"
		response.StatusCode = 503
	}
	return &response
}

type WiresTokenBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}

func (inst *Client) GetToken(body *WiresTokenBody) (token *Token, response *Response) {
	path := fmt.Sprintf(Paths.Auth.Path)
	response = &Response{}
	resp, err := inst.Rest.R().
		SetBody(body).
		SetResult(&Token{}).
		Post(path)
	return resp.Result().(*Token), response.buildResponse(resp, err)
}

type NodesBody struct {
	Nodes interface{} `json:"nodes"`
	Pos   []float64   `json:"pos"`
}

func (inst *Client) Upload(body *NodesBody, token string) (ok bool, response *Response) {
	path := fmt.Sprintf(Paths.Upload.Path)
	response = &Response{}
	resp, err := inst.Rest.R().
		SetBody(body).
		SetHeaders(map[string]string{
			"token": token,
		}).
		SetAuthToken(token).
		Post(path)
	if resp.IsSuccess() {
		return true, response.buildResponse(resp, err)
	}
	return false, response.buildResponse(resp, err)
}

type WiresExport struct {
	Objects     interface{}
	Errors      []interface{} `json:"errors"`
	ContainerId string        `json:"containerId"`
	Total       int           `json:"total"`
	Message     string        `json:"message"`
}

func (inst *Client) Backup(token string) (data *WiresExport, err error) {
	path := fmt.Sprintf(Paths.Export.Path)
	resp, err := inst.Rest.R().
		SetHeaders(map[string]string{
			"token": token,
		}).
		SetResult(&WiresExport{}).
		SetAuthToken(token).
		Get(path)
	if resp.IsSuccess() {
		data = resp.Result().(*WiresExport)
		return data, err
	} else {
		data = resp.Error().(*WiresExport)
		return data, err
	}
}
