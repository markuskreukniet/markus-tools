from llama_cpp import Llama

from src.utils.utils import create_model_configuration, create_prompt, NUMBER_OF_CONTEXT_TOKENS

def prompt_llama_cpp_model(content, number_of_input_tokens, model_file_path):
  configuration = create_model_configuration(NUMBER_OF_CONTEXT_TOKENS, number_of_input_tokens)

  model = Llama(
    model_path=model_file_path,
    n_ctx=configuration.number_of_context_tokens
  )

  response = model(
    prompt=create_prompt(content),
    stop=None,

    max_tokens=configuration.max_output_tokens,

    seed=configuration.seed, # TODO: 0 is everytime the same result
    tfs_z=configuration.tfs_z,
    temperature=configuration.temperature,

    top_k=configuration.top_k,
    top_p=configuration.top_p,
    min_p=configuration.min_p,

    # It does not have 'repeat_last_n'.
    repeat_penalty=configuration.repeat_penalty,

    mirostat_mode=configuration.mirostat_mode,
    mirostat_eta=configuration.mirostat_eta,
    mirostat_tau=configuration.mirostat_tau
  )

  choices = response.get("choices", [])

  if isinstance(choices, list) and choices:
    return choices[0].get("text", "")
  else:
    return ""
