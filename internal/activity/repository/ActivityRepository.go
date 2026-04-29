package repository

import (
	xtrememodel "github.com/globalxtreme/go-core/v2/model"
	"gorm.io/gorm"
	"net/url"
	"service/internal/pkg/config"
	"service/internal/pkg/core"
	"service/internal/pkg/model"
)

/** --- INTERFACE --- */

type ActivityRepository interface {
	Find(parameters url.Values) ([]model.Activity, interface{}, error)
}

func NewActivityRepository() ActivityRepository {
	return activityRepository{}
}

/** --- MAIN REPOSITORY --- */

type activityRepository struct {
	Transaction *gorm.DB
}

func (repo activityRepository) Find(parameters url.Values) ([]model.Activity, interface{}, error) {
	query := repo.filterByParam(parameters)
	activities, pagination, err := xtrememodel.Paginate(query.Order("id DESC"), parameters, model.Activity{})
	if err != nil {
		return nil, nil, err
	}

	return activities, pagination, nil
}

/** --- UNEXPORTED FUNCTIONS --- */

func (repo activityRepository) filterByParam(parameters url.Values) *gorm.DB {
	fromDate, toDate := core.SetDateRangeFromAPI(parameters)

	query := config.PgSQL.Where(`"createdAt" BETWEEN ? AND ?`, fromDate, toDate)

	if feature := parameters.Get("feature"); len(feature) > 0 {
		query = query.Where("feature = ?", feature)
	}

	if action := parameters.Get("action"); len(action) > 0 {
		query = query.Where("action = ?", action)
	}

	if searchReq := parameters.Get("search"); len(searchReq) > 3 {
		searchVal := "%" + searchReq + "%"
		query = query.Where(`description LIKE ? OR "subFeature" LIKE ?`, searchVal, searchVal)
	}

	return query
}
