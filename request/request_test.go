package request

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithTimeout(t *testing.T) {
	asserts := assert.New(t)
	options := newDefaultOption()
	WithTimeout(time.Duration(5) * time.Second).apply(options)
	asserts.Equal(time.Duration(5)*time.Second, options.timeout)
}

func TestWithHeader(t *testing.T) {
	asserts := assert.New(t)
	options := newDefaultOption()
	WithHeader(map[string][]string{"Origin": {"123"}}).apply(options)
	asserts.Equal(http.Header{"Origin": []string{"123"}}, options.header)
}

func TestWithContentLength(t *testing.T) {
	asserts := assert.New(t)
	options := newDefaultOption()
	WithContentLength(10).apply(options)
	asserts.EqualValues(10, options.contentLength)
}

func TestWithContext(t *testing.T) {
	asserts := assert.New(t)
	options := newDefaultOption()
	WithContext(context.Background()).apply(options)
	asserts.NotNil(options.ctx)
}

func TestHTTPClient_Request(t *testing.T) {
	asserts := assert.New(t)
	client := NewClient()
	// 正常
	{
		resp := client.Request(
			"POST",
			"/test",
			strings.NewReader(""),
			WithContentLength(0),
			WithEndpoint("https://www.baidu.com"),
			WithTimeout(time.Duration(1)*time.Microsecond),
			WithoutHeader([]string{"origin", "origin"}),
		)
		asserts.Error(resp.Err)
		asserts.Nil(resp.Response)
	}
	// 正常 带有ctx
	{
		resp := client.Request(
			"GET",
			"https://www.baidu.com",
			strings.NewReader(""),
			WithTimeout(time.Duration(10)*time.Second),
			WithContext(context.Background()),
			WithoutHeader([]string{"s s", "s s"}),
		)
		asserts.NoError(resp.Err)
		asserts.NotNil(resp.Response)
	}

}

func TestResponse_GetResponse(t *testing.T) {
	asserts := assert.New(t)

	// 直接返回错误
	{
		resp := Response{
			Err: errors.New("error"),
		}
		content, err := resp.GetResponse()
		asserts.Empty(content)
		asserts.Error(err)
	}

	// 正常
	{
		resp := Response{
			Response: &http.Response{Body: ioutil.NopCloser(strings.NewReader("123"))},
		}
		content, err := resp.GetResponse()
		asserts.Equal("123", content)
		asserts.NoError(err)
	}
}

func TestResponse_CheckHTTPResponse(t *testing.T) {
	asserts := assert.New(t)

	// 直接返回错误
	{
		resp := Response{
			Err: errors.New("error"),
		}
		res := resp.CheckHTTPResponse(200)
		asserts.Error(res.Err)
	}

	// 404 错误
	{
		resp := Response{
			Response: &http.Response{StatusCode: 404},
		}
		res := resp.CheckHTTPResponse(200)
		asserts.Error(res.Err)
	}

	// 通过
	{
		resp := Response{
			Response: &http.Response{StatusCode: 200},
		}
		res := resp.CheckHTTPResponse(200)
		asserts.NoError(res.Err)
	}
}

func TestResponse_GetRSCloser(t *testing.T) {
	asserts := assert.New(t)

	// 直接返回错误
	{
		resp := Response{
			Err: errors.New("error"),
		}
		res, err := resp.GetRSCloser()
		asserts.Error(err)
		asserts.Nil(res)
	}

	// 正常
	{
		resp := Response{
			Response: &http.Response{ContentLength: 3, Body: ioutil.NopCloser(strings.NewReader("123"))},
		}
		res, err := resp.GetRSCloser()
		asserts.NoError(err)
		content, err := ioutil.ReadAll(res)
		asserts.NoError(err)
		asserts.Equal("123", string(content))
		offset, err := res.Seek(0, 0)
		asserts.NoError(err)
		asserts.Equal(int64(0), offset)
		offset, err = res.Seek(0, 2)
		asserts.NoError(err)
		asserts.Equal(int64(3), offset)
		_, err = res.Seek(1, 2)
		asserts.Error(err)
		asserts.NoError(res.Close())
	}

}

func TestNopRSCloser_SetFirstFakeChunk(t *testing.T) {
	asserts := assert.New(t)
	rsc := NopRSCloser{
		status: &rscStatus{},
	}
	rsc.SetFirstFakeChunk()
	asserts.True(rsc.status.IgnoreFirst)

	rsc.SetContentLength(20)
	asserts.EqualValues(20, rsc.status.Size)
}

func TestHTTPClient_TPSLimit(t *testing.T) {
	a := assert.New(t)
	client := NewClient()

	finished := make(chan struct{})
	go func() {
		client.Request(
			"POST",
			"/test",
			strings.NewReader(""),
			WithTPSLimit("TestHTTPClient_TPSLimit", 1, 1),
		)
		close(finished)
	}()
	select {
	case <-finished:
	case <-time.After(10 * time.Second):
		a.Fail("Request should be finished instantly.")
	}

	finished = make(chan struct{})
	go func() {
		client.Request(
			"POST",
			"/test",
			strings.NewReader(""),
			WithTPSLimit("TestHTTPClient_TPSLimit", 1, 1),
		)
		close(finished)
	}()
	select {
	case <-finished:
	case <-time.After(2 * time.Second):
		a.Fail("Request should be finished in 1 second.")
	}

}
