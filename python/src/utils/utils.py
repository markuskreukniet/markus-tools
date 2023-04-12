from dataclasses import dataclass

FILE_NAME = "\"fileName\":"
SENTENCE1 = "The content of a text file follows."
SENTENCE2 = "Give a good file name for that file in a JSON format."
INSTRUCTION = f"{SENTENCE1} {SENTENCE2} The file name should be the property's {FILE_NAME} value as one string.\n\n"
NUMBER_OF_CONTEXT_TOKENS = 2048 # default

def is_blank(s):
  return not s.strip()

@dataclass
class ModelConfiguration:
  number_of_context_tokens: int # The model's maximum token limit includes input (prompt) and generated output.
  max_output_tokens: int

  # In Ollama, 'seed=0' sets a random seed; in llama.cpp, omit the seed for randomness.
  # However, when I omit the seed, llama.cpp gives unusable results.
  seed: int
  tfs_z: int
  temperature: float

  top_k: int
  top_p: float
  min_p: float

  repeat_last_n: int
  repeat_penalty: float

  # If mirostat_mode is 0, the Mirostat algorithm is disabled, and the model does not use mirostat_eta or mirostat_tau.
  mirostat_mode: int
  mirostat_eta: float
  mirostat_tau: float

def create_model_configuration(number_of_context_tokens, number_of_input_tokens):
  return ModelConfiguration(
    number_of_context_tokens=number_of_context_tokens,
    max_output_tokens=number_of_context_tokens - number_of_input_tokens,

    seed=0,
    tfs_z=1,
    temperature=0.1,

    top_k=40,
    top_p=0.6,
    min_p=0.1,

    repeat_last_n=64,
    repeat_penalty=1.1,

    mirostat_mode=0,
    mirostat_eta=0.1,
    mirostat_tau=0.5
  )

def create_prompt(content):
  return f"{INSTRUCTION}{content}"
