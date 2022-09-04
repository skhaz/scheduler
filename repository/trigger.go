package repository

import (
	"fmt"
	"time"

	"github.com/skhaz/scheduler/model"
)

const (
	Order = "created_at"
)

type TriggerRepository struct {
	GormRepository
}

func (r *TriggerRepository) List(after time.Time, limit int) (any, error) {
	var c model.TriggerCollection

	err := r.db.Limit(limit).Order(Order).Where(fmt.Sprintf("%s > ?", Order), after).Limit(limit).Find(&c).Error

	return c, err
}

func (r *TriggerRepository) Get(id any) (any, error) {
	var e *model.Trigger

	err := r.db.Where("id = ?", id).First(&e).Error

	return e, err
}

func (r *TriggerRepository) Create(entity any) (any, error) {
	e := entity.(*model.Trigger)

	err := r.db.Create(e).Error

	return e, err
}

func (r *TriggerRepository) Update(id any, entity any) (bool, error) {
	e := entity.(*model.Trigger)

	if err := r.db.Model(e).Where("id = ?", id).Updates(e).Error; err != nil {
		return false, err
	}

	return true, nil
}

func (r *TriggerRepository) Delete(id any) (bool, error) {
	if err := r.db.Delete(&model.Trigger{}, "id = ?", id).Error; err != nil {
		return false, err
	}

	return true, nil
}
