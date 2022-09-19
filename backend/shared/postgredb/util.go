package postgredb

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
)

//lint:ignore U1000 Ignore unused function temporarily for debugging
// main functions in shared directory are used as example usage.
// to run main function, rename the package as main and then run: go run shared/lib_name/*.go
func main() {
	example_connectionPool(context.Background())

	// 	var name string
	// 	var weight int64
	// 	err = conn.QueryRow(ctx, "select name, weight from test_table where id=$1", 42).Scan(&name, &weight)
	// 	if err != nil {
	// 		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	// 		os.Exit(1)
	// 	}

	// 	fmt.Println(name, weight)
}

func AcquireConnectionToMainDatabase(ctx context.Context) (conn *pgx.Conn, err error) {
	dbUserName := viper.GetString("DB_USERNAME")
	dbPassword := viper.GetString("DB_PASSWORD")
	dbIpAddress := viper.GetString("DB_IP_ADDRESS")
	dbPort := viper.GetString("DB_PORT")
	dbName := viper.GetString("DB_NAME")
	databaseUrl, loggableDatabaseUrl := GenerateDatabaseUrl(
		dbUserName,
		dbPassword,
		dbIpAddress,
		dbPort,
		dbName,
	)
	return AcquireConnection(ctx, databaseUrl, loggableDatabaseUrl)
}

func GenerateDatabaseUrl(
	dbUserName,
	dbPassword,
	dbIpAddress,
	dbPort,
	dbName string,
) (databaseUrl, loggableDatabaseUrl string) {
	databaseUrl = fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v",
		dbUserName,
		dbPassword,
		dbIpAddress,
		dbPort,
		dbName,
	)
	loggableDatabaseUrl = fmt.Sprintf(
		"%v:%v/%v as user: %v",
		dbIpAddress,
		dbPort,
		dbName,
		dbUserName,
	)
	return
}

func AcquireConnection(ctx context.Context, databaseUrl, loggableDatabaseUrl string) (conn *pgx.Conn, err error) {
	conn, err = pgx.Connect(ctx, databaseUrl)
	if err != nil {
		err = fmt.Errorf("error connecting to database %v: %v", loggableDatabaseUrl, err)
		return conn, err
	}
	return conn, err
}

func example_connectToDb(ctx context.Context) {
	conn, err := AcquireConnectionToMainDatabase(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close(ctx)
}

func example_helloWorld(ctx context.Context) {
	conn, err := AcquireConnectionToMainDatabase(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close(ctx)

	var greeting string
	err = conn.QueryRow(ctx, "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(greeting)
}

func example_connectionPool(ctx context.Context) {
	dbUserName := viper.GetString("DB_USERNAME")
	dbPassword := viper.GetString("DB_PASSWORD")
	dbIpAddress := viper.GetString("DB_IP_ADDRESS")
	dbPort := viper.GetString("DB_PORT")
	dbName := viper.GetString("DB_NAME")
	databaseUrl, loggableDatabaseUrl := GenerateDatabaseUrl(
		dbUserName,
		dbPassword,
		dbIpAddress,
		dbPort,
		dbName,
	)
	fmt.Printf("connecting to %v\n", loggableDatabaseUrl)
	dbpool, err := pgxpool.Connect(ctx, databaseUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error connecting to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	var greeting string
	err = dbpool.QueryRow(ctx, "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(greeting)
}

func example_singleInsert(ctx context.Context) {
	conn, err := AcquireConnectionToMainDatabase(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close(ctx)

	sqlStr := fmt.Sprintf(`
		INSERT INTO public.users
		(id, "name", email, "hashed_password")
		VALUES('%v', '%v', '%v', '%v');
	`, "uuid", "name", "email", "hashedpassword")

	res, err := conn.Exec(ctx, sqlStr)
	if err != nil {
		err = fmt.Errorf("error executing query: %v", err)
		fmt.Println(err)
		return
	}
	fmt.Println("rows affected:", res.RowsAffected())
}

func example_singleSelect(ctx context.Context) {
	conn, err := AcquireConnectionToMainDatabase(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close(ctx)

	var id string
	var name string
	var email string
	var hashedPassword string
	err = conn.QueryRow(
		ctx,
		"SELECT id, name, email, hashed_password FROM public.users where email=$1",
		"local-kd@yopmail.com").
		Scan(&id, &name, &email, &hashedPassword)
	if err != nil {
		err = fmt.Errorf("error executing query: %v", err)
		fmt.Println(err)
		return
	}
}
