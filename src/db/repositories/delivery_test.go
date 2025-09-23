package repositories

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"

	"orders/src/db/models"
	"orders/src/metrics"
)

func newTestRepo(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, *DeliveryRepository) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	hist := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "test_db_query_duration_seconds",
		Help: "help",
	}, []string{"query", "service"})
	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "test_db_query_errors_total",
		Help: "help",
	}, []string{"query", "service"})

	m := &metrics.Metrics{
		DBQueryDuration: hist,
		DBQueryErrors:   counter,
	}

	repo := NewDeliveryRepo(sqlxDB, m)

	return sqlxDB, mock, &repo
}

func TestCreateDelivery_Success(t *testing.T) {
	sqlxDB, mock, repoPtr := newTestRepo(t)
	defer sqlxDB.Close()

	in := &models.Delivery{
		Name:    "Ivan",
		Phone:   "+70000000000",
		Zip:     "12345",
		City:    "Moscow",
		Address: "Lenina 1",
		Region:  "Moscow",
		Email:   "ivan@example.com",
	}

	columns := []string{"id", "name", "phone", "zip", "city", "address", "region", "email"}
	rows := sqlmock.NewRows(columns).
		AddRow(42, in.Name, in.Phone, in.Zip, in.City, in.Address, in.Region, in.Email)

	mock.ExpectQuery(`(?s)^INSERT INTO delivery.*RETURNING id, name, phone, zip, city, address, region, email;?$`).WillReturnRows(rows)

	ctx := context.Background()
	delivery, err := (*repoPtr).CreateDelivery(ctx, in)
	require.NoError(t, err)
	require.Equal(t, 42, delivery.ID)
	require.Equal(t, in.Name, delivery.Name)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestGetDeliveryByID_Success(t *testing.T) {
	sqlxDB, mock, repoPtr := newTestRepo(t)
	defer sqlxDB.Close()

	columns := []string{"id", "name", "phone", "zip", "city", "address", "region", "email", "order_id"}
	rows := sqlmock.NewRows(columns).
		AddRow(7, "Petr", "+71111111111", "54321", "SPb", "Nevsky 1", "SPb region", "petr@example.com", 1001)

	mock.ExpectQuery(`(?s)^select \*\s+from delivery\s+where id = \$1;?$`).WithArgs(7).WillReturnRows(rows)

	ctx := context.Background()
	delivery, err := (*repoPtr).GetDeliveryByID(ctx, 7)
	require.NoError(t, err)
	require.Equal(t, 7, delivery.ID)
	require.Equal(t, "Petr", delivery.Name)
	require.Equal(t, 1001, delivery.OrderID)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
func TestGetDeliveryByID_NotFound(t *testing.T) {
	sqlxDB, mock, repoPtr := newTestRepo(t)
	defer sqlxDB.Close()

	columns := []string{"id", "name", "phone", "zip", "city", "address", "region", "email", "order_id"}
	emptyRows := sqlmock.NewRows(columns)

	mock.ExpectQuery(`(?s)^select \*\s+from delivery\s+where id = \$1;?$`).WithArgs(999).WillReturnRows(emptyRows)

	ctx := context.Background()
	_, err := (*repoPtr).GetDeliveryByID(ctx, 999)

	require.Error(t, err)
	require.Equal(t, sql.ErrNoRows, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
