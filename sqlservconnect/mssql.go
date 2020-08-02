package sqlservconnect

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

//defining resource type
type Movie struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Year string `json:"year"`
}

var db *sql.DB 

//server connection properties
var server = "localhost"
var port = 1433
var user = ""
var password = ""
var database = "tempdb"

func StartDatabaseServer() { 
	//connect to the ms sql server
	// Build connection string
	connection := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",server, user, password, port, database)

	var err error

	// Create connection pool
	db, err = sql.Open("sqlserver", connection)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	
	ctx := context.Background()
	err = db.PingContext(ctx) //check if database is alive or not
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Connected to SQL SERVER!\n")
}

// Create movie
func CreateMovie(name string, year string) (error) {
	ctx := context.Background()
	tsql := "INSERT INTO TestSchema.Movies (Name, Year) VALUES (@Name, @Year);"

	stmt, err := db.Prepare(tsql) //prepare statement to insert record in database
	if err != nil {
		return err
	}
	stmt.QueryRowContext(ctx, sql.Named("Name", name),sql.Named("Year", year)) //add named paramters and execute
	stmt.Close()
	return nil
}

// Read all Movies
func ReadMovies() ([]Movie, error) {
	ctx := context.Background()
	tsql := fmt.Sprintf("SELECT Id, Name, Year FROM TestSchema.Movies;")

	// Execute query
	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //close after reading all records

	var Movies []Movie
	// Iterate through the result set.
	for rows.Next() {
		var name, year string
		var id int
		var movie Movie
		// Get values from row.
		err := rows.Scan(&id, &name, &year)
		if err != nil {
			return nil, err
		}
		movie.Id, movie.Name , movie.Year = id , name , year
		Movies = append(Movies,movie)  //add into result array of type Movie
	}
	return Movies, nil
}

//read single movie
func ReadSingleMovie(id int) (*Movie, error) {
	var movie Movie

	ctx := context.Background()
	tsql := fmt.Sprintf("SELECT Id, Name, Year FROM TestSchema.Movies WHERE Id = @Id ;")

	// Execute query
	rows, err := db.QueryContext(ctx, tsql,sql.Named("Id",id)) //add named parameter and execute
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Get values from row.
	if rows.Next() {
		err = rows.Scan(&movie.Id, &movie.Name, &movie.Year) //scan result into struct movie
		if err != nil {
			return nil, err
		}
		return &movie, nil //return address of the struct 
	}
	return nil,nil
}

// Update Movie
func UpdateMovie(id int,name string, year string) (error) {
	ctx := context.Background()
	tsql := fmt.Sprintf("UPDATE TestSchema.Movies SET Name = @Name , Year = @Year WHERE Id = @Id")

	_, err := db.ExecContext(ctx,tsql,sql.Named("Name", name),sql.Named("Year", year),sql.Named("Id", id)) // Execute with named parameters
	if err != nil {
		return err
	}
	return nil
}

// Delete movie
func DeleteMovie(id int) (error) {
	ctx := context.Background()
	tsql := fmt.Sprintf("DELETE FROM TestSchema.Movies WHERE Id = @Id;")

	_, err := db.ExecContext(ctx, tsql, sql.Named("Id", id)) //Execute with Named Parameters
	if err != nil {
		return err
	}
	return nil
}
