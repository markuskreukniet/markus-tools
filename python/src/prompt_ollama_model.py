import http.client
import json
from io import StringIO
from pathlib import Path

FILE_NAME = "\"fileName\":"

# TODO: error handling
def prompt_ollama_model():
    prompt = f"The content of a text file follows. Give a good file name for that file in a JSON format. The file name should be the property's {FILE_NAME} value as one string.\n\nThe cat (Felis catus), or domestic cat, is a small carnivorous mammal and the only domesticated species in the Felidae family. Domesticated around 7500 BC in the Near East, cats are valued as pets and for controlling vermin. They are agile hunters with retractable claws, sharp teeth, excellent night vision, and a keen sense of smell. Though social, cats hunt alone, often at dawn and dusk. They communicate through vocalizations (meowing, purring, hissing) and body language, can hear high-frequency sounds, and use pheromones for signaling."

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
            "num_ctx": 2048,
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

def create_file_name_without_extension():
  prompt_response = prompt_ollama_model()

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

def change_file_name(file_path):
  path = Path(file_path)

  # path.rename("")
  new_file_name = create_file_name_without_extension() + ''.join(path.suffixes)

  return path.parent / new_file_name

# TODO: it should count tokens instead of white spaces
def get_txt_content(file_path):
  max_white_space_count = 2048 # TODO:
  white_space_count = 0
  string_builder = StringIO()

  with open(file_path, "r") as lines:
    for line in lines:
      string_builder.write(line) # TODO: could add to much
      white_space_count += sum(1 for char in line if char.isspace())
      if white_space_count >= max_white_space_count:
        break

  return string_builder.getvalue()

print(change_file_name("C:\\Users\\testUser\\Desktop\\test\\test.txt"))
