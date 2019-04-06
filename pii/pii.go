package pii

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/adhityasan/gomongo-rest/config"

	"github.com/adhityasan/gomongo-rest/pii/piimage"

	"github.com/gorilla/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbhost = config.Of.Mongo.Host
var dbport = config.Of.Mongo.Port
var dburl = config.Of.Mongo.URL
var dbname = config.Of.DBModules["pii"].Db
var dbcoll = config.Of.DBModules["pii"].Coll

// Pii stands for Personal Identifying Information
type Pii struct {
	ID                primitive.ObjectID   `schema:"_id,omitempty" bson:"_id,omitempty" json:"_id,omitempty"`
	Nik               string               `schema:"NIK,omitempty" bson:"nik,omitempty" json:"NIK,omitempty"`
	EktpStatus        bool                 `schema:"EKTP_STATUS,omitempty" bson:"ektp_status,omitempty" json:"EKTP_STATUS,omitempty"`
	NamaLengkap       string               `schema:"NAMA_LENGKAP,omitempty" bson:"nama_lengkap,omitempty" json:"NAMA_LENGKAP,omitempty"`
	NoHp              string               `schema:"NO_HP,omitempty" bson:"no_hp,omitempty" json:"NO_HP,omitempty"`
	TanggalLahir      string               `schema:"TANGGAL_LAHIR,omitempty" bson:"tanggal_lahir,omitempty" json:"TANGGAL_LAHIR,omitempty"`
	TempatLahir       string               `schema:"TEMPAT_LAHIR,omitempty" bson:"tempat_lahir,omitempty" json:"TEMPAT_LAHIR,omitempty"`
	PendidikanAkhir   string               `schema:"PENDIDIKAN_AKHIR,omitempty" bson:"pendidikan_akhir,omitempty" json:"PENDIDIKAN_AKHIR,omitempty"`
	NoKK              string               `schema:"NO_KK,omitempty" bson:"no_kk,omitempty" json:"NO_KK,omitempty"`
	Alamat            string               `schema:"ALAMAT,omitempty" bson:"alamat,omitempty" json:"ALAMAT,omitempty"`
	Rt                string               `schema:"RT,omitempty" bson:"rt,omitempty" json:"RT,omitempty"`
	Rw                string               `schema:"RW,omitempty" bson:"rw,omitempty" json:"RW,omitempty"`
	Kecamatan         string               `schema:"KECAMATAN,omitempty" bson:"kecamatan,omitempty" json:"KECAMATAN,omitempty"`
	Kabupaten         string               `schema:"KABUPATEN,omitempty" bson:"kabupaten,omitempty" json:"KABUPATEN,omitempty"`
	Provinsi          string               `schema:"PROVINSI,omitempty" bson:"provinsi,omitempty" json:"PROVINSI,omitempty"`
	Agama             string               `schema:"AGAMA,omitempty" bson:"agama,omitempty" json:"AGAMA,omitempty"`
	Pekerjaan         string               `schema:"PEKERJAAN,omitempty" bson:"pekerjaan,omitempty" json:"PEKERJAAN,omitempty"`
	StatusPerkawinan  string               `schema:"STATUS_PERKAWINAN,omitempty" bson:"status_perkawinan,omitempty" json:"STATUS_PERKAWINAN,omitempty"`
	FotoKTP           *piimage.ImageStruct `schema:"FOTO_KTP,omitempty" bson:"foto_ktp,omitempty" json:"FOTO_KTP,omitempty"`
	FotoSelfie        *piimage.ImageStruct `schema:"FOTO_SELFIE,omitempty" bson:"foto_selfie,omitempty" json:"FOTO_SELFIE,omitempty"`
	FotoSelfieWithKTP *piimage.ImageStruct `schema:"FOTO_SELFIE_KTP" bson:"foto_selfie_ktp,omitempty" json:"FOTO_SELFIE_KTP,omitempty"`
	PasfotoKTP        *piimage.ImageStruct `schema:"PASFOTO_KTP,omitempty" bson:"pasfoto_ktp,omitempty" json:"PASFOTO_KTP,omitempty"`
}

// GetLocalPii func to get Personal Information based on given nik (param)
// return (Pii, nil) Struct , and (nil, error) if data is not exist or somethong went wrong
func GetLocalPii(nik string) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI(dburl))
	collection := client.Database(dbname).Collection(dbcoll)
	cursor, _ := collection.Find(ctx, bson.M{"nik": nik})

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person Pii
		cursor.Decode(&person)
		fmt.Println("nik : ", person.Nik)
	}

}

// DecodeFormPost decode the formPost data in requests form-data and assign it to Pii Struct
func DecodeFormPost(r *http.Request) (*Pii, error) {

	r.ParseMultipartForm(10 << 20)

	fd := r.PostForm
	newPii := new(Pii)
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	parsedtgllahir, errParse := time.Parse("2006-01-02", fd.Get("TANGGAL_LAHIR"))
	if errParse != nil {
		return nil, errors.New("Cannot parse tanggallahir decode")
	}

	tgllahirString := parsedtgllahir.String()
	fd.Set("tanggal_lahir", tgllahirString)

	err := decoder.Decode(newPii, fd)

	if err != nil {
		errdetail := fmt.Sprintf("Fail to decode request form-data to new Pii data Struct : %s\n", err)
		return nil, errors.New(errdetail)
	}

	newPii.FotoKTP, err = piimage.ImageStructHandler("FOTO_KTP", r)
	newPii.FotoSelfie, err = piimage.ImageStructHandler("FOTO_SELFIE", r)
	newPii.FotoSelfieWithKTP, err = piimage.ImageStructHandler("FOTO_SELFIE_KTP", r)
	newPii.PasfotoKTP, err = piimage.ImageStructHandler("PASFOTO_KTP", r)

	return newPii, nil
}

// Save current Personal Identifying Information
func (p *Pii) Save() (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI(dburl))
	collection := client.Database(dbname).Collection(dbcoll)
	res, err := collection.InsertOne(ctx, p)
	defer cancel()

	if err != nil {
		return nil, err
	}

	id := res.InsertedID

	return id, nil
}

// Exist Check Pii data existance in local database
func (p *Pii) Exist() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI(dburl))
	collection := client.Database(dbname).Collection(dbcoll)
	a := new(Pii)
	err := collection.FindOneAndReplace(context.TODO(), bson.M{"nik": p.Nik}, p).Decode(a)
	defer cancel()

	if err != nil {
		return false, err
	}

	fmt.Printf("%v\n", a)
	return true, nil
}
