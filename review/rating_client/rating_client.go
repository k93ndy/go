package rating_client

import (
	"time"

	pb "./pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func GetRate(dstAddr string, timeoutSecond int, reviewId int32) (*pb.Rate, error) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(dstAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer func() {
		closeErr := conn.Close()
		if err == nil {
			// Otherwise, we will return the primary error and ignore the error from Close.
			err = closeErr
		}
	}()
	c := pb.NewRatingClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeoutSecond))
	defer cancel()
	r, err := c.GetRate(ctx, &pb.ReviewInfo{ReviewId: reviewId})
	if err != nil {
		return nil, err
	}

	return r, err
}
