package agent

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedirect(t *testing.T) {
	raw := "https://www.baidu.com/s?ie=utf-8&f=8&rsv_bp=0&rsv_idx=1&tn=baidu&wd=123123&rsv_pq=a39292a4000155f6&rsv_t=dd12SyZUqarPhbQ2yknWg9sWwIlZqtwd9ynJjQm2J8853zxqgcDuImm3icU&rqlang=cn&rsv_enter=1&rsv_sug3=6&rsv_sug1=5&rsv_sug7=100&rsv_sug2=0&inputT=1212&rsv_sug4=1212"
	u, err := url.Parse(raw)
	assert.Nil(t, err)
	newURL := "http://www.163.com/test" + u.RequestURI()
	n, err := url.Parse(newURL)
	assert.Nil(t, err)
	assert.Equal(t, n.RequestURI(), "/test"+u.RequestURI())
}
