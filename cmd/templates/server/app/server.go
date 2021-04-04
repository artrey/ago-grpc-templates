package app

import (
	"context"
	templatesV1Pb "github.com/artrey/ago-grpc-templates/pkg/api/proto/v1"
	"github.com/artrey/ago-grpc-templates/pkg/templates"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
)

type Server struct {
	templatesV1Pb.UnimplementedTemplatesServiceServer
	templatesSvc *templates.Service
}

func NewServer(templatesSvc *templates.Service) *Server {
	return &Server{templatesSvc: templatesSvc}
}

func (s *Server) Create(ctx context.Context, request *templatesV1Pb.TemplateCreateRequest) (*templatesV1Pb.TemplateResponse, error) {
	template, err := s.templatesSvc.Create(ctx, request.Title, request.Phone)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return serviceTemplateToGRPCResponse(template), nil
}

func (s *Server) List(ctx context.Context, empty *emptypb.Empty) (*templatesV1Pb.TemplatesResponse, error) {
	allTemplates, err := s.templatesSvc.List(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	items := make([]*templatesV1Pb.TemplateResponse, len(allTemplates))
	for i, template := range allTemplates {
		items[i] = serviceTemplateToGRPCResponse(template)
	}

	return &templatesV1Pb.TemplatesResponse{Items: items}, nil
}

func (s *Server) GetById(ctx context.Context, request *templatesV1Pb.TemplateByIdRequest) (*templatesV1Pb.TemplateResponse, error) {
	template, err := s.templatesSvc.Get(ctx, request.Id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return serviceTemplateToGRPCResponse(template), nil
}

func (s *Server) UpdateById(ctx context.Context, request *templatesV1Pb.TemplateUpdateRequest) (*templatesV1Pb.TemplateResponse, error) {
	template, err := s.templatesSvc.Update(ctx, request.Id, request.Title, request.Phone)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return serviceTemplateToGRPCResponse(template), nil
}

func (s *Server) DeleteById(ctx context.Context, request *templatesV1Pb.TemplateByIdRequest) (*templatesV1Pb.TemplateResponse, error) {
	template, err := s.templatesSvc.Delete(ctx, request.Id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return serviceTemplateToGRPCResponse(template), nil
}

func serviceTemplateToGRPCResponse(template *templates.Template) *templatesV1Pb.TemplateResponse {
	return &templatesV1Pb.TemplateResponse{
		Id:    template.Id,
		Title: template.Title,
		Phone: template.Phone,
		Created: &timestamppb.Timestamp{
			Seconds: template.Created,
			Nanos:   0,
		},
		Updated: &timestamppb.Timestamp{
			Seconds: template.Updated,
			Nanos:   0,
		},
	}
}
