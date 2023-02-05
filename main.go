package main

import (
	"context"
	"fmt"

	"github.com/skhaz/scheduler/controller"
	"github.com/skhaz/scheduler/database"
	"github.com/skhaz/scheduler/repository"
	"github.com/skhaz/scheduler/workflow"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

func main() {
	viper.AutomaticEnv()

	logger, _ := zap.NewDevelopment()

	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
		viper.Get("POSTGRES_HOST"),
		viper.Get("POSTGRES_PORT"),
		viper.Get("POSTGRES_USER"),
		viper.Get("POSTGRES_PASSWORD"),
		viper.Get("POSTGRES_DB"),
	)

	db, err := database.Connect(dsn)
	if err != nil {
		panic(err)
	}

	registry := repository.NewRepositoryRegistry(
		db,
		&repository.TriggerRepository{},
	)

	c := ctrl.GetConfigOrDie()
	clientset := kubernetes.NewForConfigOrDie(c)
	api, err := dynamic.NewForConfig(c)
	if err != nil {
		panic(err)
	}

	var ctx = context.Background()

	wf := workflow.NewWorkflow(ctx, api, clientset)

	server := controller.InitServer()
	server.SetLogger(logger)
	server.SetRepositoryRegistry(registry)
	server.SetWorkflow(wf)
	server.Run()
}
