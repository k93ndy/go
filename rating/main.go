package main

import (
    "log"
    "net"

    pb "./pb"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    "golang.org/x/net/context"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
    "github.com/spf13/viper"

    "github.com/sirupsen/logrus"
    "github.com/grpc-ecosystem/go-grpc-middleware"
    "github.com/grpc-ecosystem/go-grpc-middleware/tags"
    "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
)

type Rate struct {
    gorm.Model
    ReviewId    int32
    Maximum     int32
    Current     int32
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
}

func (s *server) GetRate(ctx context.Context, in *pb.ReviewInfo) (*pb.Rate, error) {
    db, err := s.gormConnect()
    if (err != nil) {
        panic(err.Error())
    }
    defer db.Close()

    var rate Rate
    db.First(&rate, in.ReviewId)

    return &pb.Rate {
        ReviewId:    rate.ReviewId,
        Maximum:     rate.Maximum,
        Current:     rate.Current,
    }, nil
}

// server side interceptor
//func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
//    log.Printf("before handling. Info: %+v", info)
//    resp, err := handler(ctx, req)
//    log.Printf("after handling. resp: %+v", resp)
//    return resp, err
//}


func main() {
    //db, err := gormConnect()
    //if (err != nil) {
    //    panic(err.Error())
    //}
    //defer db.Close()

    //db.AutoMigrate(&Rate{})
    //db.Create(&Rate{ReviewId: 1, Maximum: 100, Current:70})

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

    //s := grpc.NewServer(grpc.UnaryInterceptor(UnaryServerInterceptor))
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

    pb.RegisterRatingServer(s, &server{
        dbType: viper.GetString("db.type"),
        dbUser: viper.GetString("db.user"),
        dbPassword: viper.GetString("db.password"),
        dbProtocol: viper.GetString("db.protocol"),
        dbName: viper.GetString("db.name"),
    })
    // Register reflection service on gRPC server.
    reflection.Register(s)
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}

