package db

import (
	"errors"
	"gorm.io/gorm"
	"tableflow/go/pkg/model"
	"tableflow/go/pkg/tf"
)

func GetOrganizationOfUserWithWorkspaces(userID string) (*model.Organization, error) {
	if len(userID) == 0 {
		return nil, errors.New("no user ID provided")
	}
	var organization model.Organization
	err := tf.DB.Preload("Workspaces").
		Where("id in (select organization_id from organization_users where user_id = ?)", model.ParseID(userID)).
		First(&organization).Error
	if err != nil {
		return nil, err
	}
	if !organization.ID.Valid {
		return nil, gorm.ErrRecordNotFound
	}
	return &organization, nil
}

// GetOrganizationFromWorkspace TODO: Consider adding organization ID to other objects so it doesn't have to be retrieved
func GetOrganizationFromWorkspace(workspaceID string) (*model.Organization, error) {
	if len(workspaceID) == 0 {
		return nil, errors.New("no workspace ID provided")
	}
	var organization model.Organization
	err := tf.DB.Joins("join workspaces on organizations.id = workspaces.organization_id").
		Where("workspaces.id = ?", model.ParseID(workspaceID)).
		First(&organization).Error
	if err != nil {
		return nil, err
	}
	if !organization.ID.Valid {
		return nil, gorm.ErrRecordNotFound
	}
	return &organization, nil
}
