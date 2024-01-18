package agent

import (
	"context"
	"log"
	"net"

	datanet "github.com/davidroman0O/gogog/data/net"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "agent",
		Short: ".",
		Long:  `.`,
		Run: func(cmd *cobra.Command, args []string) {
			lis, err := net.Listen("tcp", ":50051")
			if err != nil {
				log.Fatalf("failed to listen: %v", err)
			}
			s := grpc.NewServer()
			datanet.RegisterAgentServer(s, &server{})
			log.Println("Server is listening on port 50051...")

			if err := s.Serve(lis); err != nil {
				log.Fatalf("failed to serve: %v", err)
			}
		},
	}
}

type server struct {
	datanet.UnimplementedAgentServer
}

func (s *server) SetAuth(ctx context.Context, in *datanet.UserData) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func grpcServer() chan error
