// package model

// type Pekerjaan struct {
// 	ID                  int     `json:"id"`
// 	AlumniID            int     `json:"alumni_id"`
// 	NamaPerusahaan      string  `json:"nama_perusahaan"`
// 	PosisiJabatan       string  `json:"posisi_jabatan"`
// 	BidangIndustri      string  `json:"bidang_industri"`
// 	LokasiKerja         string  `json:"lokasi_kerja"`
// 	GajiRange           string  `json:"gaji_range,omitempty"`
// 	TanggalMulaiKerja   string  `json:"tanggal_mulai_kerja"`             // YYYY-MM-DD
// 	TanggalSelesaiKerja *string `json:"tanggal_selesai_kerja,omitempty"` // YYYY-MM-DD
// 	StatusPekerjaan     string  `json:"status_pekerjaan"`
// 	IsDellete           string  `json:"isdellete"`
// 	Deskripsi           string  `json:"deskripsi_pekerjaan,omitempty"`
// }
// type isdell struct {
// 	IsDellete string `json:"isdellete"`
// }
// type JumlahPekerjaanPerTahun struct {
// 	Tahun  int `json:"tahun"`
// 	Jumlah int `json:"jumlah"`
// }

// type PekerjaanResponse struct {
// 	Data []Pekerjaan `json:"data"`
// 	Meta MetaInfo    `json:"meta"`
// }

package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Pekerjaan struct {
    ID                  primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    LegacyID         int                `bson:"id" json:"id"`
    AlumniID            int `bson:"alumni_id" json:"alumni_id"` 
    NamaPerusahaan      string             `bson:"nama_perusahaan" json:"nama_perusahaan"`
    PosisiJabatan       string             `bson:"posisi_jabatan" json:"posisi_jabatan"`
    BidangIndustri      string             `bson:"bidang_industri" json:"bidang_industri"`
    LokasiKerja         string             `bson:"lokasi_kerja" json:"lokasi_kerja"`
    GajiRange           string             `bson:"gaji_range,omitempty" json:"gaji_range,omitempty"`
    TanggalMulaiKerja   string             `bson:"tanggal_mulai_kerja" json:"tanggal_mulai_kerja"`
    TanggalSelesaiKerja *string            `bson:"tanggal_selesai_kerja,omitempty" json:"tanggal_selesai_kerja,omitempty"`
    StatusPekerjaan     string             `bson:"status_pekerjaan" json:"status_pekerjaan"`
    IsDellete            string               `bson:"isdellete" json:"isdellete"`
    Deskripsi           string             `bson:"deskripsi_pekerjaan,omitempty" json:"deskripsi_pekerjaan,omitempty"`
    CreatedAt           any          `bson:"created_at" json:"created_at"`
	UpdatedAt           any         `bson:"updated_at" json:"updated_at"`
}

type JumlahPekerjaanPerTahun struct {
    Tahun  int `json:"tahun"`
    Jumlah int `json:"jumlah"`
}

type PekerjaanResponse struct {
    Data []Pekerjaan `json:"data"`
    Meta MetaInfo    `json:"meta"`
}
