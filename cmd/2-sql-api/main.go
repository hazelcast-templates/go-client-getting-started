package main

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hazelcast/hazelcast-go-client"
	"github.com/hazelcast/hazelcast-go-client/logger"
	"github.com/hazelcast/hazelcast-go-client/serialization"

	gettingstarted "getting-started"
)

type CityDTO struct {
	city       string
	country    string
	population int32
}

type CitySerializer struct{}

func (s CitySerializer) Type() reflect.Type {
	return reflect.TypeOf(CityDTO{})
}

func (s CitySerializer) TypeName() string {
	return "CityDTO"
}

func (s CitySerializer) Write(writer serialization.CompactWriter, value interface{}) {
	city := value.(CityDTO)

	writer.WriteString("City", &city.city)
	writer.WriteString("Country", &city.country)
	writer.WriteInt32("Population", city.population)
}

func (s CitySerializer) Read(reader serialization.CompactReader) interface{} {
	return CityDTO{
		city:       *reader.ReadString("city"),
		country:    *reader.ReadString("country"),
		population: reader.ReadInt32("population"),
	}
}

func createMapping(ctx context.Context, client hazelcast.Client) error {
	fmt.Println("Creating the mapping...")

	// Mapping is required for your distributed map to be queried over SQL.
	// See: https://docs.hazelcast.com/hazelcast/latest/sql/mapping-to-maps
	mappingQuery := `
        CREATE OR REPLACE MAPPING
        cities (
            __key INT,
            country VARCHAR,
            city VARCHAR,
            population INT) TYPE IMAP
        OPTIONS (
            'keyFormat' = 'int',
            'valueFormat' = 'compact',
            'valueCompactTypeName' = 'CityDTO')
    `

	_, err := client.SQL().Execute(ctx, mappingQuery)
	if err != nil {
		return err
	}

	fmt.Println("OK.\n")
	return nil
}

func populateCities(ctx context.Context, client hazelcast.Client) error {
	fmt.Println("Inserting data...")

	// Mapping is required for your distributed map to be queried over SQL.
	// See: https://docs.hazelcast.com/hazelcast/latest/sql/mapping-to-maps
	insertQuery := `
		INSERT INTO cities
		(__key, city, country, population) VALUES
		(1, 'London', 'United Kingdom', 9540576),
		(2, 'Manchester', 'United Kingdom', 2770434),
		(3, 'New York', 'United States', 19223191),
		(4, 'Los Angeles', 'United States', 3985520),
		(5, 'Istanbul', 'Türkiye', 15636243),
		(6, 'Ankara', 'Türkiye', 5309690),
		(7, 'Sao Paulo ', 'Brazil', 22429800)
    `

	_, err := client.SQL().Execute(ctx, "DELETE from cities")
	if err != nil {
		return err
	}
	_, err = client.SQL().Execute(ctx, insertQuery)
	if err != nil {
		return err
	}

	fmt.Println("OK.\n")
	return nil
}

func fetchCities(ctx context.Context, client hazelcast.Client) error {
	fmt.Println("Fetching cities...")

	result, err := client.SQL().Execute(ctx, "SELECT __key, this FROM cities")
	if err != nil {
		return err
	}
	defer result.Close()

	fmt.Println("OK.")
	fmt.Println("--Results of SELECT __key, this FROM cities")
	fmt.Printf("| %4s | %20s | %20s | %15s |\n", "id", "country", "city", "population")

	iter, err := result.Iterator()
	for iter.HasNext() {
		row, err := iter.Next()

		key, err := row.Get(0)
		cityDTO, err := row.Get(1)

		fmt.Printf("| %4d | %20s | %20s | %15d |\n", key.(int32), cityDTO.(CityDTO).country, cityDTO.(CityDTO).city, cityDTO.(CityDTO).population)

		if err != nil {
			return err
		}
	}

	fmt.Println("\n!! Hint !! You can execute your SQL queries on your Viridian cluster over the management center. \n 1. Go to 'Management Center' of your Hazelcast Viridian cluster. \n 2. Open the 'SQL Browser'. \n 3. Try to execute 'SELECT * FROM cities'.")
	return nil
}

///////////////////////////////////////////////////////

func main() {
	// Connection details for cluster
	config, err := gettingstarted.ClientConfig()
	if err != nil {
		panic(err)
	}

	// Register Compact Serializers
	config.Serialization.Compact.SetSerializers(CitySerializer{})

	// Other environment propreties
	config.Logger.Level = logger.DebugLevel

	ctx := context.TODO()
	// create the client and connect to the cluster
	client, err := hazelcast.StartNewClientWithConfig(ctx, config)
	if err != nil {
		panic(err)
	}

	if err := createMapping(ctx, *client); err != nil {
		panic(fmt.Errorf("creating mapping: %w", err))
	}
	if err := populateCities(ctx, *client); err != nil {
		panic(fmt.Errorf("populating cities: %w", err))
	}
	if err := fetchCities(ctx, *client); err != nil {
		panic(fmt.Errorf("fetching cities: %w", err))
	}

	if err := client.Shutdown(ctx); err != nil {
		panic(err)
	}
}
