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

func encodeFeedRequestConfig(u *url.URL) (ds.FeedRequestConfig, error) {
	cfg := ds.FeedRequestConfig{}
	var err error
	if u.Query().Has("creator") {
		err = cfg.SetCreatorFilter(u.Query().Get("creator"))
		if err != nil {
			return cfg, err
		}
	}

	if u.Query().Has("moderator") {
		err = cfg.SetModeratorFilter(u.Query().Get("moderator"))
		if err != nil {
			return cfg, err
		}
	}

	if u.Query().Has("creatorProfile") {
		cfg.SetCreatorProfileFilter(u.Query().Get("creatorProfile"))
	}

	if u.Query().Has("moderatorProfile") {
		cfg.SetModeratorProfileFilter(u.Query().Get("moderatorProfile"))
	}

	if u.Query().Has("status") {
		cfg.SetStatusFilter(u.Query().Get("status"))
	}

	if u.Query().Has("createdAt") {
		err = cfg.SetCreatedFilter(u.Query().Get("createdAt"))
		if err != nil {
			return cfg, err
		}
	}
	if u.Query().Has("formatedAt") {
		err = cfg.SetFormatedFilter(u.Query().Get("formatedAt"))
		if err != nil {
			return cfg, err
		}
	}
	if u.Query().Has("formatedAfter") {
		err = cfg.SetFormatedAfter(u.Query().Get("formatedAfter"))
		if err != nil {
			return cfg, err
		}
	}
	if u.Query().Has("formatedBefore") {
		err = cfg.SetFormatedBefore(u.Query().Get("formatedBefore"))
		if err != nil {
			return cfg, err
		}
	}
	if u.Query().Has("completedAt") {
		err = cfg.SetCompletedFilter(u.Query().Get("completedAt"))
		if err != nil {
			return cfg, err
		}
	}
	return cfg, nil
}
