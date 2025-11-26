package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Alumni struct {
    ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    NIM string `bson:"nim" json:"nim"`
    Nama       string             `bson:"nama" json:"nama"`
    Jurusan    string             `bson:"jurusan" json:"jurusan"`
    Angkatan   int                `bson:"angkatan" json:"angkatan"`
    TahunLulus int                `bson:"tahun_lulus" json:"tahun_lulus"`
    Email      string             `bson:"email" json:"email"`
    NoTelepon  int             `bson:"no_telepon,omitempty" json:"no_telepon,omitempty"`
    Alamat     string             `bson:"alamat,omitempty" json:"alamat,omitempty"`
    CreatedAt  string          `bson:"created_at" json:"created_at"`
    UpdatedAt  string         `bson:"updated_at" json:"updated_at"`
}

type MetaInfo struct {
    Page   int    `json:"page"`
    Limit  int    `json:"limit"`
    Total  int    `json:"total"`
    Pages  int    `json:"pages"`
    SortBy string `json:"sortBy"`
    Order  string `json:"order"`
    Search string `json:"search"`
}

type AlumniResponse struct {
    Data []Alumni `json:"data"`
    Meta MetaInfo `json:"meta"`
}
