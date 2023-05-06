package influxdb

import (
	"context"
	"fmt"
	"time"

	"github.com/cloud-disk/internal/config"
	"github.com/cloud-disk/internal/log"

	influx "github.com/influxdata/influxdb-client-go"
)

var InfluxClient influx.Client

type QueryResult struct {
	Time  time.Time   `json:"time"`
	Point string      `json:"point"`
	Value interface{} `json:"value"`
}

func InitInfluxdb() {
	host := config.AppCfg.InfluxCfg.Host
	port := config.AppCfg.InfluxCfg.Port
	userName := config.AppCfg.InfluxCfg.UserName
	password := config.AppCfg.InfluxCfg.Password

	influxURL := fmt.Sprintf("%s:%s%s%s", host, port, userName, password)
	authToken := ""
	InfluxClient = influx.NewClient(influxURL, authToken)
}

// Query 查询influxdb
func Query(queryStr string) ([]*QueryResult, error) {
	queryAPI := InfluxClient.QueryAPI("xuehu96")
	queryRet, err := queryAPI.Query(context.Background(), queryStr)
	if err != nil {
		log.Error("query influxdb error: %s", err)
		return nil, err
	}

	var result []*QueryResult
	for queryRet.Next() {
		queryResult := QueryResult{
			Time:  queryRet.Record().Time(),
			Value: queryRet.Record().Value(),
			Point: queryRet.Record().Values()["pointId"].(string),
		}
		result = append(result, &queryResult)
	}
	return result, nil
}

// Save 保存数据到influxdb
func Save(measurement string, tags map[string]string, fields map[string]interface{}, time time.Time) {
	point := influx.NewPoint(measurement, tags, fields, time)
	writeAPI := InfluxClient.WriteAPI("xuehu96", "test")
	writeAPI.WritePoint(point)
	writeAPI.Flush()

	InfluxClient.Close()
}
