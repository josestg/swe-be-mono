package enduserrestful

import (
	"fmt"

	"github.com/josestg/swe-be-mono/internal/app"
)

// Run is the entrypoint of the enduser-restful application.
func Run(app app.Info) error {
	fmt.Println(app.Name, "is running...")
	fmt.Printf("Build Info: %+v\n", app)
	defer fmt.Println(app.Name, "has stopped.")
	return nil
}