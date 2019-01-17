package main

import (
    "log"
    "./rating_client"

    "net"
    pb "./pb"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    "golang.org/x/net/context"
    "google.golang.org/grpc"
    "google.golang.org/grpc/metadata"
    "google.golang.org/grpc/reflection"
    "github.com/spf13/viper"

    "github.com/sirupsen/logrus"
    "github.com/grpc-ecosystem/go-grpc-middleware"
    "github.com/grpc-ecosystem/go-grpc-middleware/tags"
    "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
)

type Review struct {
    gorm.Model
    ProductId       int32
    ReviewerName    string
    Content         string
}

func (s *server) gormConnect() (*gorm.DB, error) {
    CONNECT := s.dbUser + ":" + s.dbPassword + "@" + s.dbProtocol + "/" + s.dbName + "?parseTime=true"
    db, err := gorm.Open(s.dbType, CONNECT)
    
    return db, err
}

type server struct{
    dbType string
    dbUser string
    dbPassword string
    dbProtocol string
    dbName string
    ratingEndpoint string
    ratingTimeout int
}

func buildReviewMessages(reviews []Review, ratingService string, timeoutSecond int, tracingHeaders map[string]string) (*pb.ReviewMessages, error) {
    var msgs []*pb.ReviewMessage

    //bad implementation, should use a bulk get
    for _, review := range reviews {
        result, err := rating_client.GetRate(ratingService, timeoutSecond, int32(review.ID), tracingHeaders)
        if err != nil {
            log.Println(err.Error())
        }
        var rate *pb.Rate
        if err != nil {
            rate = nil
        } else {
            rate = &pb.Rate{
                Maximum: result.Maximum,
                Current: result.Current,
            }
        }
        msgs = append(msgs, &pb.ReviewMessage{
            ReviewId: int32(review.ID),
            ProductId: review.ProductId,
            ReviewerName: review.ReviewerName,
            Content: review.Content,
            Rate: rate,
        })
    }

    return &pb.ReviewMessages{
        ReviewMessages: msgs,
    }, nil
}

//extract headers used for tracing
func extractTracingHeaders(ctx *context.Context) (tracingHeaders map[string]string) {
    targetHeaders := []string{
        "x-request-id",
        "x-b3-traceid",
        "x-b3-spanid",
        "x-b3-parentspanid",
        "x-b3-sampled",
        "x-b3-flags",
        "x-ot-span-context",
    }
    tracingHeaders = make(map[string]string)
    md, ok := metadata.FromIncomingContext(*ctx)
    if ok {
        for _, value := range targetHeaders {
            _, ok := md[value]
            if ok {
                tracingHeaders[value] = md[value][0]
            }
        }
    }
    return
}

func (s *server) GetMostHelpfulReviews(ctx context.Context, in *pb.ProductInfo) (*pb.ReviewMessages, error) {
    db, err := s.gormConnect()
    if (err != nil) {
        log.Println(err.Error())
    }
    defer db.Close()

    var reviews []Review
    db.Where("product_id = ?", in.ProductId).Find(&reviews)
   
    result, err := buildReviewMessages(reviews, s.ratingEndpoint, s.ratingTimeout, extractTracingHeaders(&ctx))
    if (err != nil) {
        log.Println(err.Error())
    }
    return result, err

    //var msgs []*pb.ReviewMessage
    //msgs = append(msgs, &pb.ReviewMessage{
    //    ReviewId: 1,
    //    ProductId: 1,
    //    ReviewerName: "first reviewer",
    //    Content: "first review",
    //    Rate: &pb.Rate{
    //        Maximum: 100,
    //        Current: 80,
    //    },
    //})
    //return &pb.ReviewMessages{
    //    ReviewMessages: msgs,
    //}, nil
}

//const (
//    srvPort = ":7100"
//)

func main() {
    //db, err := gormConnect()
    //if (err != nil) {
    //    panic(err.Error())
    //}
    //defer db.Close()

    //db.AutoMigrate(&Review{})
    //insert record for test use
    //db.Create(&Review{ProductId: 2, ReviewerName: "Jpi&dRtx", Content: "you should never buy it."})

    logrusLogger := logrus.New()
    logrusLogger.SetFormatter(&logrus.JSONFormatter{})
    logrusEntry := logrus.NewEntry(logrusLogger)
    customFunc := grpc_logrus.DefaultCodeToLevel
    opts := []grpc_logrus.Option{
        grpc_logrus.WithLevels(customFunc),
    }
    grpc_logrus.ReplaceGrpcLogger(logrusEntry)

    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    err := viper.ReadInConfig()
    if err != nil {
        panic(err.Error())
    }

    lis, err := net.Listen("tcp", viper.GetString("srvPort"))
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    // s := grpc.NewServer()
    s := grpc.NewServer(
        grpc_middleware.WithUnaryServerChain(
            grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
            grpc_logrus.UnaryServerInterceptor(logrusEntry, opts...),
        ),
    //    grpc_middleware.WithStreamServerChain(
    //        grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
    //        grpc_logrus.StreamServerInterceptor(logrusEntry, opts...),
    //    ),
    )

    pb.RegisterReviewServer(s, &server{
        dbType: viper.GetString("db.type"),
        dbUser: viper.GetString("db.user"),
        dbPassword: viper.GetString("db.password"),
        dbProtocol: viper.GetString("db.protocol"),
        dbName: viper.GetString("db.name"),
        ratingEndpoint: viper.GetString("rating.endpoint"),
        ratingTimeout: viper.GetInt("rating.timeout"),
    })
    // Register reflection service on gRPC server.
    reflection.Register(s)
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }

    //rate, err := rating_client.GetRate("localhost:8000", 2, 1)
    //if err != nil {
    //    panic(err)
    //}
    //log.Println(rate.ReviewId, rate.Maximum, rate.Current)
}

