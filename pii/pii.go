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
	ID                primitive.ObjectID   `schema:"_id" bson:"_id,omitempty" json:"_id,omitempty"`
	Nik               string               `schema:"NIK,omitempty" bson:"nik,omitempty" json:"NIK,omitempty"`
	EktpStatus        bool                 `schema:"EKTP_STATUS,omitempty" bson:"ektp_status,omitempty" json:"EKTP_STATUS,omitempty"`
	NamaLengkap       string               `schema:"NAMA_LENGKAP,omitempty" bson:"nama_lengkap,omitempty" json:"NAMA_LENGKAP,omitempty"`
	NamaLengkapIbu    string               `schema:"NAMA_LENGKAP_IBU,omitempty" bson:"nama_lengkap_ibu,omitempty" json:"NAMA_LENGKAP_IBU,omitempty"`
	NoHp              string               `schema:"NOMOR_HANDPHONE,omitempty" bson:"nomor_handphone,omitempty" json:"NOMOR_HANDPHONE,omitempty"`
	TanggalLahir      string               `schema:"TANGGAL_LAHIR,omitempty" bson:"tanggal_lahir,omitempty" json:"TANGGAL_LAHIR,omitempty"`
	TempatLahir       string               `schema:"TEMPAT_LAHIR,omitempty" bson:"tempat_lahir,omitempty" json:"TEMPAT_LAHIR,omitempty"`
	PendidikanAkhir   string               `schema:"PENDIDIKAN_AKHIR,omitempty" bson:"pendidikan_akhir,omitempty" json:"PENDIDIKAN_AKHIR,omitempty"`
	NoKK              string               `schema:"NOMOR_KARTU_KELUARGA,omitempty" bson:"nomor_kartu_keluarga,omitempty" json:"NOMOR_KARTU_KELUARGA,omitempty"`
	Alamat            string               `schema:"ALAMAT,omitempty" bson:"alamat,omitempty" json:"ALAMAT,omitempty"`
	Rt                string               `schema:"RT,omitempty" bson:"rt,omitempty" json:"RT,omitempty"`
	Rw                string               `schema:"RW,omitempty" bson:"rw,omitempty" json:"RW,omitempty"`
	NomorKelurahan    string               `schema:"NOMOR_KELURAHAN,omitempty" bson:"nomor_kelurahan,omitempty" json:"NOMOR_KELURAHAN,omitempty"`
	Kelurahan         string               `schema:"KELURAHAN,omitempty" bson:"kelurahan,omitempty" json:"KELURAHAN,omitempty"`
	NomorKecamatan    string               `schema:"NOMOR_KECAMATAN,omitempty" bson:"nomor_kecamatan,omitempty" json:"NOMOR_KECAMATAN,omitempty"`
	Kecamatan         string               `schema:"KECAMATAN,omitempty" bson:"kecamatan,omitempty" json:"KECAMATAN,omitempty"`
	NomorKabupaten    string               `schema:"NOMOR_KABUPATEN,omitempty" bson:"nomor_kabupaten,omitempty" json:"NOMOR_KABUPATEN,omitempty"`
	Kabupaten         string               `schema:"KABUPATEN,omitempty" bson:"kabupaten,omitempty" json:"KABUPATEN,omitempty"`
	NomorProvinsi     string               `schema:"NOMOR_PROVINSI,omitempty" bson:"nomor_provinsi,omitempty" json:"NOMOR_PROVINSI,omitempty"`
	Provinsi          string               `schema:"PROVINSI,omitempty" bson:"provinsi,omitempty" json:"PROVINSI,omitempty"`
	Agama             string               `schema:"AGAMA,omitempty" bson:"agama,omitempty" json:"AGAMA,omitempty"`
	Pekerjaan         string               `schema:"PEKERJAAN,omitempty" bson:"pekerjaan,omitempty" json:"PEKERJAAN,omitempty"`
	StatusPerkawinan  string               `schema:"STATUS_PERKAWINAN,omitempty" bson:"status_perkawinan,omitempty" json:"STATUS_PERKAWINAN,omitempty"`
	FotoKTP           *piimage.ImageStruct `schema:"FOTO_KTP,omitempty" bson:"foto_ktp,omitempty" json:"FOTO_KTP,omitempty"`
	FotoSelfie        *piimage.ImageStruct `schema:"FOTO_SELFIE,omitempty" bson:"foto_selfie,omitempty" json:"FOTO_SELFIE,omitempty"`
	FotoSelfieWithKTP *piimage.ImageStruct `schema:"FOTO_SELFIE_KTP" bson:"foto_selfie_ktp,omitempty" json:"FOTO_SELFIE_KTP,omitempty"`
	PasfotoKTP        *piimage.ImageStruct `schema:"PASFOTO_KTP,omitempty" bson:"pasfoto_ktp,omitempty" json:"PASFOTO_KTP,omitempty"`
}

func openPiiCollection() (context.Context, context.CancelFunc, *mongo.Client, *mongo.Collection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dburl))
	collection := client.Database(dbname).Collection(dbcoll)

	if err != nil {
		log.Println(err)
	}

	return ctx, cancel, client, collection, err
}

// DecodeFormPost decode the formPost data in requests form-data and assign it to Pii Struct
func DecodeFormPost(r *http.Request) (*Pii, error) {

	r.ParseMultipartForm(10 << 20)

	fd := r.PostForm
	newPii := new(Pii)
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	if fd.Get("TANGGAL_LAHIR") != "" {
		parsedtgllahir, errParse := time.Parse("2006-01-02", fd.Get("TANGGAL_LAHIR"))
		if errParse != nil {
			return nil, errors.New("Cannot parse tanggallahir decode")
		}
		tgllahirString := parsedtgllahir.String()
		fd.Set("TANGGAL_LAHIR", tgllahirString)
	}

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

	exist, _ := p.Exist()
	if exist {
		return nil, errors.New("Pii data exist, Pii.ID has been set")
	}

	ctx, cancel, _, collection, _ := openPiiCollection()
	res, err := collection.InsertOne(ctx, p)
	defer cancel()

	if err != nil {
		return nil, err
	}

	newid := &p.ID
	*newid = res.InsertedID.(primitive.ObjectID)

	return p.ID, nil
}

// Exist Check Pii data existance in local database
func (p *Pii) Exist() (bool, error) {
	_, cancel, _, collection, _ := openPiiCollection()
	decodepoint := new(Pii)
	err := collection.FindOne(context.TODO(), bson.M{"nik": p.Nik}).Decode(decodepoint)
	defer cancel()

	if err != nil {
		log.Println(err)
		return false, err
	}

	pointerID := &p.ID
	*pointerID = decodepoint.ID

	return true, nil
}

// GrepData gerp all pii data by it current _id or nik
func (p *Pii) GrepData() error {
	_, cancel, client, _ := openConnection()
	collection := client.Database(dbname).Collection(dbcoll)
	decodepoint := new(Pii)
	err := collection.FindOne(context.TODO(), bson.M{"nik": p.Nik}).Decode(decodepoint)
	defer cancel()

	if err != nil {
		fmt.Println("FindOne Result: ", decodepoint)
		return err
	}

	fmt.Printf("%v\n", decodepoint)
	return nil
}

func openConnection() (context.Context, context.CancelFunc, *mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dburl))
	return ctx, cancel, client, err
}
