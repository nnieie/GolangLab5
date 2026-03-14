package tracer

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

var (
	metricExportInterval = 5 * time.Second

	// 请求计数器
	RequestCounter metric.Int64Counter

	// 请求延迟直方图
	RequestDuration metric.Float64Histogram

	// 错误计数器
	ErrorCounter metric.Int64Counter

	// 互动服务
	InteractionLikeCounter    metric.Int64Counter // 点赞/取消赞计数 (Tags: action="like"|"unlike")
	InteractionCommentCounter metric.Int64Counter // 评论发布数 (Tags: action="add"|"delete")

	// 视频服务
	VideoPublishCounter metric.Int64Counter // 视频投稿量

	// 聊天服务
	ChatMessageCounter metric.Int64Counter // 消息发送量

	// 用户服务
	UserRegisterCounter metric.Int64Counter // 新用户注册量
	UserLoginCounter    metric.Int64Counter // 用户登录成功/失败情况 (Tags: status="success"|"fail")

	// Kafka
	MQProduceCounter metric.Int64Counter // 消息发送成功/失败次数 (Tags: topic, status="success"|"fail")
	MQConsumeCounter metric.Int64Counter // 消费者处理成功/失败次数 (Tags: topic, status="success"|"fail", error_type)
)

// InitMetrics 初始化 Metrics
func InitMetrics(serviceName string, collectorAddr string) (func(context.Context) error, error) {
	ctx := context.Background()

	// 创建 Exporter
	exporter, err := otlpmetrichttp.New(ctx,
		otlpmetrichttp.WithEndpoint(collectorAddr),
		otlpmetrichttp.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create metric exporter: %w", err)
	}

	// 创建 Resource
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion("1.0.0"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// 定义直方图分桶 View
	durationView := sdkmetric.NewView(
		sdkmetric.Instrument{Name: "http_request_duration_seconds"},
		sdkmetric.Stream{
			Aggregation: sdkmetric.AggregationExplicitBucketHistogram{
				Boundaries: []float64{0.005, 0.01, 0.02, 0.035, 0.05, 0.075, 0.1, 0.25, 0.5, 1, 2.5, 5},
			},
		},
	)

	// 创建 MeterProvider
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(metricExportInterval))),
		sdkmetric.WithResource(res),
		sdkmetric.WithView(durationView),
	)

	otel.SetMeterProvider(mp)

	// 获取 Meter
	meter := mp.Meter(serviceName)

	// 创建指标
	RequestCounter, err = meter.Int64Counter(
		"http_requests_total",
		metric.WithDescription("Total number of HTTP requests"),
	)
	if err != nil {
		return nil, err
	}

	RequestDuration, err = meter.Float64Histogram(
		"http_request_duration_seconds",
		metric.WithDescription("HTTP request duration in seconds"),
	)
	if err != nil {
		return nil, err
	}

	ErrorCounter, err = meter.Int64Counter(
		"http_errors_total",
		metric.WithDescription("Total number of HTTP errors"),
	)
	if err != nil {
		return nil, err
	}

	// 业务指标初始化

	InteractionLikeCounter, err = meter.Int64Counter(
		"biz_interaction_like_total",
		metric.WithDescription("Total number of video likes/unlikes"),
	)
	if err != nil {
		return nil, err
	}
	InteractionLikeCounter.Add(ctx, 0, metric.WithAttributes(attribute.String("action", "like")))
	InteractionLikeCounter.Add(ctx, 0, metric.WithAttributes(attribute.String("action", "unlike")))

	InteractionCommentCounter, err = meter.Int64Counter(
		"biz_interaction_comment_total",
		metric.WithDescription("Total number of published/deleted comments"),
	)
	if err != nil {
		return nil, err
	}
	InteractionCommentCounter.Add(ctx, 0, metric.WithAttributes(attribute.String("action", "add")))
	InteractionCommentCounter.Add(ctx, 0, metric.WithAttributes(attribute.String("action", "delete")))

	VideoPublishCounter, err = meter.Int64Counter(
		"biz_video_publish_total",
		metric.WithDescription("Total number of published videos"),
	)
	if err != nil {
		return nil, err
	}

	ChatMessageCounter, err = meter.Int64Counter(
		"biz_chat_message_sent_total",
		metric.WithDescription("Total number of chat messages sent"),
	)
	if err != nil {
		return nil, err
	}

	UserRegisterCounter, err = meter.Int64Counter(
		"biz_user_register_total",
		metric.WithDescription("Total number of new user registrations"),
	)
	if err != nil {
		return nil, err
	}

	UserLoginCounter, err = meter.Int64Counter(
		"biz_user_login_total",
		metric.WithDescription("Total number of user logins (success or fail)"),
	)
	if err != nil {
		return nil, err
	}

	MQProduceCounter, err = meter.Int64Counter(
		"biz_mq_produce_total",
		metric.WithDescription("Total number of Kafka messages produced"),
	)
	if err != nil {
		return nil, err
	}

	MQConsumeCounter, err = meter.Int64Counter(
		"biz_mq_consume_total",
		metric.WithDescription("Total number of Kafka messages consumed"),
	)
	if err != nil {
		return nil, err
	}

	return mp.Shutdown, nil
}
