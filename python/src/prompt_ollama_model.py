import http.client
import json
from http.client import HTTPException
from io import StringIO
from pathlib import Path

from src.utils.utils import is_blank

FILE_NAME = "\"fileName\":"
SENTENCE1 = "The content of a text file follows."
SENTENCE2 = "Give a good file name for that file in a JSON format."
INSTRUCTION = f"{SENTENCE1} {SENTENCE2} The file name should be the property's {FILE_NAME} value as one string.\n\n"
TOKEN_INPUT_LIMIT = 2048 # default

# TODO: use python-docx

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
            "num_ctx": TOKEN_INPUT_LIMIT, # token input limit
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

def change_file_name(file_path, content):
  name = create_file_name_without_extension(content)

  if is_blank(name):
    return

  path = Path(file_path)
  new_file_name = name + ''.join(path.suffixes)

  path.rename(path.parent / new_file_name)

def get_txt_content(file_path, max_token_count):
  token_count = 0
  string_builder = StringIO()

  # Each line ends with the "\n" character, except the last line, if the file does not end with a newline.
  with open(file_path, "r") as lines:
    for line in lines:
      for token in basic_western_token_generator(line):
        string_builder.write(token)
        token_count += 1
        if token_count == max_token_count:
          return string_builder.getvalue()

  return string_builder.getvalue()

def approximate_western_token_count(text):
  token_count = 0

  for _ in basic_western_token_generator(text):
    token_count += 1

  return token_count

# This function generates tokens from Western text.
# It outputs tokens for words, whitespace characters, and punctuation marks.
# Note: This function does not support sub-word tokenization.
# For example, "unhappiness" is treated as a single token, not two tokens ("un" and "happiness").
def basic_western_token_generator(text):
  index = 0

  def is_space_or_punctuation(c):
    return c.isspace() or c in {',', '.', '?', '!', ';', ':', '(', ')', '[', ']'}

  while index < len(text):
    if is_space_or_punctuation(text[index]):
      yield text[index]
      index += 1
    else:
      string_builder = StringIO()
      while index < len(text) and not is_space_or_punctuation(text[index]):
        string_builder.write(text[index])
        index += 1
      yield string_builder.getvalue()

def change_file_name_by_content(file_path):
  content = get_txt_content(file_path, TOKEN_INPUT_LIMIT - approximate_western_token_count(INSTRUCTION))

  if is_blank(content):
    return

  return change_file_name(file_path, content)

change_file_name_by_content("C:\\Users\\testUser\\Desktop\\test\\test.txt")
