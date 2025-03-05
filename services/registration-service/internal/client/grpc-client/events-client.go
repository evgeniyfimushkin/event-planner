package grpcclient

import (
	"context"
	"log/slog"
	"time"

	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/protos/events"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type EventClient struct {
    api events.EventServiceClient
    log *slog.Logger
}

func NewEventClient(
    ctx context.Context,
    log *slog.Logger,
    addr string,
    timeout time.Duration,
    retriesCount int,
) (*EventClient, error) {

    retryOpts := []retry.CallOption{
        retry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
        retry.WithMax(uint(retriesCount)),
        retry.WithPerRetryTimeout(timeout),
    }

    logOpts := []logging.Option{
        logging.WithLogOnEvents(logging.PayloadReceived, logging.PayloadSent),
    }

    cc, err := grpc.DialContext(ctx, addr,
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithChainUnaryInterceptor(
            logging.UnaryClientInterceptor(InterceptorLogger(log), logOpts...),
            retry.UnaryClientInterceptor(retryOpts...),
        ),
    )
    if err != nil {
        return nil, err
    }
    return &EventClient{
        api: events.NewEventServiceClient(cc),
    }, nil
}

func InterceptorLogger(l *slog.Logger) logging.Logger {
    return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any){
        l.Log(ctx, slog.Level(lvl), msg, fields...)
    })
}

func (c *EventClient) CheckAndReserve(ctx context.Context, eventID uint32) (*events.CheckAndReserveResponse, error) {
    resp, err := c.api.CheckAndReserve(ctx, &events.CheckAndReserveRequest{
        EventId: eventID,
    })
    if err != nil {
        c.log.Error("failed to call CheckAndReserve", "error", err)
        return nil, err
    }
    return resp, nil
}

