package reporter

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/lgrisa/lib/utils/timeutil"
	"github.com/sirupsen/logrus"
	"go.uber.org/atomic"
	"net/http"
	"strings"
	"time"
)

const (
	serverType               = "server"
	serverDevelopmentVersion = "dev"
)

var (
	defaultErrorReportedHook *errorReportedHook
	redisExpiresInterval     = timeutil.Day * 2
	bigCacheValue            = []byte("true")
	redisValue               = "true"
)

func InitErrorReporter() {
	if startconfig.StartConfig.ErrReporterToken == "" {
		return
	}

	if defaultErrorReportedHook != nil {
		return
	}

	defaultErrorReportedHook = newErrorReportedHook()

	// 非开发版本才hook
	if build.GetBuildGitTag() != serverDevelopmentVersion {
		logrus.AddHook(defaultErrorReportedHook)
	}
}

func SetEnableReport(enable bool) {
	if defaultErrorReportedHook != nil {
		defaultErrorReportedHook.enableReport.Store(enable)
	}
}

// 记录和上报飞书
func OnErrorReported(errorType, version, content string) {
	if defaultErrorReportedHook == nil {
		return
	}

	defaultErrorReportedHook.OnErrorReported(errorType, version, content)
}

type errorReportedHook struct {
	name  string
	token string

	redisClient redis.UniversalClient

	levels []logrus.Level

	errorMd5Cache *bigcache.BigCache

	// 服务器启动成功之后开始上报
	enableReport *atomic.Bool
}

func newErrorReportedHook() *errorReportedHook {
	// 创建redisClient
	redisConf := startconfig.StartConfig.Redis
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    []string{redisConf.Addr},
		Password: redisConf.Password,
		DB:       redisConf.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	ping := client.Ping(ctx)
	if err := ping.Err(); err != nil {
		logrus.WithError(err).Panicf("测试连接到redis失败: %s", redisConf.Addr)
	}

	logrus.Infof("连接到redis成功: %s", redisConf.Addr)

	errorMd5Cache, err := bigcache.NewBigCache(bigcache.Config{
		Shards:           1024,
		HardMaxCacheSize: 50,
		MaxEntrySize:     1024,
		LifeWindow:       redisExpiresInterval + time.Minute,
		CleanWindow:      timeutil.Day,
	})
	if err != nil {
		logrus.WithError(err).Panicf("创建错误上报md5缓存失败")
	}

	hook := &errorReportedHook{
		name:          hostname.GetValue(),
		token:         startconfig.StartConfig.ErrReporterToken,
		levels:        []logrus.Level{logrus.ErrorLevel},
		redisClient:   client,
		errorMd5Cache: errorMd5Cache,
		enableReport:  atomic.NewBool(false),
	}

	return hook
}

func (hook *errorReportedHook) Fire(entry *logrus.Entry) error {
	if entry.Level != logrus.ErrorLevel {
		return nil
	}

	if !hook.enableReport.Load() {
		return nil
	}

	content := entry.Message
	for name, msg := range entry.Data {
		if err, ok := msg.(error); ok {
			errStack := fmt.Sprintf("%v:%v", name, err)
			content = content + "\n" + errStack
		}
	}

	hook.OnErrorReported(serverType, build.GetBuildGitTag(), content)
	return nil
}

func (hook *errorReportedHook) Levels() []logrus.Level {
	return hook.levels
}

func (hook *errorReportedHook) OnErrorReported(errorType, version, content string) {
	asyncall.AsyncCall(func() {
		if !defaultErrorReportedHook.errorReported(content) {
			return
		}

		// 成功记录，上报飞书
		defaultErrorReportedHook.reportFeishuText(errorType, version, content)
	})
}

func (hook *errorReportedHook) errorReported(content string) bool {
	// 先在缓存中判断一下
	if _, err := hook.errorMd5Cache.Get(content); err == nil || err != bigcache.ErrEntryNotFound {
		// 已缓存
		return false
	}
	hook.errorMd5Cache.Set(content, bigCacheValue)

	// 使用错误内容生成一个md5key
	md5key := fmt.Sprintf("ErrorReported_%v", md5.String([]byte(content)))

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := hook.redisClient.SetNX(ctx, md5key, redisValue, redisExpiresInterval).Result()
	if err != nil {
		logrus.WithError(err).Info("错误上报设置redis key 错误")
		return false
	}

	return result
}

func (hook *errorReportedHook) reportFeishuText(errorType, version, content string) {
	sb := strings.Builder{}
	sb.WriteString("🔴 ErrorReported")
	sb.WriteString("\n")
	sb.WriteString(hook.name)
	sb.WriteString("\n")
	sb.WriteString("\n")
	sb.WriteString("errorType:")
	sb.WriteString("\n")
	sb.WriteString(errorType)
	sb.WriteString("\n")
	sb.WriteString("\n")
	sb.WriteString("version:")
	sb.WriteString("\n")
	sb.WriteString(version)
	sb.WriteString("\n")
	sb.WriteString("\n")
	sb.WriteString("content:")
	sb.WriteString("\n")
	sb.WriteString(content)

	msg := &message{}
	msg.MsgType = "text"
	msg.Content.Text = sb.String()

	postData, err := jsoniter.MarshalToString(msg)
	if err != nil {
		logrus.WithError(err).Info("错误上报转json失败")
		return
	}

	if _, err = http.Post(hook.token, "application/json", strings.NewReader(postData)); err != nil {
		logrus.WithError(err).Info("错误上报到飞书失败")
		return
	}
}
