import http.client
import json
from io import StringIO
from pathlib import Path

FILE_NAME = "\"fileName\":"
INSTRUCTION = f"The content of a text file follows. Give a good file name for that file in a JSON format. The file name should be the property's {FILE_NAME} value as one string."

# TODO: error handling
def prompt_ollama_model(content):
    prompt = f"{INSTRUCTION}\n\n{content}"

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
            "num_ctx": 2048, # token input limit
            "min_p": 0.0
        }),
    )

    result = ""
    response = connection.getresponse()

    if response.status == 200:
        result = json.loads(response.read().decode('utf-8')).get("response", "")
    else:
      if connection:
        connection.close()
        raise Exception("no 200 status code") # TODO: Exception not specific

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

  # print(path.parent / new_file_name)

# TODO: it should count tokens instead of white spaces
def get_txt_content(file_path, max_white_space_count):
  white_space_count = 0
  string_builder = StringIO()

  with open(file_path, "r") as lines:
    for line in lines:
      string_builder.write(line) # TODO: could add to much
      white_space_count += sum(1 for char in line if char.isspace())
      if white_space_count >= max_white_space_count:
        break

  return string_builder.getvalue()

def is_blank(s):
  return not s.strip()

# TODO:
# words, whitespaces, and punctuations. It does not handle sub words such as unhappiness, which are two tokens, 'un' and 'happiness'. It does handle not non western languages such as Japanese.
def approximate_western_token_count(text):
  token_count = 0
  index = 0

  def is_space_or_punctuation(c):
    return c.isspace() or c in {',', '.', '?', '!', ';', ':', '(', ')', '[', ']'}

  while index < len(text):
    if is_space_or_punctuation(text[index]):
      token_count += 1
      index += 1
    else:
      token_count += 1
      while index < len(text) and not is_space_or_punctuation(text[index]):
        index += 1

  return token_count

def change_file_name_by_content(file_path):
  content = get_txt_content(file_path, 2048) # TODO:

  if is_blank(content):
    return

  return change_file_name(file_path, content)

change_file_name_by_content("C:\\Users\\testUser\\Desktop\\test\\test.txt")
