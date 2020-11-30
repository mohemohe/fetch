package fetch

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type (
	Option struct {
		Method  string
		Headers map[string]string
		Body    io.Reader
	}
	Result struct {
		res        *http.Response
		Status     int
		StatusText string
		OK         bool
	}
)

func Fetch(url string, option ...Option) (*Result, error) {
	if len(option) != 1 {
		option = make([]Option, 1)
		option[0] = Option{
			Method:  http.MethodGet,
			Headers: make(map[string]string),
		}
	}
	opt := option[0]
	if opt.Method == "" {
		opt.Method = http.MethodGet
	}

	req, err := http.NewRequest(opt.Method, url, opt.Body)
	if err != nil {
		return nil, err
	}
	for k, v := range opt.Headers {
		req.Header.Set(k, v)
	}

	client := new(http.Client)
	res, err := client.Do(req)
	if res == nil {
		return nil, errors.New("http client initialize error")
	}
	statusText := strconv.Itoa(res.StatusCode)
	f := &Result{
		res:        res,
		Status:     res.StatusCode,
		StatusText: statusText,
		OK:         strings.HasPrefix(statusText, "2"),
	}
	return f, err
}

func (result *Result) Response() *http.Response {
	return result.res
}

func (result *Result) Bytes() ([]byte, error) {
	defer result.res.Body.Close()
	return ioutil.ReadAll(result.res.Body)
}

func (result *Result) Text() (string, error) {
	bytes, err := result.Bytes()
	return string(bytes), err
}

func (result *Result) JSON(out interface{}) error {
	bytes, err := result.Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, out)
}

func (result *Result) Map() (map[string]interface{}, error) {
	m := map[string]interface{}{}
	err := result.JSON(&m)
	return m, err
}
