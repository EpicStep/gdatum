// Copyright 2025 Stepan Rabotkin.
// SPDX-License-Identifier: Apache-2.0.

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/EpicStep/gdatum/internal/domain"
)

const (
	namespaceName                     = "gdatum"
	serverStatsCollectorSubsystemName = "servers_stats_collector"
)

type CollectorMetrics struct {
	serversCollected      *prometheus.GaugeVec
	collectionErrorsTotal *prometheus.CounterVec
	insertErrorsTotal     prometheus.Counter
}

func NewCollectorMetrics(registerer prometheus.Registerer) *CollectorMetrics {
	if registerer == nil {
		registerer = prometheus.DefaultRegisterer
	}

	factory := promauto.With(registerer)
	return &CollectorMetrics{
		serversCollected: factory.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespaceName,
				Subsystem: serverStatsCollectorSubsystemName,
				Name:      "servers_collected",
				Help:      "Number of servers collected from each multiplayer platform",
			},
			[]string{"multiplayer"}),
		collectionErrorsTotal: factory.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespaceName,
				Subsystem: serverStatsCollectorSubsystemName,
				Name:      "collection_errors_total",
				Help:      "Total number of server collection errors by multiplayer",
			},
			[]string{"multiplayer"}),
		insertErrorsTotal: factory.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespaceName,
				Subsystem: serverStatsCollectorSubsystemName,
				Name:      "insert_errors_total",
				Help:      "Total number of errors when inserting server data to repository",
			},
		),
	}
}

func (m *CollectorMetrics) RecordServersCollected(multiplayer domain.Multiplayer, count int) {
	m.serversCollected.
		WithLabelValues(string(multiplayer)).
		Set(float64(count))
}

func (m *CollectorMetrics) RecordCollectionError(multiplayer domain.Multiplayer) {
	m.collectionErrorsTotal.WithLabelValues(string(multiplayer)).Inc()
}

func (m *CollectorMetrics) RecordInsertError() {
	m.insertErrorsTotal.Inc()
}
