package main

import (
	"context"
	"fmt"
	templatesV1Pb "github.com/artrey/ago-grpc-templates/pkg/api/proto/v1"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"gopkg.in/alecthomas/kingpin.v2"
	"net"
	"os"
	"time"
)

type CLI struct {
}

var (
	app = kingpin.New("client", "A command-line client to templates grpc server.")

	host = app.Flag("host", "Host of grpc server.").Default("localhost").String()
	port = app.Flag("port", "Port of grpc server.").Default("9999").String()

	create      = app.Command("create", "Create new template.")
	createTitle = create.Arg("title", "Title of template.").Required().String()
	createPhone = create.Arg("phone", "Phone for template.").Required().String()

	list = app.Command("list", "List templates.")

	get   = app.Command("get", "Get template by ID.")
	getId = get.Arg("id", "Id of template.").Required().Int64()

	update      = app.Command("update", "Update template by ID.")
	updateId    = update.Arg("id", "Id of template.").Required().Int64()
	updateTitle = update.Arg("title", "New title of template.").Required().String()
	updatePhone = update.Arg("phone", "New phone for template.").Required().String()

	remove   = app.Command("remove", "Remove template by ID.")
	removeId = remove.Arg("id", "Id of template.").Required().Int64()
)

func (cli *CLI) Run() (err error) {
	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))

	conn, err := grpc.Dial(net.JoinHostPort(*host, *port), grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			if err == nil {
				err = cerr
			}
		}
	}()

	client := templatesV1Pb.NewTemplatesServiceClient(conn)
	ctx, _ := context.WithTimeout(context.Background(), time.Second)

	switch cmd {
	case create.FullCommand():
		response, err := client.Create(ctx, &templatesV1Pb.TemplateCreateRequest{
			Title: *createTitle,
			Phone: *createPhone,
		})
		if err != nil {
			return err
		}
		fmt.Printf("Template created: %+v\n", response)

	case list.FullCommand():
		response, err := client.List(ctx, &empty.Empty{})
		if err != nil {
			return err
		}
		fmt.Println("List of templates:")
		for _, t := range response.Items {
			fmt.Printf("%+v\n", t)
		}

	case get.FullCommand():
		response, err := client.GetById(ctx, &templatesV1Pb.TemplateByIdRequest{Id: *getId})
		if err != nil {
			return err
		}
		fmt.Printf("Requested template: %+v\n", response)

	case update.FullCommand():
		response, err := client.UpdateById(ctx, &templatesV1Pb.TemplateUpdateRequest{
			Id:    *updateId,
			Title: *updateTitle,
			Phone: *updatePhone,
		})
		if err != nil {
			return err
		}
		fmt.Printf("Updated template: %+v\n", response)

	case remove.FullCommand():
		response, err := client.DeleteById(ctx, &templatesV1Pb.TemplateByIdRequest{Id: *removeId})
		if err != nil {
			return err
		}
		fmt.Printf("Removed template: %+v\n", response)
	}

	return nil
}
