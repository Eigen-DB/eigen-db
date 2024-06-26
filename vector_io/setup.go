package vector_io

import "eigen_db/cfg"

func SetupDB(config *cfg.Config) {
	params := config.HNSWParams
	instantiateVectorStore(
		params.Dimensions,
		params.SimilarityMetric,
		params.SpaceSize,
		params.M,
		params.EfConstruction,
	)
}
