package sql

import (
	"log"
	"time"

	"github.com/authorizerdev/authorizer/server/db/models"
	"github.com/authorizerdev/authorizer/server/graph/model"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

// AddVerification to save verification request in database
func (p *provider) AddVerificationRequest(verificationRequest models.VerificationRequest) (models.VerificationRequest, error) {
	if verificationRequest.ID == "" {
		verificationRequest.ID = uuid.New().String()
	}

	verificationRequest.Key = verificationRequest.ID
	verificationRequest.CreatedAt = time.Now().Unix()
	verificationRequest.UpdatedAt = time.Now().Unix()
	result := p.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}, {Name: "identifier"}},
		DoUpdates: clause.AssignmentColumns([]string{"token", "expires_at", "nonce", "redirect_uri"}),
	}).Create(&verificationRequest)

	if result.Error != nil {
		log.Println(`error saving verification request record`, result.Error)
		return verificationRequest, result.Error
	}

	return verificationRequest, nil
}

// GetVerificationRequestByToken to get verification request from database using token
func (p *provider) GetVerificationRequestByToken(token string) (models.VerificationRequest, error) {
	var verificationRequest models.VerificationRequest
	result := p.db.Where("token = ?", token).First(&verificationRequest)

	if result.Error != nil {
		log.Println(`error getting verification request:`, result.Error)
		return verificationRequest, result.Error
	}

	return verificationRequest, nil
}

// GetVerificationRequestByEmail to get verification request by email from database
func (p *provider) GetVerificationRequestByEmail(email string, identifier string) (models.VerificationRequest, error) {
	var verificationRequest models.VerificationRequest

	result := p.db.Where("email = ? AND identifier = ?", email, identifier).First(&verificationRequest)

	if result.Error != nil {
		log.Println(`error getting verification token:`, result.Error)
		return verificationRequest, result.Error
	}

	return verificationRequest, nil
}

// ListVerificationRequests to get list of verification requests from database
func (p *provider) ListVerificationRequests(pagination model.Pagination) (*model.VerificationRequests, error) {
	var verificationRequests []models.VerificationRequest

	result := p.db.Limit(int(pagination.Limit)).Offset(int(pagination.Offset)).Order("created_at DESC").Find(&verificationRequests)
	if result.Error != nil {
		log.Println("error getting verification requests:", result.Error)
		return nil, result.Error
	}

	responseVerificationRequests := []*model.VerificationRequest{}
	for _, v := range verificationRequests {
		responseVerificationRequests = append(responseVerificationRequests, v.AsAPIVerificationRequest())
	}

	var total int64
	totalRes := p.db.Model(&models.VerificationRequest{}).Count(&total)
	if totalRes.Error != nil {
		return nil, totalRes.Error
	}

	paginationClone := pagination
	paginationClone.Total = total

	return &model.VerificationRequests{
		VerificationRequests: responseVerificationRequests,
		Pagination:           &paginationClone,
	}, nil
}

// DeleteVerificationRequest to delete verification request from database
func (p *provider) DeleteVerificationRequest(verificationRequest models.VerificationRequest) error {
	result := p.db.Delete(&verificationRequest)

	if result.Error != nil {
		log.Println(`error deleting verification request:`, result.Error)
		return result.Error
	}

	return nil
}
