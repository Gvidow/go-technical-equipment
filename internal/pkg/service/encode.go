package service

import (
	"errors"
	"net/url"
	"strconv"
	"time"

	"github.com/gvidow/go-technical-equipment/internal/app/ds"
)

func encodeFeedConfig(u *url.URL) (ds.FeedEquipmentConfig, error) {
	cfg := ds.FeedEquipmentConfig{}
	var err error
	if u.Query().Has("createdAfter") {
		date, err := time.Parse("02.01.2006", u.Query().Get("createdAfter"))
		if err != nil {
			return cfg, err
		}
		cfg.SetDateCreateFilter(date)
	}
	if u.Query().Has("title") {
		cfg.SetTitleFilter(u.Query().Get("title"))
	}

	switch u.Query().Get("status") {
	case "all":
		cfg.Status = ds.All
	case "delete":
		cfg.Status = ds.Delete
	case "active", "":
		cfg.Status = ds.Active
	default:
		return ds.FeedEquipmentConfig{}, errors.New("bad status in query params")
	}

	if u.Query().Has("inStock") {
		cfg.InStock, err = strconv.ParseBool(u.Query().Get("inStock"))
		if err != nil {
			return ds.FeedEquipmentConfig{}, err
		}
	}
	return cfg, nil
}
