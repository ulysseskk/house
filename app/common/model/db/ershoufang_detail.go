package db

import (
	"context"
	"fmt"
	"strconv"

	"github.com/abyss414/house/app/common/log"

	"github.com/abyss414/house/app/common/connector"
)

const (
	ErshoufangHouseIdHouseInfoRedisKey = "house_id_detail_id"
)

type ErShouFangDetail struct {
	ID                    uint64                  `db:"id" json:"id"`
	HouseID               string                  `db:"house_id" json:"house_id"`
	Platform              string                  `json:"platform"`
	Url                   string                  `json:"url"`
	Title                 string                  `db:"title" json:"title"`
	DirectionText         string                  `json:"direction_text"`
	XiaoquID              uint64                  `db:"xiaoqu_id" json:"xiaoqu_id"`
	XiaoquName            string                  `db:"xiaoqu_name" json:"xiaoqu_name"`
	RegionID              uint64                  `db:"region_id" json:"region_id"`
	Forward               uint64                  `db:"forward" json:"forward"`
	BuildingArea          float64                 `db:"building_area" json:"building_area"`
	HouseType             string                  `db:"house_type" json:"house_type"`
	Fixture               string                  `db:"fixture" json:"fixture"`
	Warm                  string                  `db:"warm" json:"warm"`
	Floor                 string                  `db:"floor" json:"floor"`
	BuildingStructure     string                  `db:"building_structure" json:"building_structure"`
	BuildingType          string                  `db:"building_type" json:"building_type"`
	Elevator              uint64                  `db:"elevator" json:"elevator"`
	OnBoardAt             string                  `db:"on_board_at" json:"on_board_at"`
	LastSale              string                  `db:"last_sale" json:"last_sale"`
	KeepYears             string                  `db:"keep_years" json:"keep_years"`
	Mortgage              string                  `db:"mortgage" json:"mortgage"`
	AreaText              string                  `json:"area_text"`
	InnerArea             float64                 `db:"inner_area" json:"inner_area"`
	Ownership             string                  `db:"ownership" json:"ownership"`
	HouseUsage            string                  `db:"house_usage" json:"house_usage"`
	PropertyRight         string                  `db:"property_right" json:"property_right"`
	ScrapAt               int64                   `db:"scrap_at" json:"scrap_at"`
	HouseCredentialStatus string                  `db:"house_credential_status" json:"house_credential_status"`
	ElevatorRate          string                  `db:"elevator_rate" json:"elevator_rate"`
	HouseStructureType    string                  `db:"house_structure_type" json:"house_structure_type"`
	KeySellingPoints      string                  `db:"key_selling_points" json:"key_selling_points"`
	XiaoquIntro           string                  `db:"xiaoqu_intro" json:"xiaoqu_intro"`
	SellingDetail         string                  `db:"selling_detail" json:"selling_detail"`
	TexAnalyse            string                  `db:"tex_analyse" json:"tex_analyse"`
	FixtureDesc           string                  `db:"fixture_desc" json:"fixture_desc"`
	OwnershipMortgage     string                  `db:"ownership_mortgage" json:"ownership_mortgage"`
	DiagramID             uint64                  `db:"diagram_id" json:"diagram_id"`
	PriceHistory          *ErShouFangPriceHistory `json:"price_history" gorm:"-"`
}

func (detail *ErShouFangDetail) Upsert(ctx context.Context) error {
	if detail.HouseID == "" {
		return nil
	}
	redisClient, err := connector.GetRedisClient()
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("获取redis client from pool错误")
		return err
	}
	defer redisClient.Close()
	redisKey := fmt.Sprintf("%s%s", detail.Platform, ErshoufangHouseIdHouseInfoRedisKey)
	exist, err := redisClient.HExists(redisKey, detail.HouseID)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("scrapper", detail.Platform).Errorf("HExists错误, out key %s.Inner key %s", ZiroomHouseIdRentDetailIDRedisKey, detail.HouseID)
		return err
	}
	if exist {
		id, err := redisClient.HGet(redisKey, detail.HouseID)
		if err != nil {
			log.GlobalLogger().WithError(err).Errorf("HGET失败")
			return err
		}
		detail.ID, err = strconv.ParseUint(id, 10, 64)
		if err != nil {
			log.GlobalLogger().WithError(err).Errorf("Parse id 失败。从redis中获取到的id为%s", id)
			return err
		}
		err = detail.Update(ctx)
		if err != nil {
			return err
		}
	} else {
		err := detail.Insert(ctx)
		if err != nil {
			return err
		}
		_, err = redisClient.HSet(redisKey, detail.HouseID, strconv.FormatUint(detail.ID, 10))
		if err != nil {
			log.GlobalLogger().WithError(err).Errorf("HSET失败")
			return err
		}
	}
	if detail.PriceHistory != nil {
		err := detail.PriceHistory.Insert(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (detail *ErShouFangDetail) Insert(ctx context.Context) error {
	if err := connector.GetMysqlConnector(ctx).Create(detail).Error; err != nil {
		log.GlobalLogger().WithError(err).Errorf("插入Ershoufang Detai失败")
		return err
	}
	return nil
}

func (detail *ErShouFangDetail) Update(ctx context.Context) error {
	if err := connector.GetMysqlConnector(ctx).Save(detail).Error; err != nil {
		log.GlobalLogger().WithError(err).Errorf("更新Ershoufang Detai失败")
		return err
	}
	return nil
}
