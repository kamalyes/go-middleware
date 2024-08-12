/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-12 23:28:08
 * @FilePath: \go-middleware\record\access\access_handler.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package access

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kamalyes/go-core/global"
	"github.com/kamalyes/go-core/jwt"
	"github.com/kamalyes/go-toolbox/ip"
	"github.com/mssola/user_agent"
	"go.uber.org/zap"
)

// parseUserAgent
/**
 *  @Description: 解析userAgent
 *  @param c
 *  @param record
 */
func parseUserAgent(c *gin.Context, record *AccessRecordModel) {
	// 获取userAgent
	userAgent := c.Request.UserAgent()
	ua := user_agent.New(c.Request.UserAgent())
	rex := regexp.MustCompile(`\(([^)]+)\)`)
	params := rex.FindAllStringSubmatch(userAgent, -1)
	param := strings.Replace(params[0][0], ")", "", 1)
	uaInfo := strings.Split(param, ";")
	engineName, engineVersion := ua.Engine()
	browserName, browserVersion := ua.Browser()

	record.UserAgent = userAgent
	record.Platform = ua.Platform()
	record.OS = ua.OS()
	record.Engine = engineName + "/" + engineVersion
	record.BrowserName = browserName
	record.BrowserVersion = browserVersion
	record.Brand = ""
	record.ProductModel = strings.TrimSpace(uaInfo[2])
}

// parseBody
/**
 *  @Description: 获取Body
 *  @param c
 *  @param record
 */
func parseBody(c *gin.Context, record *AccessRecordModel) {
	// 获取Request.Body
	body, _ := io.ReadAll(c.Request.Body)
	// 将其转为String
	data := string(body)
	// 替换字符串及字符串分隔
	data = strings.Replace(data, "\n", "", -1)
	dataLs := strings.Split(data, "\r")
	record.Body = data
	if len(dataLs) == 0 {
		return
	}
	// 拼接请求参数
	var build strings.Builder
	for i, v := range dataLs {
		if strings.Contains(v, "Content-Disposition") {
			build.WriteString(v)
			if !strings.Contains(v, "filename") && len(dataLs) >= i+2 {
				build.WriteString("; value=")
				build.WriteString(dataLs[i+2])
			}
			build.WriteString("\n")
		}
	}
	data = build.String()
	return
}

// parseIp
/**
 *  @Description: 解析Ip地址信息
 *  @param c
 *  @param record
 */
func parseIp(c *gin.Context, record *AccessRecordModel) error {
	// 获取高德接口地址、Key和数字签名
	configKeys := []string{"amap.url", "amap.key", "amap.sign"}
	configMap := make(map[string]string)
	for _, key := range configKeys {
		value := strings.TrimSpace(global.VP.GetString(key))
		if value == "" {
			global.LOG.Error("配置项 " + key + " 未配置")
			return errors.New("缺少高德接口配置信息")
		} else {
			configMap[key] = value
		}
	}

	record.Ip = c.ClientIP()
	data, err := ip.GetGPSByIpAmap(
		configMap["amap.key"],
		configMap["amap.sign"],
		configMap["amap.url"],
		record.Ip)
	if err != nil {
		return nil
	}
	if data["country"] != nil {
		record.Country = data["country"].(string)
	}
	if data["location"] != nil {
		record.Country = data["location"].(string)
	}
	if data["province"] != nil {
		record.Country = data["province"].(string)
	}
	if data["city"] != nil {
		record.Country = data["city"].(string)
	}
	return nil
}

// parseToken
/**
 *  @Description: 解析Token
 *  @param c
 *  @param record
 */
func parseToken(c *gin.Context, record *AccessRecordModel) error {
	token := c.Request.Header.Get("ACCESS_TOKEN")
	if token != "" {
		j := jwt.NewJWT()
		// 解析token包含的信息
		claims, err := j.ResolveToken(token)
		if err == nil {
			record.UserId = claims.UserId
			record.UserName = claims.UserName
			record.MerchantNo = claims.MerchantNo
			record.AppProductId = claims.AppProductId
			record.PlatformType = claims.PlatformType
		}
	}
	return nil
}

// getHeader
/**
 *  @Description: 获取header
 *  @param c
 *  @return headerStr
 */
func getHeader(c *gin.Context) (headerStr string) {
	// 获取Herader
	header := c.Request.Header
	if len(header) == 0 {
		return
	}
	// 拼接Herder
	var build strings.Builder
	for k, v := range header {
		build.WriteString(k)
		build.WriteString(":")
		for i, v0 := range v {
			build.WriteString(v0)
			if i != len(v)-1 {
				build.WriteString(",")
			}
		}
		build.WriteString(";\n")
	}
	headerStr = build.String()
	return
}

// AccessRecordMiddleware 访问记录 retainDays 请求记录保留的时间
func AccessRecordMiddleware(retainDays int) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body []byte
		record := AccessRecordModel{}
		httpMethod := c.Request.Method
		record.Path = c.Request.URL.Path
		// ContentType
		record.ContentType = c.ContentType()
		// 解析userAgent
		parseUserAgent(c, &record)
		// 解析Body
		parseBody(c, &record)
		// IP
		parseIp(c, &record)
		// Token
		parseToken(c, &record)
		if strings.Contains(c.ContentType(), "multipart/form-data") {
			// 如果该请求是文件上传不记录请求体
			record.Body = "文件上传"
		} else {
			record.Body = string(body)
		}
		writer := responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer
		now := time.Now()
		c.Next()
		record.Error = c.Errors.ByType(gin.ErrorTypePrivate).String()
		record.Status = c.Writer.Status()
		// 延时，毫秒
		record.Latency = time.Now().Sub(now).Milliseconds()
		record.Response = writer.body.String()
		// 保存基础信息
		record.ID = global.CreateId()
		record.CreateTime = global.CreateTime()
		if httpMethod == http.MethodGet {
			var str bytes.Buffer
			m, _ := json.Marshal(record)
			_ = json.Indent(&str, m, "", "    ")
			global.LOG.Info("访问记录" + str.String())
		}
		if err := AccessRecordServiceApp.CreateAccessRecord(record, retainDays); err != nil {
			global.LOG.Error("create access record error:", zap.Any("err", err))
		}
	}
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
