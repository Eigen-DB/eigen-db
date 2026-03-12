package main

import "controller/api"

func main() {
	// customers := []string{"spotify", "apple", "google", "meta"}
	// for _, c := range customers {
	// 	j := utils.JailFactory(c)
	// 	if err := j.Start(); err != nil {
	// 		panic(err)
	// 	}
	// }

	r := api.SetupRouter()
	if err := r.Run("0.0.0.0:8080"); err != nil {
		panic(err)
	}
}
