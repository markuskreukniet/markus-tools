import http.client
import json
from http.client import HTTPException
from pathlib import Path

from src.read_file_token_content import basic_western_token_generator, get_file_content
from src.utils.utils import is_blank

FILE_NAME = "\"fileName\":"
SENTENCE1 = "The content of a text file follows."
SENTENCE2 = "Give a good file name for that file in a JSON format."
INSTRUCTION = f"{SENTENCE1} {SENTENCE2} The file name should be the property's {FILE_NAME} value as one string.\n\n"
NUMBER_OF_CONTEXT_TOKENS = 2048 # default

def prompt_ollama_model(content):
    prompt = f"{INSTRUCTION}{content}"

    connection = http.client.HTTPConnection("localhost:11434")

    connection.request(
        method="POST",
        url="/api/generate",
        headers={
            "Content-Type": "application/json"
        },

        body=json.dumps({
            "model": "llama3.2:3B",
            "prompt": prompt,
            "stream": False,
            "temperature": 0.1,
            "top_p": 0.6,

            # defaults:
            "mirostat": 0,
            "mirostat_eta": 0.1,
            "mirostat_tau": 0.5,
            "repeat_last_n": 64,
            "repeat_penalty": 1.1,
            "num_predict": -1,
            "tfs_z": 1,
            "seed": 0,
            "top_k": 40,
            "num_ctx": NUMBER_OF_CONTEXT_TOKENS,
            "min_p": 0.0
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

def create_file_name_without_extension(content):
  prompt_response = prompt_ollama_model(content)

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

def change_file_name(path, content):
  name = create_file_name_without_extension(content)

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

  content = get_file_content(
    file_path, round(NUMBER_OF_CONTEXT_TOKENS * 0.6) - approximate_western_token_count(INSTRUCTION)
  )

  if is_blank(content):
    return

  change_file_name(path, content)

change_file_name_by_content("C:\\Users\\testUser\\Desktop\\test\\test.txt")
