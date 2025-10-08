package dto

type Airport struct {
	ID           int     `db:"id" json:"id"`
	Type         *string `db:"type" json:"type,omitempty"`
	FacilityName *string `db:"facility_name" json:"facility_name,omitempty"`
	FAA          *string `db:"faa" json:"faa_ident,omitempty"`
	ICAO         string  `db:"icao" json:"icao_ident"`
	Region       *string `db:"region" json:"region,omitempty"`
	State        *string `db:"state" json:"state_full,omitempty"`
	County       *string `db:"county" json:"county,omitempty"`
	City         *string `db:"city" json:"city,omitempty"`
	Ownership    *string `db:"ownership" json:"ownership,omitempty"`
	Use          *string `db:"use" json:"use,omitempty"`
	Manager      *string `db:"manager" json:"manager,omitempty"`
	ManagerPhone *string `db:"manager_phone" json:"manager_phone,omitempty"`
	Latitude     *string `db:"latitude" json:"latitude,omitempty"`
	Longitude    *string `db:"longitude" json:"longitude,omitempty"`
	Status       string  `db:"status" json:"status"`
}

type AirportDataResponse map[string][]Airport
