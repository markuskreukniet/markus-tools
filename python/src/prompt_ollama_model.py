import http.client
import json
from http.client import HTTPException
from pathlib import Path

from src.prompt_llama_cpp_model import NUMBER_OF_CONTEXT_TOKENS
from src.read_file_token_content import basic_western_token_generator, get_file_content
from src.utils.utils import is_blank, create_model_configuration, FILE_NAME, INSTRUCTION

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
      "prompt": f"{INSTRUCTION}{content}",
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

def create_file_name_without_extension(content, number_of_input_tokens):
  prompt_response = prompt_ollama_model(content, number_of_input_tokens)

  if prompt_response == "":
    return ""

  index = prompt_response.find(FILE_NAME)

  if index == -1:
    return ""

  def slice_and_find_index(response, start_index):
    response = response[start_index:]
    return response, response.find("\"")

  prompt_response, index = slice_and_find_index(prompt_response, index + len(FILE_NAME))
  prompt_response, index = slice_and_find_index(prompt_response, index+1)

  return prompt_response[:index]

def change_file_name(path, content, number_of_input_tokens):
  name = create_file_name_without_extension(content, number_of_input_tokens)

  if is_blank(name):
    return

  new_file_name = name + ''.join(path.suffixes)

  path.rename(path.parent / new_file_name)

def approximate_western_token_count(text):
  token_count = 0

  for _ in basic_western_token_generator(text):
    token_count += 1

  return token_count

def change_file_name_by_content(file_path):
  path = Path(file_path)

  if not path.is_file() and path.stat().st_size > 0:
    return

  number_of_input_tokens = round(NUMBER_OF_CONTEXT_TOKENS * 0.6)
  content = get_file_content(file_path, number_of_input_tokens - approximate_western_token_count(INSTRUCTION))

  if is_blank(content):
    return

  change_file_name(path, content, number_of_input_tokens)

change_file_name_by_content("C:\\Users\\testUser\\Desktop\\test\\test.txt")
