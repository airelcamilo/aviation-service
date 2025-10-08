package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"aviation-service/internal/dto"

	"github.com/jmoiron/sqlx"
)

//go:generate moq -out ../mock/airport_repository_mock.go -pkg=mock . IAirportRepository
type IAirportRepository interface {
	GetAll(ctx context.Context, limit, offset int) ([]dto.Airport, error)
	GetAllPending(ctx context.Context) ([]dto.Airport, error)
	GetById(ctx context.Context, id int) (*dto.Airport, error)
	GetByICAOOrFacilityName(ctx context.Context, icao, facilityName string, limit, offset int) ([]dto.Airport, error)
	Insert(ctx context.Context, airport *dto.Airport) (*dto.Airport, error)
	UpdateById(ctx context.Context, airport *dto.Airport) (*dto.Airport, error)
	UpdateByICAO(ctx context.Context, airports []dto.Airport) error
	Delete(ctx context.Context, id int) error
}

type AirportRepository struct {
	db *sqlx.DB
}

func NewAirportRepository(db *sqlx.DB) *AirportRepository {
	return &AirportRepository{db: db}
}

func (r *AirportRepository) GetAll(ctx context.Context, limit, offset int) ([]dto.Airport, error) {
	var airports []dto.Airport
	query := `SELECT id, type, facility_name, faa, icao, region, state, county, city, ownership, use, 
			  manager, manager_phone, latitude, longitude, status 
			  FROM airport
			  LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &airports, query, limit, offset)
	return airports, err
}

func (r *AirportRepository) GetAllPending(ctx context.Context) ([]dto.Airport, error) {
	var airports []dto.Airport
	query := `SELECT id, type, facility_name, faa, icao, region, state, county, city, ownership, use, 
			  manager, manager_phone, latitude, longitude, status 
			  FROM airport
			  WHERE status = $1`

	err := r.db.SelectContext(ctx, &airports, query, "PENDING")
	return airports, err
}

func (r *AirportRepository) Insert(ctx context.Context, airport *dto.Airport) (*dto.Airport, error) {
	query := `INSERT INTO airport (
				type, facility_name, faa, icao, region, state, county, city, ownership, use, 
				manager, manager_phone, latitude, longitude, status 
			  ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
				RETURNING id, type, facility_name, faa, icao, region, state, county, city, ownership, use, 
			  	manager, manager_phone, latitude, longitude, status`
	var created dto.Airport
	err := r.db.GetContext(ctx, &created, query, airport.Type, airport.FacilityName, airport.FAA,
		airport.ICAO, airport.Region, airport.State, airport.County, airport.City, airport.Ownership,
		airport.Use, airport.Manager, airport.ManagerPhone, airport.Latitude, airport.Longitude, airport.Status)
	return &created, err
}

func (r *AirportRepository) GetByICAOOrFacilityName(ctx context.Context, icao string, facilityName string, limit, offset int) ([]dto.Airport, error) {
	var airports []dto.Airport
	query := `SELECT id, type, facility_name, faa, icao, region, state, county, city, ownership, use, 
			  manager, manager_phone, latitude, longitude, status 
			  FROM airport`
	var conditions []string
	args := []interface{}{}

	if icao != "" {
		args = append(args, icao)
		conditions = append(conditions, fmt.Sprintf("icao = $%d", len(args)))
	}
	if facilityName != "" {
		args = append(args, "%"+facilityName+"%")
		conditions = append(conditions, fmt.Sprintf("facility_name ILIKE $%d", len(args)))
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	args = append(args, limit, offset)
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)-1, len(args))
	err := r.db.SelectContext(ctx, &airports, query, args...)
	return airports, err
}

func (r *AirportRepository) GetById(ctx context.Context, id int) (*dto.Airport, error) {
	var airport dto.Airport
	query := `SELECT id, type, facility_name, faa, icao, region, state, county, city, ownership, use, 
			  manager, manager_phone, latitude, longitude, status 
			  FROM airport 
			  WHERE id = $1`
	err := r.db.GetContext(ctx, &airport, query, id)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &airport, err
}

func (r *AirportRepository) UpdateById(ctx context.Context, airport *dto.Airport) (*dto.Airport, error) {
	query := `UPDATE airport SET
			type = $1,
			facility_name = $2,
			faa = $3,
			icao = $4,
			region = $5,
			state = $6,
			county = $7,
			city = $8,
			ownership = $9,
			use = $10,
			manager = $11,
			manager_phone = $12,
			latitude = $13,
			longitude = $14,
			status = $15
			WHERE id = $16
			RETURNING id, type, facility_name, faa, icao, region, state, county, city, ownership, use, 
			manager, manager_phone, latitude, longitude, status`
	var updated dto.Airport
	err := r.db.GetContext(ctx, &updated, query, airport.Type, airport.FacilityName, airport.FAA,
		airport.ICAO, airport.Region, airport.State, airport.County, airport.City, airport.Ownership,
		airport.Use, airport.Manager, airport.ManagerPhone, airport.Latitude, airport.Longitude, airport.Status, airport.ID)
	return &updated, err
}

func (r *AirportRepository) UpdateByICAO(ctx context.Context, airports []dto.Airport) error {
	values := []interface{}{}
    placeholders := []string{}

    for i, apt := range airports {
        base := i*15 + 1
        placeholders = append(placeholders, fmt.Sprintf(
            "($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
            base, base+1, base+2, base+3, base+4, base+5, base+6,
            base+7, base+8, base+9, base+10, base+11, base+12, base+13, base+14,
        ))

        values = append(values,
            apt.ICAO, 
            apt.Type,
            apt.FacilityName,
            apt.FAA,
            apt.Region,
            apt.State,
            apt.County,
            apt.City,
            apt.Ownership,
            apt.Use,
            apt.Manager,
            apt.ManagerPhone,
            apt.Latitude,
            apt.Longitude,
            apt.Status,
        )
    }

    query := `
        UPDATE airport AS a SET
            type = v.type,
            facility_name = v.facility_name,
            faa = v.faa,
            region = v.region,
            state = v.state,
            county = v.county,
            city = v.city,
            ownership = v.ownership,
            use = v.use,
            manager = v.manager,
            manager_phone = v.manager_phone,
            latitude = v.latitude,
            longitude = v.longitude,
            status = v.status
        FROM (VALUES
    ` + strings.Join(placeholders, ",") + `
        ) AS v(icao, type, facility_name, faa, region, state, county, city, ownership, use,
                manager, manager_phone, latitude, longitude, status)
        WHERE a.icao = v.icao
    `

    _, err := r.db.ExecContext(ctx, query, values...)
    return err
}

func (r *AirportRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM airport WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No airport found with id %d", id)
	}
	return err
}
