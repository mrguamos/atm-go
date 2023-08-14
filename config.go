package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

type configService struct {
	db *sqlx.DB
}

type Config struct {
	Key   string `db:"key" json:"key"`
	Value string `db:"value" json:"value,omitempty"`
}

func (s *configService) getConfigs(keys ...string) ([]Config, error) {
	config := []Config{}
	if len(keys) > 0 {
		query, args, _ := sqlx.In("SELECT * FROM config WHERE key IN (?);", keys)
		query = s.db.Rebind(query)
		err := s.db.Select(&config, query, args...)
		return config, err
	}
	return s.loadConfigs()
}

func (s *configService) loadConfigs() ([]Config, error) {
	config := []Config{}
	err := s.db.Select(&config, "SELECT * FROM config")
	return config, err
}

func (s *configService) updateConfigs(configs []Config) error {
	tx, _ := s.db.Begin()
	oldConfig := make(map[string]string)
	for _, c := range configs {
		oldConfig[c.Key] = c.Value
		_, err := tx.Exec("UPDATE config SET value = $1 WHERE key = $2", c.Value, c.Key)
		if err != nil {
			tx.Rollback()
			for k, v := range oldConfig {
				viper.SetDefault(k, v)
			}
			return err
		}
		viper.SetDefault(c.Key, c.Value)
	}
	tx.Commit()
	return nil
}
