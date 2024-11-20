package nacronsritammarat

import (
	"log"
	"zuck-my-clothe/zuck-my-clothe-backend/platform"
	"zuck-my-clothe/zuck-my-clothe-backend/repository"
	"zuck-my-clothe/zuck-my-clothe-backend/usecases"

	"github.com/robfig/cron/v3"
)

type KonNaCron struct {
	Cron        *cron.Cron
	CronUsecase usecases.KonCronUsecase
}

func SummonKonCron(db *platform.Postgres) KonNaCron {
	c := cron.New()
	paymentRepo := repository.CreateNewPaymentRepository(db)
	orderDetailRepo := repository.CreateOrderDetailRepository(db)
	usecase := usecases.CreateNewKonCronUsecase(paymentRepo,orderDetailRepo)
	scheduler := KonNaCron{Cron: c, CronUsecase: usecase}

	c.AddFunc("@every 1m", func() {
		if err := scheduler.CronUsecase.CleanupExpiredPayment(); err != nil {
			log.Default()
		}
		if err := scheduler.CronUsecase.CleanUpExpiredOrder(); err != nil {
			log.Default()
		}
	})

	return scheduler
}

func (s *KonNaCron) StartKonKron() {
	s.Cron.Start()
}
