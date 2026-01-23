package application

import (
	"github.com/rafaeldepontes/imagopher/internal/application/model"
	"github.com/rafaeldepontes/imagopher/internal/image/processor/controller"
)

func NewApplication() *model.Application {
	return &model.Application{
		ImageController: controller.NewController(),
	}
}
