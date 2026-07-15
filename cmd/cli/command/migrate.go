package command

import (
	"fmt"
	"go-api/domain/entity"
	"go-api/infrastructure/config"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

func NewMigrateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Migrate the database",
		Long:  "Migrate the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			env := config.Load()
			db := config.ConnectDatabase(env)

			err := db.Transaction(func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(
					&entity.User{},
					&entity.Insight{},
					&entity.Media{},
					&entity.Signal{},
				); err != nil {
					return err
				}

				// GORM AutoMigrate does not drop NOT NULL on existing columns.
				return tx.Exec(`ALTER TABLE medias ALTER COLUMN insight_id DROP NOT NULL`).Error
			})

			if err != nil {
				return fmt.Errorf(
					"🚨 failed to migrate database: %w",
					err,
				)
			}

			fmt.Println("🎉 Database migrations completed successfully")
			return nil
		},
	}
}
