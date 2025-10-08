package seeding

import (
    "fmt"
    "os"
    "path/filepath"
    "sort"

    "github.com/jmoiron/sqlx"
	"aviation-service/pkg/logger"
)

func RunSeeding(db *sqlx.DB, seedingsDir string) error {
    files, err := filepath.Glob(filepath.Join(seedingsDir, "*.sql"))
    if err != nil {
        return fmt.Errorf("Failed to read seeding: %s", err)
    }
    if len(files) == 0 {
        return fmt.Errorf("No seeding file: %s", err)
    }
    sort.Strings(files)

    for _, file := range files {
        filename := filepath.Base(file)

        sqlBytes, err := os.ReadFile(file)
        if err != nil {
            return fmt.Errorf("Failed to read file %s: %s", filename, err)
        }
        sql := string(sqlBytes)

        logger.Infow("Applying seeding", "fileName", filename)
        if _, err := db.Exec(sql); err != nil {
            return fmt.Errorf("Failed to execute seeding %s: %s", filename, err)
        }
    }

    logger.Info("All seeding applied successfully")
    return nil
}
