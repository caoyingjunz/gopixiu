/*
Copyright 2021 The Pixiu Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cloud

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/caoyingjunz/gopixiu/pkg/db/model"
)

type KubeConfigInterface interface {
	Create(ctx context.Context, obj *model.KubeConfig) (*model.KubeConfig, error)
	Update(ctx context.Context, id, resourceVersion int64, updates map[string]interface{}) error
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (*model.KubeConfig, error)
	List(ctx context.Context, cloudName string) ([]model.KubeConfig, error)
	ListByIds(ctx context.Context, ids []int64) ([]model.KubeConfig, error)
}

type kubeConfig struct {
	db *gorm.DB
}

func NewKubeConfig(db *gorm.DB) KubeConfigInterface {
	return &kubeConfig{db}
}

func (s *kubeConfig) Create(ctx context.Context, obj *model.KubeConfig) (*model.KubeConfig, error) {
	now := time.Now()
	obj.GmtCreate = now
	obj.GmtModified = now

	if err := s.db.Create(obj).Error; err != nil {
		return nil, err
	}
	return obj, nil
}

func (s *kubeConfig) Update(ctx context.Context, id, resourceVersion int64, updates map[string]interface{}) error {
	updates["gmt_modified"] = time.Now()
	updates["resource_version"] = resourceVersion + 1

	f := s.db.Model(&model.KubeConfig{}).
		Where("id = ? and resource_version = ?", id, resourceVersion).
		Updates(updates)
	if f.Error != nil {
		return f.Error
	}

	return nil
}

func (s *kubeConfig) Delete(ctx context.Context, id int64) error {
	return s.db.
		Delete(&model.KubeConfig{}, id).
		Error
}

func (s *kubeConfig) Get(ctx context.Context, id int64) (*model.KubeConfig, error) {
	var obj model.KubeConfig
	if err := s.db.First(&obj, id).Error; err != nil {
		return nil, err
	}

	return &obj, nil
}

func (s *kubeConfig) List(ctx context.Context, cloudName string) ([]model.KubeConfig, error) {
	var objs []model.KubeConfig
	if err := s.db.Where("cloud_name = ?", cloudName).Find(&objs).Error; err != nil {
		return nil, err
	}

	return objs, nil
}

func (s *kubeConfig) ListByIds(ctx context.Context, ids []int64) ([]model.KubeConfig, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var objs []model.KubeConfig
	if err := s.db.Where("id IN ?", ids).Find(&objs).Error; err != nil {
		return nil, err
	}
	return objs, nil
}
