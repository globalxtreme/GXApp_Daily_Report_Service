package core

import (
	"fmt"
	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"service/internal/pkg/config"
	"time"
)

type FirstRepository[F any, M any] interface {
	FirstByForm(form F, args ...func(query *gorm.DB) *gorm.DB) M
}

type FindRepository[F any, M any] interface {
	FindByForm(form F) []M
}

type PaginateRepository[F any, M any] interface {
	PaginateByForm(form F) ([]M, interface{})
}

type NumberPoolRepository interface {
	TransactionInterface
	TakenNumberPool(number ...string) string
}

func GetIncrementMonthly(model interface{}) int64 {
	var totalData int64
	config.PgSQL.Unscoped().
		Where("EXTRACT(MONTH FROM \"createdAt\") = ?", time.Now().Month()).
		Model(&model).
		Count(&totalData)

	return totalData + 1
}

func Truncate(db *gorm.DB, tables ...schema.Tabler) {
	if len(tables) > 0 {
		for _, table := range tables {
			err := db.Exec(fmt.Sprintf("truncate table %s restart identity cascade", table.TableName())).Error
			if err != nil {
				xtremepkg.LogError(fmt.Sprintf("Truncate invalid: %v", err), false)
			}
		}
	}
}

func TakenNumberPool(repo NumberPoolRepository, tx *gorm.DB) (string, func()) {
	var numberPool string
	tx.Transaction(func(tx *gorm.DB) error {
		repo.SetTransaction(tx)
		numberPool = repo.TakenNumberPool()

		return nil
	})

	errFunc := func() {
		if r := recover(); r != nil {
			repo.TakenNumberPool(numberPool)
			panic(r)
		}
	}

	return numberPool, errFunc
}
