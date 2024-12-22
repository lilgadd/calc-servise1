package main
import(
	"github.com/lilgadd/calc.go/internal/application"
)

func main(){
	app := application.New()
	//app.Run()
	app.RunServer()
}