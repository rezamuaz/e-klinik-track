package dto

import "e-klinik/infra/pg"

func ToApexChartSkp7Hari(rows []pg.GetCapaianSKP7HariTerakhirRow) ChartResponse {
	var categories []string
	var diverifikasi, ditolak, belum []float64

	for _, r := range rows {
		categories = append(categories, r.Tanggal.Time.Format("2006-01-02"))
		diverifikasi = append(diverifikasi, float64(r.Diverifikasi))
		ditolak = append(ditolak, float64(r.Ditolak))
		belum = append(belum, float64(r.BelumDiverifikasi))
	}

	series := []ChartSeries{
		{Name: "Diverifikasi", Data: diverifikasi},
		{Name: "Ditolak", Data: ditolak},
		{Name: "Belum Diverifikasi", Data: belum},
	}

	return ChartResponse{
		Series:     series,
		Categories: categories,
	}
}

type ChartSeries struct {
	Name string    `json:"name"`
	Data []float64 `json:"data"`
}

type ChartResponse struct {
	Series     []ChartSeries `json:"series"`
	Categories []string      `json:"categories"`
}
