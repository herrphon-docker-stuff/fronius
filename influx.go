// Copyright 2015 Tamás Gulácsi
//
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

// Package main of fronius gets the data from Solar.Web
package main

import (
	"context"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type influxClient struct {
	influxApi       *api.WriteAPIBlocking
	RetentionPolicy string

	Logger log.Logger
}

func newInfluxClient(influxUrl, influxOrg, influxBucket, retentionPolicy string, logger log.Logger) (influxClient, error) {
	authToken := os.Getenv("INFLUX_AUTH_TOKEN")
	// Create a new client using an InfluxDB server base URL and an authentication token
	client := influxdb2.NewClient(influxUrl, authToken)
	// Use blocking write client for writes to desired bucket
	writeAPI := client.WriteAPIBlocking(influxOrg, influxBucket)

	return influxClient{influxApi: &writeAPI, RetentionPolicy: retentionPolicy, Logger: logger}, nil
}

type dataPoint struct {
	time.Time
	Name  string
	Value float64
	Unit  string
}

func (c influxClient) Put(measurement string, points ...dataPoint) error {
	ip := make([]*write.Point, len(points))

	for i, p := range points {
		ip[i] = influxdb2.NewPoint(
			measurement,
			map[string]string{"name": p.Name},
			map[string]interface{}{"energy": p.Value, "unit": p.Unit},
			p.Time)
	}

	influxApi := *c.influxApi
	err := influxApi.WritePoint(context.Background(), ip...)

	return err
}
