package utils

import (
	"aviation-service/internal/dto"
	"errors"
	"regexp"
	"strings"
)

type AirportValidator struct {}

func NewAirportValidator() *AirportValidator {
	return &AirportValidator{}
}

func (v *AirportValidator) IsComplete(req *dto.Airport) bool {
	return req.Type != nil && *req.Type != "" &&
		req.FacilityName != nil && *req.FacilityName != "" &&
		req.FAA != nil && *req.FAA != "" &&
		req.ICAO != "" &&
		req.Region != nil && *req.Region != "" &&
		req.State != nil && *req.State != "" &&
		req.County != nil && *req.County != "" &&
		req.City != nil && *req.City != "" &&
		req.Ownership != nil && *req.Ownership != "" &&
		req.Use != nil && *req.Use != "" &&
		req.Manager != nil && *req.Manager != "" &&
		req.ManagerPhone != nil && *req.ManagerPhone != "" &&
		req.Latitude != nil && *req.Latitude != "" &&
		req.Longitude != nil && *req.Longitude != ""
}

func (v *AirportValidator) Validate(req *dto.Airport) error {
	if strings.TrimSpace(req.ICAO) == "" {
		return errors.New("ICAO is required")
	}

	// Ownership & Use must be PU or PR
	if req.Ownership != nil && *req.Ownership != "" && *req.Ownership != "PU" && *req.Ownership != "PR" {
		return errors.New("Ownership must be PU or PR")
	}
	if  req.Use != nil && *req.Use != "" && *req.Use != "PU" && *req.Use != "PR" {
		return errors.New("Use must be PU or PR")
	}

	// Latitude regex
	if req.Latitude != nil && *req.Latitude != "" {
		latPattern := regexp.MustCompile(`^\d{2}-\d{2}-\d{2}(\.\d+)?[NS]$`)
		if !latPattern.MatchString(*req.Latitude) {
			return errors.New("Invalid latitude format (expected DD-MM-SS.sssN/S)")
		}
	}

	// Longitude regex
	if req.Longitude != nil && *req.Longitude != "" {
		longPattern := regexp.MustCompile(`^\d{3}-\d{2}-\d{2}(\.\d+)?[EW]$`)
		if !longPattern.MatchString(*req.Longitude) {
			return errors.New("Invalid longitude format (expected DDD-MM-SS.sssE/W)")
		}
	}

	// Manager phone regex
	if req.ManagerPhone != nil && *req.ManagerPhone != "" {
		phonePattern := regexp.MustCompile(`^\+?[0-9\-\(\)\s]{7,20}$`)
		if !phonePattern.MatchString(*req.ManagerPhone) {
			return errors.New("Invalid manager phone format")
		}
	}

	return nil
}