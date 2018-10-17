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

type Advertisement struct {
    gorm.Model
    Name        string `gorm:"type:varchar(100);UNIQUE"`
    Description string `gorm:"type:varchar(255)"`
    Url         string `gorm:"type:varchar(255)"`
    Image       []byte `gorm:"size:70000"`
}

func (s *server) gormConnect() *gorm.DB {

    CONNECT := s.dbUser + ":" + s.dbPassword + "@" + s.dbProtocol + "/" + s.dbName + "?parseTime=true"
    db, err := gorm.Open(s.dbType, CONNECT)

    if err != nil {
        panic(err.Error())
    }
    return db
}

type server struct{
    dbType string
    dbUser string
    dbPassword string
    dbProtocol string
    dbName string
}

func (s *server) GetRandomAdvertisement(ctx context.Context, in *pb.Empty) (*pb.AdContent, error) {
    db := s.gormConnect()
    defer db.Close()

    //return random advertisement
    var ad Advertisement
    db.Order(gorm.Expr("rand()")).First(&ad)

    return &pb.AdContent{
        Name:        ad.Name,
        Description: ad.Description,
        Url:         ad.Url,
        Image:       ad.Image,
    }, nil
}

func main() {
    // db := gormConnect()
    // defer db.Close()

    // // Migrate the schema
    // db.AutoMigrate(&Advertisement{})

    // // Create
    // db.Create(&Advertisement{Name: "Test2", Description: "Description2", Url: "#", Image: []byte("Image2")})
    // db.Create(&Advertisement{Name: "Test3", Description: "Description3", Url: "#", Image: []byte("Image3")})

    // // Read
    // var ad Advertisement
    // db.First(&ad, 1)                   // find product with id 1
    // db.First(&ad, "name = ?", "Test1") // find product with name Test1

    // // Update
    // db.Model(&ad).Update("Name", "Test1-updated")

    // Delete
    // db.Delete(&ad)

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

    pb.RegisterAdvertisementServer(s, &server{
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

