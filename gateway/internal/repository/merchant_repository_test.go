package repository

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	uuid "github.com/satori/go.uuid"
)

func TestGetIDByApiKey(t *testing.T) {
	tests := []struct {
		name   string
		input  func() (uuid.UUID, error)
		expect error
	}{
		{
			name: "GetIDByApiKey_ErrorNoRows",
			input: func() (uuid.UUID, error) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}
				defer db.Close()

				merchantId := uuid.NewV4()

				mock.ExpectQuery(regexp.QuoteMeta(getByApiKey)).WithArgs(merchantId).WillReturnError(sql.ErrNoRows)

				svc := NewMerchantRepository(db)

				return svc.GetIDByApiKey(merchantId)
			},
			expect: ErrMerchantNotFound,
		},
		{
			name: "GetIDByApiKey_ClosedConnection",
			input: func() (uuid.UUID, error) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}
				defer db.Close()

				merchantId := uuid.NewV4()

				mock.ExpectQuery(regexp.QuoteMeta(getByApiKey)).WithArgs(merchantId).WillReturnError(sql.ErrConnDone)

				svc := NewMerchantRepository(db)
				return svc.GetIDByApiKey(merchantId)
			},
			expect: sql.ErrConnDone,
		},
		{
			name: "GetIDByApiKey_Success",
			input: func() (uuid.UUID, error) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}
				defer db.Close()

				merchantId := uuid.NewV4()

				mock.ExpectQuery(
					regexp.QuoteMeta(getByApiKey)).
					WithArgs(merchantId).
					WillReturnRows(sqlmock.NewRows(
						[]string{"id"},
					).AddRow(merchantId))

				svc := NewMerchantRepository(db)
				return svc.GetIDByApiKey(merchantId)
			},
			expect: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.input()
			if err != tt.expect {
				t.Errorf("expect values are equals, but got diferent, output: %v, expect: %v", err, tt.expect)
			}
		})
	}
}
