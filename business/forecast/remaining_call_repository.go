package forecast

import (
	"database/sql"
	"kovaja/sun-forecast/utils/logger"
)

type RemainingCallRepository struct {
	db *sql.DB
}

const MAX_CALLS = 8

func (repository RemainingCallRepository) ReadRemainingCalls(todayKey string) (int, error) {
	var remaining int
	err := repository.db.QueryRow("SELECT remaining FROM remaining_calls WHERE date = $1", todayKey).Scan(&remaining)

	if err != nil {
		if err == sql.ErrNoRows {
			// this means that we don't have that day yet, so return MAX_CALLS
			// we will write it later on
			logger.Log("No record found for %s, returning MAX_CALLS", todayKey)
			remaining = MAX_CALLS
		} else {
			return 0, err
		}
	}

	return remaining, nil
}

func (repository RemainingCallRepository) UpdateRemainingCalls(todayKey string, remainingCalls int) error {
	query := `
    INSERT INTO remaining_calls (date, remaining)
    VALUES ($1, $2)
    ON CONFLICT (date) DO UPDATE SET remaining = EXCLUDED.remaining
  `
	_, err := repository.db.Exec(query, todayKey, remainingCalls)

	logger.Log("Setting remaining calls for %s as %d", todayKey, remainingCalls)
	return err
}

func InitializeRemainigCallRepository(db *sql.DB) RemainingCallRepository {
	return RemainingCallRepository{
		db: db,
	}
}
