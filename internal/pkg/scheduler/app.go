package scheduler

import (
	"context"
	"sync"
)

type App struct {
	configType string
	pathConfig string

	serviceProvider *serviceProvider
}

// NewApp ...
func NewApp(ctx context.Context, configType string, pathConfig string) (*App, error) {
	a := &App{
		configType: configType,
		pathConfig: pathConfig,
	}
	err := a.initDeps(ctx)

	return a, err
}

func (a *App) initDeps(ctx context.Context) error {

	inits := []func(context.Context) error{
		a.initServiceProvider,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider(a.configType, a.pathConfig)

	return nil
}

// Run ...
func (a *App) Run(ctx context.Context) error {
	//nolint:errcheck
	defer func() {
		a.serviceProvider.db.Close()
		a.serviceProvider.rabbitProducer.Close()
	}()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	err := a.runSchedulerService(ctx, wg)
	if err != nil {
		return err
	}
	wg.Wait()
	a.serviceProvider.GetLogger().Info("scheduler service stopped")

	return nil
}

func (a *App) runSchedulerService(ctx context.Context, wg *sync.WaitGroup) error {
	go func() {
		defer wg.Done()

		a.serviceProvider.GetLogger().Info("attempting to run scheduler service")
		a.serviceProvider.GetSchedulerService(ctx).Run(ctx)
	}()

	return nil
}
