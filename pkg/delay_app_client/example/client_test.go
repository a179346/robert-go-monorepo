package delay_app_client_test

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/a179346/robert-go-monorepo/pkg/delay_app_client"
	"github.com/a179346/robert-go-monorepo/pkg/roberthttp"
)

func Example() {
	client := delay_app_client.New(
		"http://localhost:8080",
		http.Client{Timeout: 5000 * time.Millisecond},
	)

	ctx := context.Background()

	response, err := client.Delay(ctx, 2000, "Hello, World!")
	if err != nil {
		if errResponse, ok := err.(roberthttp.DefaultResponseError[interface{}]); ok {
			fmt.Printf("errResponse %+v", errResponse)
		} else {
			fmt.Printf("err %v", err)
		}
		return
	}

	fmt.Printf("response %+v", response)
	// Output: response &{Data:Hello, World!}
}
