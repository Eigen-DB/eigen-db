package vector_io

import "eigen_db/cfg"

func SetupDB(config cfg.IConfig) {
	instantiateVectorStore(
		config.GetHNSWParamsDimensions(),
		config.GetHNSWParamsSimilarityMetric(),
		config.GetHNSWParamsSpaceSize(),
		config.GetHNSWParamsM(),
		config.GetHNSWParamsEfConstruction(),
	)
}
