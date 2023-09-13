package main



func main(){
	app := App{}
	app.Initialize(DbUser, DbPassword, DbName)
	app.Run("localhost:8080")
}