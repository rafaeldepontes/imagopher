package application

import (
	"github.com/rafaeldepontes/imagopher/internal/database/postgres"
	"github.com/rafaeldepontes/imagopher/internal/images/processor"
	"github.com/rafaeldepontes/imagopher/internal/images/processor/controller"
)

type Application struct {
	ImageController processor.Controller
}

func NewApplication() Application {
	return Application{
		ImageController: controller.NewController(),
	}
}

func (app Application) Shutdown() error {
	if err := postgres.Close(); err != nil {
		return err
	}
	return nil
}
