package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"github.com/desperadochn/Go-000/Week02/pkg/dao"
	"github.com/desperadochn/Go-000/Week02/pkg/endpoint"
	"github.com/desperadochn/Go-000/Week02/pkg/redis"
	"github.com/desperadochn/Go-000/Week02/pkg/service"
	"github.com/desperadochn/Go-000/Week02/pkg/router"
)

func main() {

	var (
		// 服务地址和服务名
		servicePort = flag.Int("service.port", 8085, "service port")
	)

	flag.Parse()

	ctx := context.Background()
	errChan := make(chan error)

	err := dao.MysqlInit("127.0.0.1", "3306", "root", "123456", "user")
	if err != nil {
		log.Fatal(err)
	}

	err = redis.InitRedis("127.0.0.1", "6379", "")
	if err != nil {
		log.Fatal(err)
	}

	userService := service.MakeUserServiceImpl(&dao.UserDAOImpl{})

	userEndpoints := &endpoint.UserEndpoints{
		endpoint.MakeRegisterEndpoint(userService),
		endpoint.MakeLoginEndpoint(userService),
	}

	r := router.MakeHttpHandler(ctx, userEndpoints)

	go func() {
		errChan <- http.ListenAndServe(":"+strconv.Itoa(*servicePort), r)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	error := <-errChan
	log.Println(error)

}
