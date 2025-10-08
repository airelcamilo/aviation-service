package repository_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"testing"

	"aviation-service/internal/dto"
	. "aviation-service/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

var (
	atype        = "Airport"
	facilityName = "Lorem ipsum"
	faa          = "LAX"
	region       = "LA"
	state        = "CALIFORNIA"
	county       = "LOS ANGELES"
	city         = "LOS ANGELES"
	ownership    = "PU"
	use          = "PU"
	manager      = "Manager1"
	managerPhone = "123456"
	latitude     = "12.34"
	longitude    = "56.78"
)

func setupMockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open sqlmock database: %s", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return sqlxDB, mock
}

func TestAirportRepository_GetAll(t *testing.T) {
	tests := []struct {
		name        string
		mockRows    *sqlmock.Rows
		mockError   error
		expectedLen int
		expectedErr error
	}{
		{
			name: "Success get all airports",
			mockRows: sqlmock.NewRows([]string{
				"id", "type", "facility_name", "faa", "icao", "region", "state",
				"county", "city", "ownership", "use", "manager", "manager_phone", "latitude", "longitude", "status",
			}).AddRow(1, "Airport", "Lorem ipsum", "LAX", "KLAX", "LA", "CALIFORNIA",
				"LOS ANGELES", "LOS ANGELES", "PU", "PU", "Manager1", "123456", 12.34, 56.78, "PENDING"),
			expectedLen: 1,
		},
		{
			name:        "Error DB",
			mockRows:    nil,
			mockError:   sql.ErrConnDone,
			expectedLen: 0,
			expectedErr: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupMockDB(t)
			defer db.Close()

			repo := NewAirportRepository(db)
			query := `SELECT (.+) FROM airport`

			if tt.mockError != nil {
				mock.ExpectQuery(query).WillReturnError(tt.mockError)
			} else {
				mock.ExpectQuery(query).WillReturnRows(tt.mockRows)
			}

			got, err := repo.GetAll(context.Background(), 20, 0)
			if err != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}

			if len(got) != tt.expectedLen {
				t.Errorf("Expected length %v, got %v", tt.expectedLen, len(got))
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unmet SQL expectations: %v", err)
			}
		})
	}
}

func TestAirportRepository_GetAllByPending(t *testing.T) {
	tests := []struct {
		name        string
		mockRows    *sqlmock.Rows
		mockError   error
		expectedLen int
		expectedErr error
	}{
		{
			name: "Success get all pending airports",
			mockRows: sqlmock.NewRows([]string{
				"id", "type", "facility_name", "faa", "icao", "region", "state",
				"county", "city", "ownership", "use", "manager", "manager_phone", "latitude", "longitude", "status",
			}).AddRow(1, "Airport", "Lorem ipsum", "LAX", "KLAX", "LA", "CALIFORNIA",
				"LOS ANGELES", "LOS ANGELES", "PU", "PU", "Manager1", "123456", 12.34, 56.78, "PENDING"),
			expectedLen: 1,
		},
		{
			name: "No pending airports",
			mockRows: sqlmock.NewRows([]string{
				"id", "type", "facility_name", "faa", "icao", "region", "state",
				"county", "city", "ownership", "use", "manager", "manager_phone", "latitude", "longitude", "status",
			}),
			expectedLen: 0,
		},
		{
			name:        "Error DB",
			mockRows:    nil,
			mockError:   sql.ErrConnDone,
			expectedLen: 0,
			expectedErr: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupMockDB(t)
			defer db.Close()

			repo := NewAirportRepository(db)
			query := `SELECT (.+) FROM airport WHERE status = (.+)`

			if tt.mockError != nil {
				mock.ExpectQuery(query).WillReturnError(tt.mockError)
			} else {
				mock.ExpectQuery(query).WillReturnRows(tt.mockRows)
			}

			got, err := repo.GetAllPending(context.Background())
			if err != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}

			if len(got) != tt.expectedLen {
				t.Errorf("Expected length %v, got %v", tt.expectedLen, len(got))
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unmet SQL expectations: %v", err)
			}
		})
	}
}

func TestAirportRepository_Insert(t *testing.T) {
	tests := []struct {
		name           string
		mockRows       *sqlmock.Rows
		mockError      error
		expectedResult *dto.Airport
		expectedErr    error
	}{
		{
			name: "Success insert airport",
			mockRows: sqlmock.NewRows([]string{
				"type", "facility_name", "faa", "icao", "region", "state",
				"county", "city", "ownership", "use", "manager", "manager_phone", "latitude", "longitude", "status",
			}).AddRow("Airport", "Lorem Ipsum", "LAX", "KLAX", "LA", "CALIFORNIA",
				"LOS ANGELES", "LOS ANGELES", "PU", "PU", "Manager1", "123456", 12.34, 56.78, "PENDING"),
			expectedResult: &dto.Airport{ID: 0, Type: &atype, FacilityName: &facilityName, FAA: &faa, ICAO: "KLAX", Region: &region, State: &state, County: &county, City: &city, Ownership: &ownership,
				Use: &use, Manager: &manager, ManagerPhone: &managerPhone, Latitude: &latitude, Longitude: &longitude, Status: "PENDING"},
		},
		{
			name:           "Error DB",
			mockRows:       nil,
			mockError:      sql.ErrConnDone,
			expectedResult: &dto.Airport{},
			expectedErr:    sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupMockDB(t)
			defer db.Close()

			repo := NewAirportRepository(db)
			query := `INSERT INTO airport (.+) VALUES (.+)`

			if tt.mockError != nil {
				mock.ExpectQuery(query).WillReturnError(tt.mockError)
			} else {
				mock.ExpectQuery(query).WillReturnRows(tt.mockRows)
			}

			got, err := repo.Insert(context.Background(), &dto.Airport{})
			if err != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}

			if got.ICAO != tt.expectedResult.ICAO {
				t.Errorf("Expected result %v, got %v", tt.expectedResult, got)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unmet SQL expectations: %v", err)
			}
		})
	}
}

func TestAirportRepository_GetByICAOOrFacilityName(t *testing.T) {
	tests := []struct {
		name         string
		mockRows     *sqlmock.Rows
		mockError    error
		expectedLen  int
		expectedErr  error
		icao         string
		facilityName string
	}{
		{
			name: "Success get airport by icao",
			mockRows: sqlmock.NewRows([]string{
				"type", "facility_name", "faa", "icao", "region", "state",
				"county", "city", "ownership", "use", "manager", "manager_phone", "latitude", "longitude", "status",
			}).AddRow("Airport", "Lorem ipsum", "LAX", "KLAX", "LA", "CALIFORNIA",
				"LOS ANGELES", "LOS ANGELES", "PU", "PU", "Manager1", "123456", 12.34, 56.78, "PENDING"),
			expectedLen: 1,
			icao:        "KLAX",
		},
		{
			name: "Success get airport by facility name",
			mockRows: sqlmock.NewRows([]string{
				"type", "facility_name", "faa", "icao", "region", "state",
				"county", "city", "ownership", "use", "manager", "manager_phone", "latitude", "longitude", "status",
			}).AddRow("Airport", "Lorem ipsum", "LAX", "KLAX", "LA", "CALIFORNIA",
				"LOS ANGELES", "LOS ANGELES", "PU", "PU", "Manager1", "123456", 12.34, 56.78, "PENDING"),
			expectedLen:  1,
			facilityName: "Lorem ipsum",
		},
		{
			name: "Success get airport by icao and facility name",
			mockRows: sqlmock.NewRows([]string{
				"type", "facility_name", "faa", "icao", "region", "state",
				"county", "city", "ownership", "use", "manager", "manager_phone", "latitude", "longitude", "status",
			}).AddRow("Airport", "Lorem ipsum", "LAX", "KLAX", "LA", "CALIFORNIA",
				"LOS ANGELES", "LOS ANGELES", "PU", "PU", "Manager1", "123456", 12.34, 56.78, "PENDING"),
			expectedLen:  1,
			icao:         "KLAX",
			facilityName: "Lorem ipsum",
		},
		{
			name:        "Error no rows",
			mockRows:    nil,
			mockError:   sql.ErrNoRows,
			expectedLen: 0,
		},
		{
			name:        "Error DB",
			mockRows:    nil,
			mockError:   sql.ErrConnDone,
			expectedLen: 0,
			expectedErr: sql.ErrConnDone,
		},
	}

	limit, offset := 20, 0
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupMockDB(t)
			defer db.Close()

			repo := NewAirportRepository(db)
			query := `SELECT (.+) FROM airport`

			args := []driver.Value{}
			if tt.icao != "" {
				args = append(args, tt.icao)
			}
			if tt.facilityName != "" {
				args = append(args, "%"+tt.facilityName+"%")
			}
			args = append(args, limit, offset)

			if tt.mockError != nil {
				mock.ExpectQuery(query).WithArgs(args...).WillReturnError(tt.mockError)
			} else {
				mock.ExpectQuery(query).WithArgs(args...).WillReturnRows(tt.mockRows)
			}

			got, err := repo.GetByICAOOrFacilityName(context.Background(), tt.icao, tt.facilityName, limit, offset)
			if err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}

			if len(got) != tt.expectedLen {
				t.Errorf("Expected len %v, got %v", tt.expectedLen, len(got))
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unmet SQL expectations: %v", err)
			}
		})
	}
}

func TestAirportRepository_GetById(t *testing.T) {
	tests := []struct {
		name           string
		id             int
		mockRows       *sqlmock.Rows
		mockError      error
		expectedResult *dto.Airport
		expectedErr    error
	}{
		{
			name: "Success get airport by id",
			id:   0,
			mockRows: sqlmock.NewRows([]string{
				"type", "facility_name", "faa", "icao", "region", "state",
				"county", "city", "ownership", "use", "manager", "manager_phone", "latitude", "longitude", "status",
			}).AddRow("Airport", "Lorem ipsum", "LAX", "KLAX", "LA", "CALIFORNIA",
				"LOS ANGELES", "LOS ANGELES", "PU", "PU", "Manager1", "123456", 12.34, 56.78, "PENDING"),
			expectedResult: &dto.Airport{ID: 0, Type: &atype, FacilityName: &facilityName, FAA: &faa, ICAO: "KLAX", Region: &region, State: &state, County: &county, City: &city,
				Ownership: &ownership, Use: &use, Manager: &manager, ManagerPhone: &managerPhone, Latitude: &latitude, Longitude: &longitude, Status: "PENDING"},
		},
		{
			name:           "Error no rows",
			mockRows:       nil,
			mockError:      sql.ErrNoRows,
			expectedResult: nil,
		},
		{
			name:           "Error DB",
			mockRows:       nil,
			mockError:      sql.ErrConnDone,
			expectedResult: &dto.Airport{},
			expectedErr:    sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupMockDB(t)
			defer db.Close()

			repo := NewAirportRepository(db)

			query := `SELECT (.+) WHERE id = (.+)`

			if tt.mockError != nil {
				mock.ExpectQuery(query).WithArgs(tt.id).WillReturnError(tt.mockError)
			} else {
				mock.ExpectQuery(query).WithArgs(tt.id).WillReturnRows(tt.mockRows)
			}

			got, err := repo.GetById(context.Background(), tt.id)
			if err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}

			if !reflect.DeepEqual(got, tt.expectedResult) {
				t.Errorf("Expected result %v, got %v", tt.expectedResult, got)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unmet SQL expectations: %v", err)
			}
		})
	}
}

func TestAirportRepository_UpdateById(t *testing.T) {
	airport := &dto.Airport{ID: 0, Type: &atype, FacilityName: &facilityName, FAA: &faa, ICAO: "KLAX", Region: &region, State: &state, County: &county, City: &city,
		Ownership: &ownership, Use: &use, Manager: &manager, ManagerPhone: &managerPhone, Latitude: &latitude, Longitude: &longitude, Status: "PENDING"}
	tests := []struct {
		name           string
		id             int
		mockRows       *sqlmock.Rows
		mockError      error
		expectedResult *dto.Airport
		expectedErr    error
	}{
		{
			name: "Success update airport by id",
			id:   0,
			mockRows: sqlmock.NewRows([]string{
				"type", "facility_name", "faa", "icao", "region", "state",
				"county", "city", "ownership", "use", "manager", "manager_phone", "latitude", "longitude", "status",
			}).AddRow("Airport", "Lorem ipsum", "LAX", "KLAX", "LA", "CALIFORNIA",
				"LOS ANGELES", "LOS ANGELES", "PU", "PU", "Manager1", "123456", 12.34, 56.78, "PENDING"),
			expectedResult: airport,
		},
		{
			name:           "Error DB",
			mockRows:       nil,
			mockError:      sql.ErrConnDone,
			expectedResult: &dto.Airport{},
			expectedErr:    sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupMockDB(t)
			defer db.Close()

			repo := NewAirportRepository(db)

			query := `UPDATE airport SET (.+) WHERE id = (.+)`

			if tt.mockError != nil {
				mock.ExpectQuery(query).WillReturnError(tt.mockError)
			} else {
				mock.ExpectQuery(query).WillReturnRows(tt.mockRows)
			}

			got, err := repo.UpdateById(context.Background(), airport)
			if err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}

			if !reflect.DeepEqual(got, tt.expectedResult) {
				t.Errorf("Expected result %v, got %v", tt.expectedResult, got)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unmet SQL expectations: %v", err)
			}
		})
	}
}

func TestAirportRepository_UpdateByICAO(t *testing.T) {
	airports := []dto.Airport{{ID: 0, Type: &atype, FacilityName: &facilityName, FAA: &faa, ICAO: "KLAX", Region: &region, State: &state, County: &county, City: &city,
		Ownership: &ownership, Use: &use, Manager: &manager, ManagerPhone: &managerPhone, Latitude: &latitude, Longitude: &longitude, Status: "PENDING"},
		{ID: 0, Type: &atype, FacilityName: &facilityName, FAA: &faa, ICAO: "KLAX", Region: &region, State: &state, County: &county, City: &city,
			Ownership: &ownership, Use: &use, Manager: &manager, ManagerPhone: &managerPhone, Latitude: &latitude, Longitude: &longitude, Status: "PENDING"}}
	tests := []struct {
		name        string
		mockError   error
		expectedErr error
	}{
		{
			name:        "Success update airport by icao",
			mockError:   nil,
			expectedErr: nil,
		},
		{
			name:        "Error DB",
			mockError:   sql.ErrConnDone,
			expectedErr: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupMockDB(t)
			defer db.Close()

			repo := NewAirportRepository(db)

			query := `UPDATE airport AS a SET`

			if tt.mockError != nil {
				mock.ExpectExec(query).WillReturnError(tt.mockError)
			} else {
				mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, int64(len(airports))))
			}

			err := repo.UpdateByICAO(context.Background(), airports)
			if err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unmet SQL expectations: %v", err)
			}
		})
	}
}

func TestAirportRepository_DELETE(t *testing.T) {
	tests := []struct {
		name        string
		id          int
		mockError   error
		mockResult  driver.Result
		expectedErr error
	}{
		{
			name:       "Success delete airport by id",
			id:         0,
			mockResult: sqlmock.NewResult(1, 1),
		},
		{
			name:        "No airport found to be deleted",
			mockResult:  sqlmock.NewResult(0, 0),
			expectedErr: fmt.Errorf("No airport found with id 0"),
		},
		{
			name:        "Error DB",
			mockError:   sql.ErrConnDone,
			expectedErr: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupMockDB(t)
			defer db.Close()

			repo := NewAirportRepository(db)

			query := `DELETE FROM airport WHERE id = (.+)`

			if tt.mockError != nil {
				mock.ExpectExec(query).WithArgs(tt.id).WillReturnError(tt.mockError)
			} else {
				mock.ExpectExec(query).WithArgs(tt.id).WillReturnResult(tt.mockResult)
			}

			err := repo.Delete(context.Background(), tt.id)
			if err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unmet SQL expectations: %v", err)
			}
		})
	}
}
