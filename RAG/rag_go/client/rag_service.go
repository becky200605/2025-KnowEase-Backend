package rag

import (
	"context"
	"flag"
	"fmt"
	"time"

	pb "KnowEase/RAG/rag_go/rag/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RagService struct {
}

func NewRAGService() *RagService {
	return &RagService{}
}

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func (rs *RagService) connect() (*grpc.ClientConn, error) {
	flag.Parse()
	// 建立连接
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect:%v", err)
	}
	return conn, nil
}
func (rs *RagService) InitQuestion(question string) (string, error) {

	conn, err := rs.connect()
	if err != nil {
		return "", err
	}
	defer conn.Close()

	client := pb.NewRAGServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()
	//syncResp, err := client.SyncData(ctx, &pb.SyncDataRequest{FromId: 0})
	//fmt.Printf("Synced %d documents\n", syncResp.SyncedCount)
	searchResp, err := client.Search(ctx, &pb.SearchRequest{
		Query: question,
	})

	if err != nil {
		return "", fmt.Errorf("could not search: %v", err)
	}
	//var PostIDs []string
	return searchResp.Answer, nil
}

func (rs *RagService) SyncQuestion(formid int64) error {
	conn, err := rs.connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewRAGServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	syncResp, err := client.SyncData(ctx, &pb.SyncDataRequest{FromId: formid})
	if err != nil {
		return fmt.Errorf("could not sync: %v", err)
	}
	fmt.Printf("Synced %d documents\n", syncResp.SyncedCount)
	return nil
}
