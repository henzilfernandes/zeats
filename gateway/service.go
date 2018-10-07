package gateway

import (
	"github.com/gocql/gocql"
	"zeats/types"
	"errors"
	"github.com/golang/glog"
	"reflect"
	"strings"
)

const (
	insertQuery = `INSERT INTO zeats.eats (product_id, name, image_closed, image_open, description, story, 
		sourcing_values, ingredients, allergy_info, dietary_certifications) VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	deleteQuery = `DELETE FROM zeats.eats WHERE product_id = ?;`
	selectQuery = `SELECT * FROM zeats.eats WHERE product_id = ?;`

	jsonTag = `json`
)

type service struct {
	cassandraSession *gocql.Session
}

type Config struct {
	CassandraSession *gocql.Session
}

type Service interface {
	InsertEats(prod *types.Product) error
	UpdateEats(prod *types.Product) error
	DeleteEats(id string) error
	FetchEats(id string) (*types.Product, error)
}

func (s *service) InsertEats(prod *types.Product) error {
	if s.cassandraSession == nil {
		return errors.New("cassandra session is down")
	}

	if qErr := s.cassandraSession.Query(insertQuery,
		prod.Id,
		prod.Name,
		prod.ImageClosed,
		prod.ImageOpen,
		prod.Description,
		prod.Story,
		prod.SourcingValues,
		prod.Ingredients,
		prod.AllergyInfo,
		prod.DietaryCertifications).Consistency(gocql.One).Exec(); qErr != nil {

		glog.V(0).Infof("Product Insert/Update error :: %s", qErr)
		return qErr
	}
	return nil
}

func (s *service) DeleteEats(id string) error {
	if s.cassandraSession == nil {
		return errors.New("cassandra session is down")
	}

	if qErr := s.cassandraSession.Query(deleteQuery, id).Consistency(gocql.All).Exec(); qErr != nil {
		glog.V(0).Infof("Product Delete error :: ProductId :: %s :: Err :: %s", qErr)
		return qErr
	}
	return nil
}

func (s *service) FetchEats(id string) (*types.Product, error) {
	if s.cassandraSession == nil {
		return nil, errors.New("cassandra session is down")
	}

	prod := &types.Product{}
	qErr := s.cassandraSession.Query(selectQuery, id).Consistency(gocql.All).Scan(
		&prod.Id, &prod.AllergyInfo, &prod.Description, &prod.DietaryCertifications,
		&prod.ImageClosed, &prod.ImageOpen, &prod.Ingredients, &prod.Name, &prod.SourcingValues, &prod.Story)
	if qErr != nil {
		glog.V(0).Infof("Product Fetch error :: ProductId :: %s :: Err :: %s", qErr)
		return nil, qErr
	}
	return prod, nil
}

func (s *service) UpdateEats(prod *types.Product) error {
	if s.cassandraSession == nil {
		return errors.New("cassandra session is down")
	}

	if prod.Id == "" {
		return errors.New("product Id not specified in update")
	}

	cassandraBatch := s.cassandraSession.NewBatch(gocql.UnloggedBatch)
	cassandraBatch.Cons = gocql.One

	v := reflect.ValueOf(prod).Elem()
	for i := 0; i < v.NumField(); i++ {
		key := strings.Split(v.Type().Field(i).Tag.Get(jsonTag), ",")[0]
		value := v.Field(i).Interface()

		switch value.(type) {
		case int64: if value == 0 {
			continue
		}
		case string: if value == "" {
			continue
		}
		}

		if key == "productId" || value == nil {
			continue
		}

		cassandraBatch.Query("UPDATE zeats.eats SET " + key + "=" + "? WHERE product_id = ?", value, prod.Id)
	}

	if qErr := s.cassandraSession.ExecuteBatch(cassandraBatch); qErr != nil {
		glog.V(0).Infof("Product Update error :: ProductId :: %s :: Err :: %s", qErr)
		return qErr
	}

	return nil
}

func newService(config *Config) Service {
	return &service{
		cassandraSession: config.CassandraSession,
	}
}