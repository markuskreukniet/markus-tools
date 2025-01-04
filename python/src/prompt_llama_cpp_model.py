from src.utils.utils import create_model_configuration, INSTRUCTION

NUMBER_OF_CONTEXT_TOKENS = 2048 # default

def prompt_llama_cpp_model(content, number_of_input_tokens):
  configuration = create_model_configuration(NUMBER_OF_CONTEXT_TOKENS, number_of_input_tokens)

  model = Llama(
    model_path = r"C:\Users\testUser\Downloads\bartowski Llama-3.2-3B-Instruct-Q4_K_M.gguf",
    prompt = f"{INSTRUCTION}{content}", # TODO: duplicate
    stop = None,

    n_ctx = configuration.number_of_context_tokens,
  )
