package paraformer

import (
	"context"
	"log"
	"net/http"
	"strings"

	httpclient "github.com/devinyf/dashscopego/httpclient"
	"github.com/google/uuid"
)

// real-time voice recognition

func ConnRecognitionClient(ctx context.Context, request *Request, token string) (*httpclient.WsClient, error) {
	// Initialize the client with the necessary parameters.
	header := http.Header{}
	header.Add("Authorization", token)

	ctx_ws, cancelFn := context.WithCancel(ctx)

	client := httpclient.NewWsClient(ParaformerWSURL, header, ctx_ws, cancelFn)

	if err := client.ConnClient(request); err != nil {
		return nil, err
	}

	return client, nil
}

func CloseRecognitionClient(cli *httpclient.WsClient) error {
	cli.CancelFn()

	if err := cli.CloseClient(); err != nil {
		log.Printf("close client error: %v", err)
		return err
	}

	return nil
}

func SendRadioData(cli *httpclient.WsClient, bytesData []byte) {
	cli.SendBinaryDates(bytesData)
}

type ResultWriter interface {
	WriteResult(str string) error
}

func HandleRecognitionResult(cli *httpclient.WsClient, fn StreamingFunc) {
	outputChan, errChan := cli.ResultChans()

	// TODO: handle errors.
BREAK_FOR:
	for {
		select {
		case output, ok := <-outputChan:
			if !ok {
				log.Println("outputChan is closed")
				break BREAK_FOR
			}

			// streaming callback func
			if err := fn(cli.Ctx, output.Data); err != nil {
				log.Println("error: ", err)
				break BREAK_FOR
			}

		case err := <-errChan:
			if err != nil {
				log.Println("error: ", err)
				break BREAK_FOR
			}
		case <-cli.Ctx.Done():
			cli.Over = true
			log.Println("Done")
			break BREAK_FOR
		}
	}

	log.Println("get recognition result...over")
}

// task_id length 32.
func GenerateTaskID() string {
	u, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	uuid := strings.ReplaceAll(u.String(), "-", "")

	return uuid
}
