package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	config "github.com/bramble555/blog/conf"
	"github.com/bramble555/blog/controller"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func TestRateLimitMiddleware(t *testing.T) {
	// 1. 设置全局配置模拟环境
	global.Config = &config.Config{
		RateLimit: config.RateLimit{
			Enable:      true,
			Store:       "memory",
			GlobalRPS:   1, // 每秒只允许 1 个请求 (生成速率)
			GlobalBurst: 1, // 桶的大小只有 1 (最多存 1 个令牌)
		},
	}
	global.Log = logrus.New() // 初始化日志，防止空指针panic

	// 2. 设置 Gin 引擎和路由
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.RateLimitMiddleware()) // 加载限流中间件
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// 3. 发起请求测试

	// [第一次请求]：应该成功
	// 原因：初始化时桶是满的（容量为1，有1个令牌）。请求拿走这唯一的一个令牌，成功通过。
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/ping", nil)
	r.ServeHTTP(w1, req1)

	if w1.Code != 200 {
		t.Errorf("第一次请求 http 状态码 = %d, 期望 200", w1.Code)
	}

	var resp1 map[string]interface{}
	json.Unmarshal(w1.Body.Bytes(), &resp1)
	if resp1["message"] != "pong" {
		t.Logf("第一次请求成功返回: %v", resp1)
	}

	// [第二次请求]：应该失败（被限流）
	// 原因：
	// 1. 全局配置是每秒生成 1 个令牌 (GlobalRPS = 1)。
	// 2. 第一个请求刚刚把桶里唯一的一个令牌拿走了，桶现在是空的 (0/1)。
	// 3. 第二次请求是紧接着发起的（几乎 0 毫秒间隔），新的令牌还没来得及生成（需要等1秒）。
	// 4. 所以没拿到令牌，RateLimitMiddleware 应该拦截并返回错误。
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/ping", nil)
	r.ServeHTTP(w2, req2)

	// 注意：中间件拦截后，HTTP 状态码还是 200 (这是框架的设定)，但业务 Code 会是错误码
	if w2.Code != 200 {
		t.Errorf("第二次请求 http 状态码 = %d, 期望 200 (期望得到业务错误码)", w2.Code)
	}

	var resp2 map[string]interface{}
	err := json.Unmarshal(w2.Body.Bytes(), &resp2)
	if err != nil {
		t.Fatalf("解析第二次响应失败: %v", err)
	}

	// 检查业务错误码是否为 CodeTooManyRequests (10013)
	// 在 controller/response.go 中定义: CodeTooManyRequests
	// JSON 反序列化数字默认为 float64，需要断言转换
	codeVal, ok := resp2["code"].(float64)
	if !ok {
		t.Errorf("未找到响应码 code 或格式非数字: %v", resp2)
	}

	expectedCode := float64(controller.CodeTooManyRequests)
	if codeVal != expectedCode {
		t.Errorf("第二次请求业务码 = %v, 期望 %v (应该被限流)", codeVal, expectedCode)
	} else {
		t.Logf("第二次请求被正确限流，返回业务码: %v", codeVal)
	}

	// [第三次请求]：等待后应该成功
	// 原因：我们睡了 1.1 秒。因为产生速率是 1秒/个，睡这么久足够生成一个新的令牌了。
	time.Sleep(1100 * time.Millisecond)

	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/ping", nil)
	r.ServeHTTP(w3, req3)

	if w3.Code != 200 {
		t.Errorf("第三次请求 http 状态码 = %d, 期望 200", w3.Code)
	}
	var resp3 map[string]interface{}
	json.Unmarshal(w3.Body.Bytes(), &resp3)
	if resp3["message"] != "pong" {
		t.Errorf("第三次请求应该成功(令牌已生成) 但返回了: %v", resp3)
	}
}
