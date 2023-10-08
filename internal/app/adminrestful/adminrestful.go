package adminrestful

import (
	"fmt"
	"github.com/josestg/swe-be-mono/internal/app"
)

// Run is the entrypoint of the admin-restful application.
func Run(app app.Info) error {
	fmt.Println(app.Name, "is running...")
	defer fmt.Println(app.Name, "has stopped.")
	return nil
}
