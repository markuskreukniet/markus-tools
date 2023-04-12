import http.client
import json
from http.client import HTTPException

from src.utils.utils import create_model_configuration, create_prompt, NUMBER_OF_CONTEXT_TOKENS

def prompt_ollama_model(content, number_of_input_tokens):
  configuration = create_model_configuration(NUMBER_OF_CONTEXT_TOKENS, number_of_input_tokens)
  connection = http.client.HTTPConnection("localhost:11434")

  connection.request(
    method="POST",
    url="/api/generate",
    headers={
      "Content-Type": "application/json"
    },

    body=json.dumps({
      "model": "llama3.2:3B",
      "prompt": create_prompt(content),
      "stream": False,

      "num_ctx": configuration.number_of_context_tokens,
      "num_predict": configuration.max_output_tokens,

      "seed": configuration.seed,
      "tfs_z": configuration.tfs_z,
      "temperature": configuration.temperature,

      "top_k": configuration.top_k,
      "top_p": configuration.top_p,
      "min_p": configuration.min_p,

      "repeat_last_n": configuration.repeat_last_n,
      "repeat_penalty": configuration.repeat_penalty,

      "mirostat": configuration.mirostat_mode,
      "mirostat_eta": configuration.mirostat_eta,
      "mirostat_tau": configuration.mirostat_tau
    }),
  )

  result = ""
  response = connection.getresponse()

  if response.status == 200:
      result = json.loads(response.read().decode("utf-8")).get("response", "")
  else:
    if connection:
      connection.close()
      raise HTTPException("no 200 status code")

  connection.close()

  return result
