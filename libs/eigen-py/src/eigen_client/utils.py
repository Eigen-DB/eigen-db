from eigen_client.supported_models import SUPPORTED_MODELS, SUPPORTED_MODEL_PROVIDERS

def validate_embedding_model(model_name: str, model_provider: str) -> dict:
    if model_name is None or model_provider is None or model_provider == "none":
        return {}
    model_names = [model['name'] for model in SUPPORTED_MODELS if model_provider == model['supported_provider']]
    if model_provider not in SUPPORTED_MODEL_PROVIDERS:
        raise ValueError(f"Invalid model provider: {model_provider}. Supported providers are: {', '.join(SUPPORTED_MODEL_PROVIDERS)}.")
    if model_name not in model_names:
        raise ValueError(f"Invalid model name: {model_name} for provider {model_provider}. Supported models for this provider are: {', '.join(model_names)}.\nIf the model is not supported, use create_index instead with the desired dimensions and metric.")
    return [model for model in SUPPORTED_MODELS if model['name'] == model_name][0]